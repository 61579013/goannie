package godler

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/cheggaaa/pb/v3"
)

// DlerURL 基础结构体
type DlerURL struct {
	URL                   string        // 下载的url
	SavePath              string        // 保存路径带后缀
	OneThreading          bool          // 强制单线程
	ThreadingCount        int64         // 下载线程数 默认20
	ThreadingDownloadSize int64         // 每个区块下载的大小 默认1048576
	TimeOut               time.Duration //下载超时时间 默认 180 秒

	Method  string            // 请求下载地址的方式 默认GET
	Headers map[string]string // 请求下载地址的Headers 默认{}

	Client http.Client // 请求时的Cilent 默认{Timeout:time.Second * 180}
	Body   io.Reader   // 请求时的Body	默认nil

	IsBar bool // 是否显示进度条等信息 默认false

	ContentLength *int64 // 当前下载文件总字节数
	Schedule      *int64 // 当前已经下载的字节数
}

// InitHeaders 初始化Header
func (dlerurl *DlerURL) InitHeaders(request *http.Request) {
	for idx, val := range dlerurl.Headers {
		request.Header.Set(idx, val)
	}
}

// IsClientEmpty 判断Client是否为空
func (dlerurl *DlerURL) IsClientEmpty() bool {
	return reflect.DeepEqual(dlerurl.Client, http.Client{})
}

// downloadOver 下载完成时的回调
func (dlerurl *DlerURL) downloadOver(file *os.File) error {
	// 下载成功，更改临时文件名称
	err := dlerurl.reTempFileName(file)
	if err != nil {
		return err
	}
	return nil
}

// Init 初始化结构体
func (dlerurl *DlerURL) Init() {
	if dlerurl.ThreadingCount <= 0 {
		dlerurl.ThreadingCount = 20
	}
	if dlerurl.ThreadingDownloadSize <= 0 {
		dlerurl.ThreadingDownloadSize = 1048576
	}
	if dlerurl.Method == "" {
		dlerurl.Method = "GET"
	}
	if dlerurl.IsClientEmpty() {
		dlerurl.Client.Timeout = time.Second * 180
	}
	if dlerurl.TimeOut <= 0 {
		dlerurl.TimeOut = time.Second * 180
	}
	schedule := int64(0)
	dlerurl.Schedule = &schedule
}

// Download 开始下载
func (dlerurl *DlerURL) Download() error {
	dlerurl.Init()
	isThreading, contentLength, err := dlerurl.isMultithreading()
	if err != nil {
		return err
	}

	dlerurl.ContentLength = &contentLength

	file, err := dlerurl.creactFile()
	if err != nil {
		return err
	}
	// 是支持多线程下载的文件
	if isThreading && (!dlerurl.OneThreading) {
		err = dlerurl.threadingDownloadStart(file)
		if err != nil {
			return err
		}
	} else {
		err = dlerurl.downloadStart(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// downloadStart 单线程下载
func (dlerurl *DlerURL) downloadStart(file *os.File) error {
	defer file.Close()
	if dlerurl.IsBar {
		fmt.Println("单线程下载")
	}
	bar := &dlerBar{IsShow: dlerurl.IsBar}
	bar.dlerBarStart(*dlerurl.ContentLength)
	request, err := http.NewRequest(dlerurl.Method, dlerurl.URL, dlerurl.Body)
	if err != nil {
		bar.dlerBarFinish()
		return err
	}
	dlerurl.InitHeaders(request)
	res, err := dlerurl.Client.Do(request)
	if err != nil {
		bar.dlerBarFinish()
		return err
	}
	defer res.Body.Close()
	buf := make([]byte, 1048576)
	for {
		n, readErr := res.Body.Read(buf)
		_, writeErr := file.Write(buf[0:n])
		if writeErr != nil {
			return writeErr
		}
		atomic.AddInt64(dlerurl.Schedule, int64(n))
		bar.dlerBarAdd(int64(n))
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			bar.dlerBarFinish()
			return readErr
		}
	}
	bar.dlerBarFinish()
	if dlerurl.IsBar {
		fmt.Println("download succeed!")
	}
	// 下载完成后的回调
	err = dlerurl.downloadOver(file)
	if err != nil {
		return err
	}
	return nil
}

// creactFile 创建并获取文件
func (dlerurl *DlerURL) creactFile() (*os.File, error) {
	var file *os.File
	var startSize int64
	fileInfo, err := os.Stat(dlerurl.SavePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件或目录不存在，创建临时文件
			file, err = dlerurl.creactTempFile()
			if err != nil {
				return nil, err
			}
		}
	} else {
		startSize = fileInfo.Size()
		if startSize == *dlerurl.ContentLength {
			return nil, errors.New("文件已存在")
		}
		// 文件没有下载完整时，删除文件
		err = os.Remove(dlerurl.SavePath)
		if err != nil {
			return nil, err
		}
		// 创建临时下载文件
		file, err = dlerurl.creactTempFile()
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

// creactTempFile 创建临时下载文件
func (dlerurl *DlerURL) creactTempFile() (*os.File, error) {
	var file *os.File
	tempFile := fmt.Sprintf("%s.tempDownload", dlerurl.SavePath)
	_, err := os.Stat(tempFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件或目录不存在
			err = os.MkdirAll(path.Dir(dlerurl.SavePath), os.ModePerm)
			if err != nil {
				return file, err
			}
			file, err = os.Create(tempFile)
			if err != nil {
				return file, err
			}
		}
	} else {
		// 如果临时文件存在，先删除再创建
		err = os.Remove(tempFile)
		if err != nil {
			return file, err
		}
		file, err = os.Create(tempFile)
		if err != nil {
			return file, err
		}
	}
	return file, nil
}

// reTempFileName 修改临时下载文件名称
func (dlerurl *DlerURL) reTempFileName(file *os.File) error {
	_ = file.Close()
	tempFile := fmt.Sprintf("%s.tempDownload", dlerurl.SavePath)
	fileInfo, err := os.Stat(dlerurl.SavePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在,直接改名字
			err = os.Rename(tempFile, dlerurl.SavePath)
			if err != nil {
				return err
			}
		}
	} else {
		startSize := fileInfo.Size()
		if startSize == *dlerurl.ContentLength {
			// 文件已经存在了，删除临时文件
			err = os.Remove(tempFile)
			if err != nil {
				return err
			}
			return errors.New("文件已存在")
		} else {
			// 文件残缺，先删除原文件，再改名字
			err = os.Remove(dlerurl.SavePath)
			if err != nil {
				return err
			}
			err = os.Rename(tempFile, dlerurl.SavePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// isMultithreading 判断url是否支持多线程
func (dlerurl *DlerURL) isMultithreading() (bool, int64, error) {
	reqHead, err := http.NewRequest("HEAD", dlerurl.URL, dlerurl.Body)
	if err != nil {
		return false, 0, err
	}
	resData, err := dlerurl.Client.Do(reqHead)
	if err != nil {
		return false, 0, err
	}
	defer resData.Body.Close()
	ranges := resData.Header.Get("Accept-Ranges")
	if ranges != "bytes" {
		return false, resData.ContentLength, nil
	}
	return true, resData.ContentLength, nil
}

// threadingDownloadStart 下载所有区块
func (dlerurl *DlerURL) threadingDownloadStart(file *os.File) error {
	if dlerurl.IsBar {
		fmt.Println("多线程下载")
	}
	defer file.Close()
	//创建维护线程数
	taskChannel := make(chan int, dlerurl.ThreadingCount)
	resChannel := make(chan error, 1)
	dlerClose := false
	var startSize int64 = 0
	bar := &dlerBar{IsShow: dlerurl.IsBar}
	bar.dlerBarStart(*dlerurl.ContentLength)

	defer func() {
		dlerClose = true
		close(taskChannel)
		close(resChannel)
		if bar.dlerBarIsStarted() {
			bar.dlerBarFinish()
		}
		bar = nil
	}()
	go func() {
		defer func() { recover() }()
		for startSize < *dlerurl.ContentLength {
			//申请任务
			taskChannel <- 0
			end := startSize + dlerurl.ThreadingDownloadSize
			go dlerurl.downloadSlice(resChannel, taskChannel, file, startSize, end, bar, &dlerClose, 0)
			startSize = end + 1
		}
	}()

	timeOut := time.After(dlerurl.TimeOut)
	for {
		time.Sleep(time.Millisecond * 500)
		select {
		case err := <-resChannel:
			bar.dlerBarFinish()
			return err
		case <-timeOut:
			bar.dlerBarFinish()
			return errors.New("下载超时")
		default:
			if *dlerurl.Schedule >= *dlerurl.ContentLength {
				bar.dlerBarFinish()
				if dlerurl.IsBar {
					fmt.Println("download succeed!")
				}
				// 下载完成后的回调
				err := dlerurl.downloadOver(file)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
}

// downloadSlice 下载一个区块
func (dlerurl *DlerURL) downloadSlice(resChannel chan<- error, taskChannel <-chan int, file *os.File, start, end int64, bar *dlerBar, dlerClose *bool, errorCount int64) {
	defer func() { <-taskChannel }()
	defer func() { recover() }()
	reqFile, err := dlerurl.getSliceFile(start, end, 0)
	if err != nil {
		if errorCount > 3 || *dlerClose {
			resChannel <- err
			return
		}
		// 容错
		dlerurl.downloadSlice(resChannel, taskChannel, file, start, end, bar, dlerClose, errorCount+1)
		return
	}
	resData, err := dlerurl.Client.Do(reqFile)
	if err != nil {
		if errorCount > 3 || *dlerClose {
			resChannel <- err
			return
		}
		// 容错
		dlerurl.downloadSlice(resChannel, taskChannel, file, start, end, bar, dlerClose, errorCount+1)
		return
	}
	defer resData.Body.Close()
	bytes, err := ioutil.ReadAll(resData.Body)
	if err != nil || *dlerClose {
		resChannel <- err
		return
	}
	_, err = file.WriteAt(bytes, start)
	if err != nil || *dlerClose {
		resChannel <- err
		return
	}
	batesCount := len(bytes)
	bar.dlerBarAdd(int64(batesCount))
	atomic.AddInt64(dlerurl.Schedule, int64(batesCount))
}

// getSliceFile 请求区块数据
func (dlerurl *DlerURL) getSliceFile(start, end, errorCount int64) (*http.Request, error) {
	reqFile, err := http.NewRequest(dlerurl.Method, dlerurl.URL, dlerurl.Body)
	if err != nil {
		if errorCount >= 3 {
			return reqFile, err
		}
		return dlerurl.getSliceFile(start, end, errorCount+1)
	}
	dlerurl.InitHeaders(reqFile)
	bytes := fmt.Sprintf("bytes=%s-%s", strconv.FormatInt(start, 10), strconv.FormatInt(end, 10))
	reqFile.Header.Set("Range", bytes)
	return reqFile, err
}

// dlerBar 简单封装进度条
type dlerBar struct {
	IsShow bool            // 是否显示
	DlerPb *pb.ProgressBar // pb对象
}

func (dbr *dlerBar) dlerBarStart(total int64) {
	if !dbr.IsShow {
		return
	}
	full := pb.Full
	dbr.DlerPb = pb.New64(total).SetTemplate(full).Start()
	dbr.DlerPb.Set(pb.Bytes, true)
}

func (dbr *dlerBar) dlerBarAdd(value int64) {
	if !dbr.IsShow {
		return
	}
	dbr.DlerPb.Add64(value)
}

func (dbr *dlerBar) dlerBarFinish() {
	if !dbr.IsShow {
		return
	}
	dbr.DlerPb.Finish()
}
func (dbr *dlerBar) dlerBarIsStarted() bool {
	if !dbr.IsShow {
		return false
	}
	return dbr.DlerPb.IsStarted()
}
