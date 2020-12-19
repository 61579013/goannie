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
	if err := binary.Update(); err != nil {
		utils.ErrInfo(err.Error())
	}
	// 检查Data目录
	if isDataPath, _ := utils.IsExist(config.AppDataPath); !isDataPath {
		if err := os.MkdirAll(config.AppDataPath, os.ModePerm); err != nil {
			utils.ErrInfo(err.Error())
			utils.ExitInfo()
		}
	}
	// 启动 redis
	cmd := exec.Command("cmd", "/c", "start", "/B", "redis-server", config.RedisConfFile, "--dir", config.AppDataPath)
	if err = cmd.Start(); err != nil {
		utils.ErrInfo(err.Error())
	}
	// 连接 redis
	CONN, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		utils.ErrInfo(err.Error())
		utils.ExitInfo()
	}
	sayHello()
}

// CONN redis服务连接
var CONN redis.Conn

func main() {
	defer CONN.Close()

}

func sayHello() {
	green := color.New(color.FgGreen).PrintlnFunc()
	magenta := color.New(color.FgMagenta).PrintfFunc()
	hiBlue := color.New(color.FgHiBlue).PrintlnFunc()

	green(config.TITLE)
	magenta("	版本: %s	更新时间: %s\n\n", config.VERSION, config.UPDATETIME)
	hiBlue("> 下载统计")
	videoIDCount()
	fmt.Println("")
}

// videoIDCount 打印过滤库个数
func videoIDCount() {
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
		fmt.Printf("%s：%s  ", item["name"], ptList[idx]["count"])
	}
	fmt.Println("")
}
