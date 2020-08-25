package main

import (
	"bufio"
	"errors"
	"fmt"
	pf "gitee.com/rock_rabbit/goannie/platforms"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
)

var goannieVersion = "v0.0.05"
var goannieUpdateTime = "2020-08-25"
var goannieTitle = `
                                        __           
   __     ___      __      ___     ___ /\_\     __   
 /'_ ` + "`" + `\  / __` + "`" + `\  /'__` + "`" + `\  /' _ ` + "`" + `\ /' _ ` + "`" + `\/\ \  /'__` + "`" + `\
/\ \L\ \/\ \L\ \/\ \L\.\_/\ \/\ \/\ \/\ \ \ \/\  __/
\ \____ \ \____/\ \__/.\_\ \_\ \_\ \_\ \_\ \_\ \____\
 \/___L\ \/___/  \/__/\/_/\/_/\/_/\/_/\/_/\/_/\/____/
   /\____/
   \_/__/`

// 平台子任务结构体
type UrlRegexp struct {
	Name       string                 // 匹配名称
	Info       string                 // 简介
	UrlRegexps []*regexp.Regexp       // URL匹配表
	Run        func(pf.RunType) error // 执行任务
}

// 平台结构体
type Platform struct {
	Name       string      // 平台名称
	UrlRegexps []UrlRegexp // URL匹配表
	CookieFile string
}

// 匹配url
func (pf Platform) isUrl(url string) (bool, UrlRegexp) {
	for _, item := range pf.UrlRegexps {
		for _, must := range item.UrlRegexps {
			if must.MatchString(url) {
				return true, item
			}
		}
	}
	return false, UrlRegexp{}
}

// 打印信息
func (pf Platform) printInfo() {
	fmt.Printf("|----------------- %s -----------------|\n", pf.Name)
	for _, item := range pf.UrlRegexps {
		fmt.Printf("task: %s\tinfo: %s\n", item.Name, item.Info)
	}
}

var platformList []Platform

func init() {
	// 初始化支持平台
	platformList = []Platform{
		{
			"腾讯视频",
			[]UrlRegexp{
				{
					"detail",
					"腾讯剧集页 https://v.qq.com/detail/5/52852.html",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://v\.qq\.com/detail/\d+/\d+\.html.*?$`),
					},
					pf.RunTxDetail,
				},
			},
			"./tengxun.txt",
		},
		{
			"爱奇艺视频",
			[]UrlRegexp{
				{
					"detail",
					"爱奇艺剧集页 https://www.iqiyi.com/a_19rrht2ok5.html",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.iqiyi\.com/a_\w+\.html.*?$`),
					},
					pf.RunIqyDetail,
				},
			},
			"./iqiyi.txt",
		},
		{
			"西瓜视频",
			[]UrlRegexp{
				{
					"one",
					"单视频 https://www.ixigua.com/6832194590221533707",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/\d+.*?$`),
						regexp.MustCompile(`^(http|https)://m\.ixigua\.com/\d+.*?$`),
						regexp.MustCompile(`^(http|https)://m\.ixigua\.com/video/\d+.*?$`),
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/.*?\?id=\d+.*?$`),
						regexp.MustCompile(`^(http|https)://toutiao\.com/group/\d+/.*?$`),
					},
					pf.RunXgOne,
				},{
					"userList",
					"TA的视频 https://www.ixigua.com/home/85383446500/video/",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/home/\d+/video.*?$`),
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/home/\d+/($|\?.*?$)`),
					},
					pf.RunXgUserList,
				},
			},
			"./xigua.txt",
		},
	}
}

func main() {
	// 设置环境变量
	if err := pf.SetGoannieEnv(); err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
	printHello()
	// 检查 annie
	err := pf.GetAnnie()
	if err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
	// 检查ffmpeg
	if err = pf.GetFfmpeg(); err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
GETSAVEPATH:
	var savePath string
	err = getSavePath(&savePath)
	if err != nil {
		printErrInfo(err.Error())
		goto GETSAVEPATH
	}
	isSavePath, err := isDir(savePath)
	if !isSavePath {
		if err != nil {
			printErrInfo(err.Error())
			goto GETSAVEPATH
		} else {
			printErrInfo("文件夹不存在")
			goto GETSAVEPATH
		}
	}
GETURL:
	var url string
	err = getUrl(&url)
	if err != nil {
		printErrInfo(err.Error())
		goto GETURL
	}
	platform, subtask, err := getUrlPlatform(url)
	if err != nil {
		// 尝试直接使用 annie
		err = pf.AnnieDownload(url, savePath, "")
		if err != nil {
			printErrInfo(err.Error())
		}
		goto GETURL
	}
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("平台：%s  子任务：%s\n", platform.Name, subtask.Name)
	color.Unset()
	runType := pf.RunType{
		Url:        url,
		SavePath:   savePath,
		CookieFile: platform.CookieFile,
	}
	err = subtask.Run(runType)
	if err != nil {
		printErrInfo(err.Error())
		goto GETURL
	}
	goto GETURL
}

// 获取url平台
func getUrlPlatform(url string) (Platform, UrlRegexp, error) {
	for _, item := range platformList {
		isUrl, subtask := item.isUrl(url)
		if isUrl {
			return item, subtask, nil
		}
	}
	return Platform{}, UrlRegexp{}, errors.New("不支持这个链接")
}

// 获取URL
func getUrl(url *string) error {
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("$ 请输入URL：")
	color.Unset()
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	*url = strings.Replace(string(data), "\n", "", -1)
	return nil
}

// 文件夹是否存在
func isDir(path string) (bool, error) {
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

// 获取保存路径
func getSavePath(savePath *string) error {
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("$ 请输入保存路径：")
	color.Unset()
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	*savePath = strings.Replace(string(data), "\n", "", -1)
	return nil
}

// 打印欢迎语
func printHello() {
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	fmt.Println(goannieTitle)
	color.Set(color.FgMagenta, color.Bold)
	fmt.Printf("	version: %s	updateTime: %s\n\n", goannieVersion, goannieUpdateTime)
	color.Set(color.FgBlue, color.Bold)
	fmt.Println("支持平台")
	for _, item := range platformList {
		item.printInfo()
		fmt.Printf("cookie 设置：在goannie.exe同级目录中新建 %s 写入cookie=xxx;xxx=xxx;格式即可。\n", item.CookieFile)
		if item.Name == "腾讯视频"{
			fmt.Println("ccode 和 ckey 设置：在goannie.exe同级目录中新建 ccode.txt 和 ckey.txt 写入其中即可。")
		}
		fmt.Println("")
	}
}

// 打印错误信息
func printErrInfo(errInfo string) {
	color.Set(color.FgRed, color.Bold)
	defer color.Unset()
	fmt.Println("错误信息：" + errInfo)
}

func exitInfo() {
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	fmt.Printf("$ 回车退出：\n")
	var s string
	_, _ = fmt.Scanln(&s)
	os.Exit(1)
}
