package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gitee.com/rock_rabbit/goannie/binary"
	pf "gitee.com/rock_rabbit/goannie/platforms"
	"github.com/fatih/color"
	"github.com/garyburd/redigo/redis"
)

var goannieVersion = "v0.0.18"
var goannieUpdateTime = "2020-09-29"
var goannieTitle = `
                                        __           
   __     ___      __      ___     ___ /\_\     __   
 /'_ ` + "`" + `\  / __` + "`" + `\  /'__` + "`" + `\  /' _ ` + "`" + `\ /' _ ` + "`" + `\/\ \  /'__` + "`" + `\
/\ \L\ \/\ \L\ \/\ \L\.\_/\ \/\ \/\ \/\ \ \ \/\  __/
\ \____ \ \____/\ \__/.\_\ \_\ \_\ \_\ \_\ \_\ \____\
 \/___L\ \/___/  \/__/\/_/\/_/\/_/\/_/\/_/\/_/\/____/
   /\____/
   \_/__/`

// URLRegexp 平台子任务结构体
type URLRegexp struct {
	Name       string                                    // 匹配名称
	Info       string                                    // 简介
	URLRegexps []*regexp.Regexp                          // URL匹配表
	Run        func(pf.RunType, map[string]string) error // 执行任务
}

// Platform 平台结构体
type Platform struct {
	Name          string      // 平台名称
	URLRegexps    []URLRegexp // URL匹配表
	CookieFile    string      // Cookie文件路径
	DefaultCookie string      // 默认Cookie
}

// isURL 匹配url
func (pf Platform) isURL(url string) (bool, URLRegexp) {
	for _, item := range pf.URLRegexps {
		for _, must := range item.URLRegexps {
			if must.MatchString(url) {
				return true, item
			}
		}
	}
	return false, URLRegexp{}
}

// printInfo 打印信息
func (pf Platform) printInfo() {
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("|-----------------\t ")
	color.Unset()
	color.Set(color.FgHiMagenta, color.Bold)
	fmt.Printf(pf.Name)
	color.Unset()
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" \t-----------------|\n")
	color.Unset()
	for _, item := range pf.URLRegexps {
		color.Set(color.FgBlue, color.Bold)
		fmt.Printf("task: ")
		color.Unset()
		fmt.Printf(item.Name)
		color.Set(color.FgBlue, color.Bold)
		fmt.Printf("\tinfo: ")
		color.Unset()
		fmt.Printf(item.Info + "\n")
	}
}

var platformList []Platform

func init() {
	// 初始化支持平台
	platformList = []Platform{
		{
			"抖音视频",
			[]URLRegexp{
				{
					"one",
					"单视频		https://www.iesdouyin.com/share/video/6877354382132808971",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.iesdouyin\.com/share/video/\d+`),
					},
					pf.RunDyOne,
				}, {
					"userList",
					"作者视频		https://www.iesdouyin.com/share/user/2836383897749943?sec_uid=xxxxx",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.iesdouyin\.com/share/user/\d+\?sec_uid=.*?`),
					},
					pf.RunDyUserList,
				}, {
					"shortURL",
					"短链接		https://v.douyin.com/JDq8uv7/",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://v\.douyin\.com/.*?`),
					},
					pf.RunDyShortURL,
				},
			},
			"douyin.txt",
			"",
		}, {
			"腾讯视频",
			[]URLRegexp{
				{
					"one",
					"单视频		https://v.qq.com/x/cover/mzc00200agq0com/r31376lllyf.html",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://v\.qq\.com/x/cover/.*?/.*?\.html.*?$`),
						regexp.MustCompile(`^(http|https)://v\.qq\.com/x/cover/.*?\.html.*?$`),
						regexp.MustCompile(`^(http|https)://v\.qq\.com/x/page/.*?/.*?\.html.*?$`),
						regexp.MustCompile(`^(http|https)://v\.qq\.com/x/page/.*?\.html.*?$`),
					},
					pf.RunTxOne,
				},
				{
					"detail",
					"腾讯剧集页	https://v.qq.com/detail/5/52852.html",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://v\.qq\.com/detail/\d+/\d+\.html.*?$`),
					},
					pf.RunTxDetail,
				}, {
					"userList",
					"作者视频		https://v.qq.com/s/videoplus/1790091432",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://v\.qq\.com/biu/videoplus\?vuid=\d+.*?$`),
						regexp.MustCompile(`^(http|https)://v\.qq\.com/s/videoplus/\d+.*?$`),
						regexp.MustCompile(`^(http|https)://v\.qq\.com/x/bu/h5_user_center\?vuid=\d+.*?$`),
					},
					pf.RunTxUserList,
				}, {
					"lookList",
					"看作者作品列表	look https://v.qq.com/s/videoplus/1790091432",
					[]*regexp.Regexp{
						regexp.MustCompile(`^look (http|https)://v\.qq\.com/biu/videoplus\?vuid=\d+.*?$`),
						regexp.MustCompile(`^look (http|https)://v\.qq\.com/s/videoplus/\d+.*?$`),
						regexp.MustCompile(`^look (http|https)://v\.qq\.com/x/bu/h5_user_center\?vuid=\d+.*?$`),
					},
					pf.RunLookTxUserList,
				},
			},
			"tengxun.txt",
			"",
		}, {
			"火锅视频",
			[]URLRegexp{
				{
					"userList",
					"作者视频		https://huoguo.qq.com/m/person.html?userid=18590596",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://huoguo\.qq\.com/m/person\.html\?userid=\d+.*?$`),
					},
					pf.RunHgUserList,
				}, {
					"lookList",
					"看作者作品列表	look https://huoguo.qq.com/m/person.html?userid=18590596",
					[]*regexp.Regexp{
						regexp.MustCompile(`^look (http|https)://huoguo\.qq\.com/m/person\.html\?userid=\d+.*?$`),
					},
					pf.RunLookHgUserList,
				},
			},
			"tengxun.txt",
			"",
		},
		{
			"爱奇艺视频",
			[]URLRegexp{
				{
					"one",
					"单视频		https://www.iqiyi.com/v_1fr4mggxzpo.html",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.iqiyi\.com/v_\w+\.html.*?$`),
					},
					pf.RunIqyOne,
				},
				{
					"detail",
					"爱奇艺剧集页	https://www.iqiyi.com/a_19rrht2ok5.html",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.iqiyi\.com/a_\w+\.html.*?$`),
					},
					pf.RunIqyDetail,
				},
			},
			"iqiyi.txt",
			"",
		},
		{
			"西瓜视频",
			[]URLRegexp{
				{
					"one",
					"单视频		https://www.ixigua.com/6832194590221533707",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/\d+.*?$`),
						regexp.MustCompile(`^(http|https)://m\.ixigua\.com/\d+.*?$`),
						regexp.MustCompile(`^(http|https)://m\.ixigua\.com/video/\d+.*?$`),
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/.*?\?id=\d+.*?$`),
						regexp.MustCompile(`^(http|https)://toutiao\.com/group/\d+/.*?$`),
					},
					pf.RunXgOne,
				}, {
					"userList",
					"TA的视频		https://www.ixigua.com/home/85383446500/video/",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/home/\d+/video.*?$`),
						regexp.MustCompile(`^(http|https)://www\.ixigua\.com/home/\d+($|/$|/\?.*?$|\?.*?$)`),
					},
					pf.RunXgUserList,
				}, {
					"lookList",
					"看作者作品列表	look https://www.ixigua.com/home/85383446500/video/",
					[]*regexp.Regexp{
						regexp.MustCompile(`^look (http|https)://www\.ixigua\.com/home/\d+/video.*?$`),
						regexp.MustCompile(`^look (http|https)://www\.ixigua\.com/home/\d+/($|\?.*?$)`),
					},
					pf.RunLookXgUserList,
				},
			},
			"xigua.txt",
			"",
		},
		{
			"好看视频",
			[]URLRegexp{
				{
					"one",
					"单视频		https://haokan.baidu.com/v?vid=3881011031260239591",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://haokan\.baidu\.com/v\?vid=\d+.*?$`),
					},
					pf.RunHkOne,
				}, {
					"userList",
					"作者视频		https://haokan.baidu.com/author/1649278643844524",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://haokan\.baidu\.com/author/\d+.*?$`),
					},
					pf.RunHkUserList,
				},
			},
			"haokan.txt",
			"",
		},
		{
			"哔哩哔哩",
			[]URLRegexp{
				{
					"one",
					"单视频		https://www.bilibili.com/video/BV1iK4y1e7uL",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://www\.bilibili\.com/video/\w+.*?$`),
						regexp.MustCompile(`^(http|https)://www\.bilibili\.com/bangumi/play/\w+.*?$`),
					},
					pf.RunBliOne,
				},
				{
					"userList",
					"TA的视频		https://space.bilibili.com/337312411",
					[]*regexp.Regexp{
						regexp.MustCompile(`^(http|https)://space\.bilibili\.com/\d+.*?$`),
					},
					pf.RunBliUserList,
				},
			},
			"bilibili.txt",
			"",
		},
	}
}

func main() {
	// 设置环境变量
	if err := pf.SetGoannieEnv(); err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
	// 检查二进制文件更新
	if err := binary.Update(); err != nil {
		printErrInfo(err.Error())
	}
	// 检查Data目录
	if isDataPath, _ := isDir(pf.AppDataPath); !isDataPath {
		if err := os.MkdirAll(pf.AppDataPath, os.ModePerm); err != nil {
			printErrInfo(err.Error())
			exitInfo()
		}
	}
	// 启动 redis
	cmd := exec.Command("redis-server", pf.RedisConfFile, "--dir", pf.AppDataPath)
	_ = cmd.Start()
	// 连接 redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
	defer conn.Close()
	printHello(conn)
	// v0.0.09 版本兼容
	pf.TOLeadRedis(conn)
GETSAVEPATH:
	var savePath string
	err = getInput("请输入保存路径", &savePath)
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
	// 是否过滤重复下载？
	var isDeWeight string
	isDeWeightBool := true
	err = getInput("是否过滤重复下载？ yes or no (yse)", &isDeWeight)
	if isDeWeight == "no" {
		isDeWeightBool = false
	}
GETURL:
	var url string
	err = getInput("请输入URL", &url)
	if err != nil {
		printErrInfo(err.Error())
		goto GETURL
	}
	platform, subtask, err := getURLPlatform(url)
	if err != nil {
		// 尝试直接使用 annie
		err = pf.AnnieDownload(url, savePath, "", "")
		if err != nil {
			printErrInfo(err.Error())
		}
		goto GETURL
	}
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("平台：%s  子任务：%s\n", platform.Name, subtask.Name)
	color.Unset()
	currenpath, _ := getCurrentPath()
	runType := pf.RunType{
		URL:           url,
		SavePath:      savePath,
		CookieFile:    fmt.Sprintf("%s%s", currenpath, platform.CookieFile),
		DefaultCookie: platform.DefaultCookie,
		IsDeWeight:    isDeWeightBool,
		RedisConn:     conn,
	}
	err = subtask.Run(runType, map[string]string{})
	if err != nil {
		printErrInfo(err.Error())
		goto GETURL
	}
	goto GETURL
}

// getURLPlatform 获取url平台
func getURLPlatform(url string) (Platform, URLRegexp, error) {
	for _, item := range platformList {
		isURL, subtask := item.isURL(url)
		if isURL {
			return item, subtask, nil
		}
	}
	return Platform{}, URLRegexp{}, errors.New("不支持这个链接")
}

// isDir 文件夹是否存在
func isDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// getInput 获取控制台输入
func getInput(info string, outStr *string) error {
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("$ %s：", info)
	color.Unset()
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	*outStr = strings.Replace(string(data), "\n", "", -1)
	return nil
}

// printHello 打印欢迎语
func printHello(conn redis.Conn) {
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	fmt.Println(goannieTitle)
	color.Set(color.FgMagenta, color.Bold)
	fmt.Printf("	version: %s	updateTime: %s\n\n", goannieVersion, goannieUpdateTime)
	color.Set(color.FgHiBlue, color.Bold)
	fmt.Println("支持平台")
	for _, item := range platformList {
		item.printInfo()
		color.Set(color.FgHiBlack, color.Bold)
		fmt.Printf("cookie 设置：在goannie.exe同级目录中新建 %s 写入name=value;name=value....格式即可。\n", item.CookieFile)
		if item.Name == "腾讯视频" {
			fmt.Println("ccode 和 ckey 设置：在goannie.exe同级目录中新建 ccode.txt 和 ckey.txt 写入其中即可。")
		}
		color.Unset()
		fmt.Println("")
	}
	color.Set(color.FgHiBlue, color.Bold)
	fmt.Println("下载统计")
	color.Unset()
	pf.PrintVideoIDCount(conn)
	fmt.Println("")
}

// printErrInfo 打印错误信息
func printErrInfo(errInfo string) {
	color.Set(color.FgRed, color.Bold)
	defer color.Unset()
	fmt.Println("错误信息：" + errInfo)
}

// exitInfo 结束程序
func exitInfo() {
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	fmt.Printf("$ 回车退出：")
	var s string
	_, _ = fmt.Scanln(&s)
	os.Exit(1)
}

// getCurrentPath 获取程序所在目录
func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\". `)
	}
	return string(path[0 : i+1]), nil
}
