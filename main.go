package main

import (
	"bufio"
	"errors"
	"fmt"
	pf "gitee.com/rock_rabbit/goannie/platforms"
	"github.com/fatih/color"
	"github.com/garyburd/redigo/redis"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var goannieVersion = "v0.0.10"
var goannieUpdateTime = "2020-08-28"
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
	Name       string                                    // 匹配名称
	Info       string                                    // 简介
	UrlRegexps []*regexp.Regexp                          // URL匹配表
	Run        func(pf.RunType, map[string]string) error // 执行任务
}

// 平台结构体
type Platform struct {
	Name          string      // 平台名称
	UrlRegexps    []UrlRegexp // URL匹配表
	CookieFile    string
	DefaultCookie string
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
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("|-----------------\t ")
	color.Unset()
	color.Set(color.FgHiMagenta, color.Bold)
	fmt.Printf(pf.Name)
	color.Unset()
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" \t-----------------|\n")
	color.Unset()
	for _, item := range pf.UrlRegexps {
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
	// 初始化支持平台https://v.qq.com/vplus/8e7410558116f585a0b6c87bb849e22c#uin=8e7410558116f585a0b6c87bb849e22c
	platformList = []Platform{
		{
			"腾讯视频",
			[]UrlRegexp{
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
			"./tengxun.txt",
			"",
		}, {
			"火锅视频",
			[]UrlRegexp{
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
			"./tengxun.txt",
			"",
		},
		{
			"爱奇艺视频",
			[]UrlRegexp{
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
			"./iqiyi.txt",
			"",
		},
		{
			"西瓜视频",
			[]UrlRegexp{
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
			"./xigua.txt",
			"",
		},
		{
			"好看视频",
			[]UrlRegexp{
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
			"./haokan.txt",
			"",
		},
		{
			"哔哩哔哩",
			[]UrlRegexp{
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
			"./bilibili.txt",
			"Hm_lvt_8a6e55dbd2870f0f5bc9194cddf32a02=1585388744; _uuid=FEA617EA-B2C3-B04F-885E-46ABFAE4F7D805867infoc; buvid3=2811A771-C069-4E94-B8FA-F782D325D66D53928infoc; sid=k0i58z3h; CURRENT_FNVAL=16; LIVE_BUVID=AUTO4815859721196971; rpdid=|(J~|~mR~u~R0J'ul)lm~kuJl; DedeUserID=73471687; DedeUserID__ckMd5=0faa2eb95831029d; SESSDATA=56f9a3c2%2C1601524186%2C0ae0b*41; bili_jct=5f7ce56db63cc720b6ceffebd133e796; CURRENT_QUALITY=120; bp_video_offset_73471687=419852075596227754; PVID=2",
		},
	}
}

func main() {
	// 设置环境变量
	if err := pf.SetGoannieEnv(); err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
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
	// 检查Aria2
	if err = pf.GetAria2(); err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
	// 检查 redis
	if err = pf.GetRedis(); err != nil {
		printErrInfo(err.Error())
		exitInfo()
	}
	// 启动 redis
	cmd := exec.Command("redis-server",pf.RedisConfFile)
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
	// 是否去重
	var isDeWeight string
	isDeWeightBool := true
	err = getIsDeWeight(&isDeWeight)
	if isDeWeight == "no" {
		isDeWeightBool = false
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
		err = pf.AnnieDownload(url, savePath, "", "")
		if err != nil {
			printErrInfo(err.Error())
		}
		goto GETURL
	}
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("平台：%s  子任务：%s\n", platform.Name, subtask.Name)
	color.Unset()
	runType := pf.RunType{
		Url:           url,
		SavePath:      savePath,
		CookieFile:    platform.CookieFile,
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

// 是否去重
func getIsDeWeight(isDeWeight *string) error {
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("$ 是否去重？ yes or no (yse)：")
	color.Unset()
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	*isDeWeight = strings.Replace(string(data), "\n", "", -1)
	return nil
}

// 打印欢迎语
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
		fmt.Printf("cookie 设置：在goannie.exe同级目录中新建 %s 写入cookie=xxx;xxx=xxx;格式即可。\n", item.CookieFile)
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
