package main

import (
	"fmt"
	"os"
	"os/exec"

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

}

func sayHello() {
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)
	hiBlue := color.New(color.FgHiBlue)
	hiWhite := color.New(color.FgHiWhite)
	green.Println(config.TITLE)
	magenta.Printf("	版本: %s	更新时间: %s\n\n", config.VERSION, config.UPDATETIME)
	hiBlue.Printf("%s %s%s %s%s\n", green.Sprint("$"), hiBlue.Sprint("作者昵称："), hiWhite.Sprint("rockrabbit"), hiBlue.Sprint("作者主页："), hiWhite.Sprint("https://www.68wu.com"))
	hiBlue.Printf("%s %s%s\n", green.Sprint("$"), hiBlue.Sprint("软件主页："), hiWhite.Sprint("https://gitee.com/rock_rabbit/goannie"))
	fmt.Println("")
	hiBlue.Printf("……………………………… %s %s\n", hiWhite.Sprint("下载统计"), hiBlue.Sprint("………………………………"))
	videoIDCount()
	fmt.Println("")
}

// videoIDCount 打印过滤库个数
func videoIDCount() {
	hiWhite := color.New(color.FgHiWhite)
	hiBlue := color.New(color.FgHiBlue)
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
		}, {
			"name":  "抖音视频",
			"pt":    "douyin",
			"count": "",
		}, {
			"name":  "优酷视频",
			"pt":    "youku",
			"count": "",
		},
	}
	for idx, item := range ptList {
		resInt, _ := redis.Int(CONN.Do("SCARD", item["pt"]))
		ptList[idx]["count"] = fmt.Sprintf("%d", resInt)
		hiBlue.Printf("%s%s  ", hiBlue.Sprintf("%s：", item["name"]), hiWhite.Sprint(ptList[idx]["count"]))
	}
	fmt.Println("")
}
