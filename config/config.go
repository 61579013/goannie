package config

import (
	"fmt"
	"os"
)

const (
	// VERSION 版本号
	VERSION = "1.0.0"
	// UPDATETIME 更新时间
	UPDATETIME = "2020-12-19"
	// TITLE 软件标题
	TITLE = `
                                        __           
   __     ___      __      ___     ___ /\_\     __   
 /'_ ` + "`" + `\  / __` + "`" + `\  /'__` + "`" + `\  /' _ ` + "`" + `\ /' _ ` + "`" + `\/\ \  /'__` + "`" + `\
/\ \L\ \/\ \L\ \/\ \L\.\_/\ \/\ \/\ \/\ \ \ \/\  __/
\ \____ \ \____/\ \__/.\_\ \_\ \_\ \_\ \_\ \_\ \____\
 \/___L\ \/___/  \/__/\/_/\/_/\/_/\/_/\/_/\/_/\/____/
   /\____/
   \_/__/`
)

var (
	// AppPath 程序app目录
	AppPath = fmt.Sprintf("%s\\goannie", os.Getenv("APPDATA"))
	// AppBinPath 程序bin目录
	AppBinPath = fmt.Sprintf("%s\\bin", AppPath)
	// AppDataPath 程序data目录
	AppDataPath = fmt.Sprintf("%s\\data", AppPath)
	// AnnieFile 程序annie存放位置
	AnnieFile = fmt.Sprintf("%s\\annie.exe", AppBinPath)
	// FfmpegFile 程序ffmpeg存放位置
	FfmpegFile = fmt.Sprintf("%s\\ffmpeg.exe", AppBinPath)
	// Aria2File 程序aria2存放位置
	Aria2File = fmt.Sprintf("%s\\aria2c.exe", AppBinPath)
	// RedisFile 程序redis存放位置
	RedisFile = fmt.Sprintf("%s\\redis-server.exe", AppBinPath)
	// RedisConfFile 程序redisconf存放位置
	RedisConfFile = fmt.Sprintf("%s\\redis.windows-service.conf", AppBinPath)
)
