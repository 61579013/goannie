package binary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitee.com/rock_rabbit/goannie/config"
	"gitee.com/rock_rabbit/goannie/downloader"
	"gitee.com/rock_rabbit/goannie/utils"
)

/* updateBinary.go
 * 更新二进制文件
 */

// VersionJSON binary_version.json格式
type VersionJSON struct {
	Annie struct {
		Version string `json:"version"`
		FileURL string `json:"fileURL"`
	} `json:"annie"`
	Ffmpeg struct {
		Version string `json:"version"`
		FileURL string `json:"fileURL"`
	} `json:"ffmpeg"`
	Aria2 struct {
		Version string `json:"version"`
		FileURL string `json:"fileURL"`
	} `json:"aria2"`
	Redis struct {
		Version string `json:"version"`
		FileURL string `json:"fileURL"`
	} `json:"redis"`
	RedisConf struct {
		Version string `json:"version"`
		FileURL string `json:"fileURL"`
	} `json:"redisConf"`
}

// Init 初始化VersionJSON
func (vj *VersionJSON) Init() {
	if vj.Annie.Version == "" {
		vj.Annie.Version = "0.0.0"
	}
	if vj.Ffmpeg.Version == "" {
		vj.Ffmpeg.Version = "0.0.0"
	}
	if vj.Aria2.Version == "" {
		vj.Aria2.Version = "0.0.0"
	}
	if vj.Redis.Version == "" {
		vj.Redis.Version = "0.0.0"
	}
	if vj.RedisConf.Version == "" {
		vj.RedisConf.Version = "0.0.0"
	}
}

var (
	// UpdateNetworkJSONFile 更新文件网络地址
	UpdateNetworkJSONFile = config.GetString("binary.UpdateNetworkJSONFile")
	// UpdateLocalJSONFile 更新文件本地地址
	UpdateLocalJSONFile = fmt.Sprintf("%s\\binary_version.json", config.AppBinPath)
	// UpdateLockFile 更新锁文件地址
	UpdateLockFile = fmt.Sprintf("%s\\update.lock", config.AppBinPath)
	// UpdateLockTimeOut 更新锁超时时间
	UpdateLockTimeOut = config.GetInt("binary.UpdateLockTimeOut")
)

// GetVersionJSON 获取 binary_version.json 内容
func GetVersionJSON() (*VersionJSON, error) {
	f, err := os.Open(UpdateLocalJSONFile)
	if err != nil {
		return &VersionJSON{}, err
	}
	defer f.Close()
	byteData, err := ioutil.ReadAll(f)
	if err != nil {
		return &VersionJSON{}, err
	}
	var reqData VersionJSON
	if err = json.Unmarshal(byteData, &reqData); err != nil {
		return &VersionJSON{}, nil
	}
	return &reqData, nil
}

// GetNetworkJSONFile 获取网络更新文件内容
func GetNetworkJSONFile() (*VersionJSON, error) {
	request, err := http.NewRequest("GET", UpdateNetworkJSONFile, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{Timeout: time.Second * 20}
	data, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()

	byteData, err := ioutil.ReadAll(data.Body)
	var reqData VersionJSON

	if err = json.Unmarshal(byteData, &reqData); err != nil {
		return nil, err
	}
	return &reqData, nil
}

// IsUpdate 是否需要更新二进制文件，并返回需要更新的文件
func IsUpdate() (*VersionJSON, error) {
	// 网络文件 与 本地文件对比
	localFile, err := GetVersionJSON()
	if err != nil {
		// 如果文件存在
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	localFile.Init()

	networkFile, err := GetNetworkJSONFile()
	if err != nil {
		return nil, err
	}
	var retData VersionJSON
	// Annie检查
	AnnieUpdate := compareStrVer(localFile.Annie.Version, networkFile.Annie.Version)
	if AnnieUpdate == VersionSmall {
		// 检查到需要更新
		retData.Annie = networkFile.Annie
	}
	// ffmpeg检查
	FfmpegUpdate := compareStrVer(localFile.Ffmpeg.Version, networkFile.Ffmpeg.Version)
	if FfmpegUpdate == VersionSmall {
		// 检查到需要更新
		retData.Ffmpeg = networkFile.Ffmpeg
	}
	// aria2检查
	Aria2Update := compareStrVer(localFile.Aria2.Version, networkFile.Aria2.Version)
	if Aria2Update == VersionSmall {
		// 检查到需要更新
		retData.Aria2 = networkFile.Aria2
	}
	// redis检查
	RedisUpdate := compareStrVer(localFile.Redis.Version, networkFile.Redis.Version)
	if RedisUpdate == VersionSmall {
		// 检查到需要更新
		retData.Redis = networkFile.Redis
	}
	// redisConf检查
	RedisConfUpdate := compareStrVer(localFile.RedisConf.Version, networkFile.RedisConf.Version)
	if RedisConfUpdate == VersionSmall {
		// 检查到需要更新
		retData.RedisConf = networkFile.RedisConf
	}
	return &retData, nil
}

// 更新锁实现规则，update.lock 文件中写入加锁时间
// 有效的更新锁：在有效时间内，update.lock文件内有有效的时间戳

// UpdateIsLock 判断更新锁
func UpdateIsLock() error {
	// 读文件
	f, err := os.Open(UpdateLockFile)
	if err != nil {
		// 文件不存在则锁不存在
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	timeNow, err := strconv.Atoi(string(content))
	if err != nil {
		return err
	}
	if (time.Now().Unix() - int64(UpdateLockTimeOut)) > int64(timeNow) {
		// 更新锁失效
		return nil
	}
	return fmt.Errorf("更新锁存在，%d 秒后更新锁过期", int64(UpdateLockTimeOut)-(time.Now().Unix()-int64(timeNow)))
}

// UpdateLock 添加更新锁
func UpdateLock() error {
	// 判断更新锁
	if err := UpdateIsLock(); err != nil {
		return err
	}
	updateLockCount := fmt.Sprintf("%d", time.Now().Unix())
	// 删除文件
	_ = os.Remove(UpdateLockFile)
	// 检查文件夹
	dir := config.AppBinPath
	isexist, err := utils.IsExist(dir)
	if err != nil {
		return err
	}
	if !isexist {
		os.MkdirAll(dir, os.ModeDir)
	}
	// 创建文件
	f, err := os.Create(UpdateLockFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(updateLockCount); err != nil {
		return err
	}
	return nil
}

// UpdateUnLock 移除更新锁
func UpdateUnLock() error {
	err := os.Remove(UpdateLockFile)
	if err != nil {
		return nil
	}
	return nil
}

// GetVersionJSONFile 返回文件，以及内容，不存在则创建
func GetVersionJSONFile() (*os.File, *VersionJSON, error) {
	f, err := os.OpenFile(UpdateLocalJSONFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, err
	}
	byteData, err := ioutil.ReadAll(f)
	var reqData VersionJSON

	if err = json.Unmarshal(byteData, &reqData); err != nil {
		return f, &VersionJSON{}, nil
	}
	return f, &reqData, nil
}

// ExecUpdate 执行更新
func ExecUpdate(vj *VersionJSON) error {
	var err error
	// 判断更新锁
	if err = UpdateIsLock(); err != nil {
		return err
	}
	// 给更新流程加锁
	if err = UpdateLock(); err != nil {
		return err
	}
	// 移除更新锁
	defer UpdateUnLock()
	// 获取文件读写权限
	f, resData, err := GetVersionJSONFile()
	if err != nil {
		return err
	}
	defer f.Close()

	resData.Init()
	// Annie
	if vj.Annie.Version != "" && vj.Annie.FileURL != "" {
		// 执行
		utils.Infof("检查到需要更新 annie version:%s 启动下载\n", vj.Annie.Version)
		if err = DownloadFile(vj.Annie.FileURL, config.AnnieFile); err != nil {
			return err
		}
		// 更新文件
		resData.Annie = vj.Annie
		jsonData, err := json.Marshal(resData)
		if err != nil {
			return err
		}
		// 清空文件
		f.Truncate(0)
		// 写文件
		if _, err := f.WriteAt(jsonData, 0); err != nil {
			return err
		}
	}

	// ffmpeg
	if vj.Ffmpeg.Version != "" && vj.Ffmpeg.FileURL != "" {
		// 执行
		utils.Infof("检查到需要更新 ffmpeg version:%s 启动下载\n", vj.Ffmpeg.Version)
		if err = DownloadFile(vj.Ffmpeg.FileURL, config.FfmpegFile); err != nil {
			return err
		}
		// 更新文件
		resData.Ffmpeg = vj.Ffmpeg
		jsonData, err := json.Marshal(resData)
		if err != nil {
			return err
		}
		// 清空文件
		f.Truncate(0)
		// 写文件
		if _, err := f.WriteAt(jsonData, 0); err != nil {
			return err
		}
	}

	// aria2
	if vj.Aria2.Version != "" && vj.Aria2.FileURL != "" {
		// 执行
		utils.Infof("检查到需要更新 aria2 version:%s 启动下载\n", vj.Aria2.Version)
		if err = DownloadFile(vj.Aria2.FileURL, config.Aria2File); err != nil {
			return err
		}
		// 更新文件
		resData.Aria2 = vj.Aria2
		jsonData, err := json.Marshal(resData)
		if err != nil {
			return err
		}
		// 清空文件
		f.Truncate(0)
		// 写文件
		if _, err := f.WriteAt(jsonData, 0); err != nil {
			return err
		}
	}

	// redis
	if vj.Redis.Version != "" && vj.Redis.FileURL != "" {
		// 执行
		utils.Infof("检查到需要更新 redis version:%s 启动下载\n", vj.Redis.Version)
		if err = DownloadFile(vj.Redis.FileURL, config.RedisFile); err != nil {
			return err
		}
		// 更新文件
		resData.Redis = vj.Redis
		jsonData, err := json.Marshal(resData)
		if err != nil {
			return err
		}
		// 清空文件
		f.Truncate(0)
		// 写文件
		if _, err := f.WriteAt(jsonData, 0); err != nil {
			return err
		}
	}

	// redisConf
	if vj.RedisConf.Version != "" && vj.RedisConf.FileURL != "" {
		// 执行
		utils.Infof("检查到需要更新 redisConf version:%s 启动下载\n", vj.RedisConf.Version)
		if err = DownloadFile(vj.RedisConf.FileURL, config.RedisConfFile); err != nil {
			return err
		}
		// 更新文件
		resData.RedisConf = vj.RedisConf
		jsonData, err := json.Marshal(resData)
		if err != nil {
			return err
		}
		// 清空文件
		f.Truncate(0)
		// 写文件
		if _, err := f.WriteAt(jsonData, 0); err != nil {
			return err
		}
	}
	return nil
}

// DownloadFile 文件下载，存在删除
func DownloadFile(filURL, outPath string) error {
	// 检查文件夹
	dir := config.AppBinPath
	isexist, err := utils.IsExist(dir)
	if err != nil {
		return err
	}
	if !isexist {
		os.MkdirAll(dir, os.ModeDir)
	}
	_ = os.Remove(outPath)
	dir, file := filepath.Split(outPath)
	if err = downloader.New(filURL, dir).SetOutputName(file).Run(); err != nil {
		return err
	}
	return nil
}

// Update 检查并执行更新
func Update() error {
	updateJSON, err := IsUpdate()
	if err != nil {
		return err
	}
	if err = ExecUpdate(updateJSON); err != nil {
		return err
	}
	return nil
}
