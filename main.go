package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitee.com/rock_rabbit/goannie/extractors"
	"gitee.com/rock_rabbit/goannie/reptiles"
	"gitee.com/rock_rabbit/goannie/request"
	"gitee.com/rock_rabbit/goannie/storage"

	"gitee.com/rock_rabbit/goannie/binary"
	"gitee.com/rock_rabbit/goannie/config"
	extractorsTypes "gitee.com/rock_rabbit/goannie/extractors/types"
	reptilesTypes "gitee.com/rock_rabbit/goannie/reptiles/types"
	"gitee.com/rock_rabbit/goannie/utils"
	"github.com/fatih/color"
	"github.com/garyburd/redigo/redis"
)

func init() {
	var err error
	// 设置环境变量
	if err = utils.SetGoannieEnv(config.AppBinPath); err != nil {
		utils.ErrInfo(err.Error())
		utils.ExitInfo()
	}
	// 检查二进制文件更新
	if config.GetBool("binary.check") {
		if err := binary.Update(); err != nil {
			utils.ErrInfo(err.Error())
		}
	}
	// 检查Data目录
	if isDataPath, _ := utils.IsExist(config.AppDataPath); !isDataPath {
		if err := os.MkdirAll(config.AppDataPath, os.ModePerm); err != nil {
			utils.ErrInfo(err.Error())
			utils.ExitInfo()
		}
	}
	// 启动 redis
	if config.GetBool("redis.start") {
		cmd := exec.Command("cmd", "/c", "start", "/B", "redis-server", config.RedisConfFile, "--dir", config.AppDataPath)
		if err = cmd.Start(); err != nil {
			utils.ErrInfo(err.Error())
		}
	}
	// 连接 redis
	if config.GetBool("redis.dial") {
		CONN, err = redis.Dial(config.GetString("redis.network"), config.GetString("redis.address"))
		if err != nil {
			utils.ErrInfo(err.Error())
			utils.ExitInfo()
		}
	}
	sayHello()
}

// CONN redis服务连接
var CONN redis.Conn

func main() {
	defer CONN.Close()
	var err error
	// 创建 redis 存储器
	verify := storage.NewRedis(CONN)
	videoIDCount(verify)
	fmt.Println("")
GETSAVEPATH:
	sayPathlist()
	var savePath string
	if savePath, err = getSavepath(); err != nil {
		utils.ErrInfo(err.Error())
		goto GETSAVEPATH
	}
GETURL:
	var url string
	if err = utils.GetStrInput("请输入URL或.txt", &url); err != nil {
		utils.ErrInfo(err.Error())
		goto GETURL
	}
	// 首先判断是否txt文件
	if filepath.Ext(url) == ".txt" {
		goto GETURL
	}
	cookiePath, defCookie := getCookiepath(url)
	cookie := setRequestOptions(cookiePath, defCookie)
	reptilesData, err := reptiles.Extract(url, reptilesTypes.Options{
		Cookie: cookie,
		Verify: verify,
	})
	if err != nil {
		utils.ErrInfo(err.Error())
		goto GETURL
	}
	for _, d := range reptilesData {
		if err = extractors.Extract(d.URL, extractorsTypes.Options{
			Cookie:   cookie,
			Verify:   verify,
			SavePath: savePath,
		}); err != nil {
			utils.ErrInfo(err.Error())
		}
	}
	goto GETURL
}

func setRequestOptions(cookie, defCookie string) string {
	cookiePath := cookie
	if cookie != "" {
		// If cookie is a file path, convert it to a string to ensure cookie is always string
		if _, fileErr := os.Stat(cookie); fileErr == nil {
			// Cookie is a file
			data, err := ioutil.ReadFile(cookie)
			if err != nil {
				utils.ErrInfo(err.Error())
				return ""
			}
			cookie = strings.TrimSpace(string(data))
		}
	}
	if cookie == cookiePath || cookie == "" {
		cookie = strings.TrimSpace(defCookie)
	}
	request.SetOptions(request.Options{
		RetryTimes: config.GetInt("app.retryTimes"),
		Cookie:     cookie,
		Refer:      config.GetString("app.refer"),
		Debug:      config.GetBool("app.debug"),
	})
	return cookie
}

func getCookiepath(u string) (string, string) {
	var domain string
	u = strings.TrimSpace(u)
	parseU, err := url.ParseRequestURI(u)
	if err != nil {
		return "", ""
	}
	if parseU.Host == "haokan.baidu.com" {
		domain = "haokan"
	} else {
		domain = utils.Domain(parseU.Host)
	}
	if _, ok := extractors.ExtractorMap[domain]; !ok {
		return "", ""
	}
	extractor := extractors.ExtractorMap[domain]
	filename := fmt.Sprintf("%s.txt", extractor.Key())
	currenpath, _ := utils.GetCurrentPath()
	cookiepath := filepath.Join(currenpath, filename)

	hiWhite := color.New(color.FgHiWhite)
	hiBlue := color.New(color.FgHiBlue)
	hiBlue.Printf("%s%s %s%s %s%s\n", hiBlue.Sprint("Name："), hiWhite.Sprint(extractor.Name()), hiBlue.Sprint("Key："), hiWhite.Sprint(extractor.Key()), hiBlue.Sprint("Cookiepath："), hiWhite.Sprint(cookiepath))
	return cookiepath, extractor.DefCookie()
}

func getSavepath() (string, error) {
	var (
		err  error
		path string
	)
	if err = utils.GetStrInput("保存路径", &path); err != nil {
		return "", err
	}
	switch path {
	case "p1", "p2", "p3", "p4", "p5":
		path = config.GetString("outpath." + path)
	}
	if config.GetBool("app.autoCreatePath") {
		os.MkdirAll(path, 0666)
	}
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
		return "", err
	}
	savePath(path)
	return path, nil
}

func savePath(s string) {
	p1 := config.GetString("outpath.p1")
	p2 := config.GetString("outpath.p2")
	p3 := config.GetString("outpath.p3")
	p4 := config.GetString("outpath.p4")
	p5 := config.GetString("outpath.p5")
	if s == p1 || s == p2 || s == p3 || s == p4 || s == p5 || s == "" {
		return
	}
	config.Set("outpath.p1", s)
	config.Set("outpath.p2", p1)
	config.Set("outpath.p3", p2)
	config.Set("outpath.p4", p3)
	config.Set("outpath.p5", p4)
	config.WriteConfig()
}

func sayHello() {
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)
	hiBlue := color.New(color.FgHiBlue)
	hiWhite := color.New(color.FgHiWhite)
	green.Println(config.TITLE)
	magenta.Printf("	版本: %s	更新时间: %s\n\n", config.VERSION, config.UPDATETIME)
	hiBlue.Printf("%s %s%s %s%s\n", green.Sprint("$"), hiBlue.Sprint("作者："), hiWhite.Sprint("rockrabbit"), hiBlue.Sprint("作者主页："), hiWhite.Sprint("https://www.68wu.com"))
	hiBlue.Printf("%s %s%s %s%s\n", green.Sprint("$"), hiBlue.Sprint("GIT仓库："), hiWhite.Sprint("https://gitee.com/rock_rabbit/goannie"), hiBlue.Sprint("开源协议："), hiWhite.Sprint("MIT"))
	fmt.Println("")
	hiBlue.Printf("……………………………… %s %s\n", hiWhite.Sprint("下载统计"), hiBlue.Sprint("………………………………"))
}

func sayPathlist() {
	hiWhite := color.New(color.FgHiWhite)
	hiBlue := color.New(color.FgHiBlue)
	hiBlue.Printf("……………………………… %s %s\n", hiWhite.Sprint("历史保存路径"), hiBlue.Sprint("………………………………"))
	p1 := config.GetString("outpath.p1")
	p2 := config.GetString("outpath.p2")
	p3 := config.GetString("outpath.p3")
	p4 := config.GetString("outpath.p4")
	p5 := config.GetString("outpath.p5")
	if p1 != "" {
		hiBlue.Printf("[p1] %s\n", hiWhite.Sprint(config.GetString("outpath.p1")))
	}
	if p2 != "" {
		hiBlue.Printf("[p2] %s\n", hiWhite.Sprint(config.GetString("outpath.p2")))
	}
	if p3 != "" {
		hiBlue.Printf("[p3] %s\n", hiWhite.Sprint(config.GetString("outpath.p3")))
	}
	if p4 != "" {
		hiBlue.Printf("[p4] %s\n", hiWhite.Sprint(config.GetString("outpath.p4")))
	}
	if p5 != "" {
		hiBlue.Printf("[p5] %s\n", hiWhite.Sprint(config.GetString("outpath.p5")))
	}
	fmt.Println("")
}

// videoIDCount 打印过滤库个数
func videoIDCount(verify storage.Storage) {
	hiWhite := color.New(color.FgHiWhite)
	hiBlue := color.New(color.FgHiBlue)
	for _, e := range extractors.ExtractorMap {
		hiBlue.Printf("%s%s  ", hiBlue.Sprintf("%s：", e.Name()), hiWhite.Sprint(verify.Count(e.Key())))
	}
	fmt.Println("")
}
