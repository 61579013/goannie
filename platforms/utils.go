package platforms

import (
	"fmt"
	"gitee.com/rock_rabbit/goannie/godler"
	"github.com/fatih/color"
	"io/ioutil"
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
func AnnieDownloadAll(urlList []map[string]string, runType RunType) {
	for _, item := range urlList {
		err := AnnieDownload(item["url"], runType.SavePath, runType.CookieFile)
		if err != nil{
			PrintErrInfo(err.Error())
		}
	}
}

// 使用annie下载
func AnnieDownload(url, savePath, cookiePath string) error {
	arg := []string{}
	ccode := GetTxtContent("./ccode.txt")
	ckey := GetTxtContent("./ckey.txt")
	if ccode != "" {
		arg = append(arg, "-ccode", ccode)
	}
	if ckey != "" {
		arg = append(arg, "-ckey", ckey)
	}
	arg = append(arg, []string{"-retry","3","-c", cookiePath, "-o", savePath, url}...)
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
	if strings.Index(os.Getenv("PATH"),AppBinPath) == -1{
		// 添加环境变量
		err := os.Setenv("PATH",fmt.Sprintf("%s;%s",os.Getenv("PATH"),AppBinPath))
		if err != nil{
			return err
		}
	}
	return nil
}
