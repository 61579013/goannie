package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/rock_rabbit/goannie/godler"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// 控制台输入
func GetCmdDataString(Info string, resData *string) error {
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("$ %s：", Info)
	color.Unset()
	_, err := fmt.Scanln(resData)
	if err != nil {
		return err
	}
	return nil
}

// 使用annie批量下载
func AnnieDownloadAll(urlList []map[string]string, runType RunType, pt string) {
	for _, item := range urlList {
		err := AnnieDownload(item["url"], runType.SavePath, runType.CookieFile, runType.DefaultCookie)
		if err != nil {
			PrintErrInfo(err.Error())
		} else {
			AddVideoID(pt, item["vid"])
		}
	}
}

// 使用annie下载
func AnnieDownload(url, savePath, cookiePath, DefaultCookie string) error {
	arg := []string{}
	ccode := GetTxtContent("./ccode.txt")
	ckey := GetTxtContent("./ckey.txt")
	if ccode != "" {
		arg = append(arg, "-ccode", ccode)
	}
	if ckey != "" {
		arg = append(arg, "-ckey", ckey)
	}
	// COOKE 设定
	onCookie := cookiePath
	if GetTxtContent(cookiePath) == "" {
		onCookie = DefaultCookie
	}
	arg = append(arg, []string{"-retry", "3", "-c", onCookie, "-o", savePath, url}...)
	cmd := exec.Command("annie", arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func Aria2Download(url, savePath, saveFile, cookiePath string, maxConnectionPerServer int) error {
	// -c 断点续传
	// -m 失败重试次数。默认值为：5
	// --retry-wait=<SEC> 失败重试间隔时间（单位：秒），默认值为：0
	// -o 命名下载文件
	// -x 设置每个下载最大的连接数
	cmd := exec.Command("aria2c",
		"-c", "-m", "5",
		"--retry-wait=10", "-x",
		fmt.Sprintf("%d", maxConnectionPerServer), `--header="cookie: `+GetTxtContent(cookiePath)+`"`,
		"-d", savePath,
		"-o", saveFile,
		"--console-log-level=warn",
		url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 打印错误信息
func PrintErrInfo(errInfo string) {
	color.Set(color.FgRed, color.Bold)
	defer color.Unset()
	fmt.Println("错误信息：" + errInfo)
}

// 打印日常信息
func PrintInfo(info string) {
	color.Set(color.FgBlue, color.Bold)
	defer color.Unset()
	fmt.Println(info)
}

func PrintInfof(info string) {
	color.Set(color.FgBlue, color.Bold)
	defer color.Unset()
	fmt.Printf(info)
}

// 请求 annie
func GetAnnie() error {
	IsAnniePath, err := IsExist(AppBinPath)
	if err != nil {
		return err
	}
	if !IsAnniePath {
		if err = os.MkdirAll(AppBinPath, os.ModePerm); err != nil {
			return nil
		}
	}
	IsAnnieFile, err := IsExist(AnnieFile)
	if err != nil {
		return err
	}
	if !IsAnnieFile {
		PrintInfo("检查到该机器没有下载 annie 启动下载")
		dlerurl := godler.DlerUrl{
			Url:      "http://image.68wu.cn/annie/annie_0.10.3_64.exe",
			SavePath: AnnieFile,
			IsBar:    true,
		}
		err := dlerurl.Download()
		if err != nil {
			return err
		}
	}
	return nil
}

// 请求 ffmpeg
func GetFfmpeg() error {
	IsAnniePath, err := IsExist(AppBinPath)
	if err != nil {
		return err
	}
	if !IsAnniePath {
		if err = os.MkdirAll(AppBinPath, os.ModePerm); err != nil {
			return nil
		}
	}
	IsFfmpegFile, err := IsExist(FfmpegFile)
	if err != nil {
		return err
	}
	if !IsFfmpegFile {
		PrintInfo("检查到该机器没有下载 ffmpeg 启动下载")
		dlerurl := godler.DlerUrl{
			Url:      "http://image.68wu.cn/ffmpeg/ffmpeg.exe",
			SavePath: FfmpegFile,
			IsBar:    true,
		}
		err := dlerurl.Download()
		if err != nil {
			return err
		}
	}
	return nil
}

// 请求 aria2
func GetAria2() error {
	IsAnniePath, err := IsExist(AppBinPath)
	if err != nil {
		return err
	}
	if !IsAnniePath {
		if err = os.MkdirAll(AppBinPath, os.ModePerm); err != nil {
			return nil
		}
	}
	IsAria2File, err := IsExist(Aria2File)
	if err != nil {
		return err
	}
	if !IsAria2File {
		PrintInfo("检查到该机器没有下载 aria2 启动下载")
		dlerurl := godler.DlerUrl{
			Url:      "http://image.68wu.cn/aria2/aria2c.exe",
			SavePath: Aria2File,
			IsBar:    true,
		}
		err := dlerurl.Download()
		if err != nil {
			return err
		}
	}
	return nil
}

// 文件夹或文件是否存在
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// 获取txt内容
func GetTxtContent(path string) string {
	IsCookie, _ := IsExist(path)
	if !IsCookie {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	cotent, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(cotent)
}

// 设置环境变量
func SetGoannieEnv() error {
	if strings.Index(os.Getenv("PATH"), AppBinPath) == -1 {
		// 添加环境变量
		err := os.Setenv("PATH", fmt.Sprintf("%s;%s", os.Getenv("PATH"), AppBinPath))
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取跳转真实地址
func GetRealUrl(url string) (string, error) {
	newClient := Client
	newClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("user-agent", UserAgentPc)
	resP, err := newClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resP.Body.Close()
	if resP.StatusCode == 301 || resP.StatusCode == 302 {
		return resP.Header.Get("location"), nil
	}
	return "", errors.New("获取跳转真实地址失败")
}

//对比库是否存在视频ID
func IsVideoID(pt, vid string) bool {
	resData, err := GetVideoID(pt)
	if err != nil {
		return false
	}
	for _, i := range *resData {
		if i == vid {
			return true
		}
	}
	return false
}

// 存储已下载的视频ID
func AddVideoID(pt, vid string) (bool, error) {
	resData, err := GetVideoID(pt)
	if err != nil {
		return false, err
	}
	for _, i := range *resData {
		if i == vid {
			return true, nil
		}
	}
	*resData = append(*resData, vid)
	AppDataFile := fmt.Sprintf("%s\\%s.json", AppDataPath, pt)
	f, err := os.OpenFile(AppDataFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}
	defer f.Close()
	content, err := json.Marshal(*resData)
	if err != nil {
		return false, err
	}
	_, err = f.Write(content)
	if err != nil {
		return false, err
	}
	return false, nil
}

// 获取已下载的视频ID
func GetVideoID(pt string) (*[]string, error) {
	AppDataFile := fmt.Sprintf("%s\\%s.json", AppDataPath, pt)
	IsAppDataPath, err := IsExist(AppDataPath)
	if err != nil {
		return nil, err
	}
	if !IsAppDataPath {
		if err = os.MkdirAll(AppDataPath, os.ModePerm); err != nil {
			return nil, nil
		}
	}
	if isDataFile, _ := IsExist(AppDataFile); !isDataFile {
		// 创建
		err := CreactVideoID(pt)
		if err != nil {
			return nil, err
		}
		return &[]string{}, nil
	}
	file, err := os.Open(AppDataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	resData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var data []string
	err = json.Unmarshal(resData, &data)
	if err != nil {
		// 这里 json 格式错误 覆盖写
		file.Close()
		err := CreactVideoID(pt)
		if err != nil {
			return nil, err
		}
		return &[]string{}, nil
	}
	return &data, nil
}

func CreactVideoID(pt string) error {
	// 创建
	AppDataFile := fmt.Sprintf("%s\\%s.json", AppDataPath, pt)
	file, err := os.Create(AppDataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte("[]"))
	if err != nil {
		return err
	}
	return nil
}

func PrintVideoIDCount() {
	ptList := []map[string]string{
		{
			"name":  "腾讯视频",
			"pt":    "tengxun",
			"count": "",
		}, {
			"name":  "爱奇艺视频",
			"pt":    "iqiyi",
			"count": "",
		}, {
			"name":  "好看视频",
			"pt":    "haokan",
			"count": "",
		}, {
			"name":  "哔哩哔哩",
			"pt":    "bilibili",
			"count": "",
		}, {
			"name":  "西瓜视频",
			"pt":    "xigua",
			"count": "",
		},
	}
	for idx,item := range ptList{
		data,err := GetVideoID(item["pt"])
		if err != nil{
			ptList[idx]["count"] = "--"
		}else{
			ptList[idx]["count"] = fmt.Sprintf("%d",len(*data))
		}
		fmt.Printf("%s：%s  ",item["name"],ptList[idx]["count"])
	}
	fmt.Println("")
}
