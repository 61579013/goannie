package main

import (
	"fmt"
	"os"
	"os/exec"

	"gitee.com/rock_rabbit/goannie/extractors"

	"gitee.com/rock_rabbit/goannie/binary"
	"gitee.com/rock_rabbit/goannie/config"
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
	// s := storage.NewRedis(CONN)
GETSAVEPATH:
	sayPathlist()
	var savePath string
	if savePath, err = getSavepath(); err != nil {
		utils.ErrInfo(err.Error())
		goto GETSAVEPATH
	}
	fmt.Println(savePath)
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
	// 在这里直接判断txt文件
	// ...
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
	videoIDCount()
	fmt.Println("")
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
func videoIDCount() {
	hiWhite := color.New(color.FgHiWhite)
	hiBlue := color.New(color.FgHiBlue)
	for _, e := range extractors.ExtractorMap {
		resInt, _ := redis.Int(CONN.Do("SCARD", e.Key()))
		hiBlue.Printf("%s%s  ", hiBlue.Sprintf("%s：", e.Name()), hiWhite.Sprint(resInt))
	}
	fmt.Println("")
}
