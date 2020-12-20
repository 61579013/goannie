package config

import (
	"fmt"
	"os"

	"gitee.com/rock_rabbit/goannie/utils"

	"github.com/spf13/viper"
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
	// Config 基本配置
	Config *viper.Viper
)

// 初始化基本配置
func init() {
	Config = viper.New()
	path, _ := utils.GetCurrentPath()
	Config.AddConfigPath(path)
	Config.SetConfigName("config")
	Config.SetConfigType("toml")
	// 设置默认参数
	Config.SetDefault("app.checkDuplication", true)
	Config.SetDefault("app.autoCreatePath", true)
	Config.SetDefault("app.isFiltrationID", true)
	Config.SetDefault("binary.check", true)
	Config.SetDefault("binary.UpdateNetworkJSONFile", "http://image.68wu.cn/goannie/binary_version.json")
	Config.SetDefault("binary.UpdateLockTimeOut", 150)
	Config.SetDefault("redis.start", true)
	Config.SetDefault("redis.dial", true)
	Config.SetDefault("redis.network", "tcp")
	Config.SetDefault("redis.address", "localhost:6379")
	Config.SetDefault("outpath.p1", "./保存目录")
	Config.SetDefault("outpath.p2", "")
	Config.SetDefault("outpath.p3", "")
	Config.SetDefault("outpath.p4", "")
	Config.SetDefault("outpath.p5", "")
	// 读取配置，如果不存在直接创建默认配置
	if err := Config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			utils.ErrInfo(err.Error())
		} else {
			if err := WriteConfig(); err != nil {
				utils.ErrInfo(err.Error())
			}
		}
	}
}

// WriteConfig 保存配置文件
func WriteConfig() error {
	var err error
	if err = Config.SafeWriteConfig(); err == nil {
		return nil
	}
	if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
		return err
	}
	if err := Config.WriteConfig(); err != nil {
		return err
	}
	return nil
}

// WatchConfig 重新读取配置
func WatchConfig() {
	Config.WatchConfig()
}

// Set 设置配置
func Set(key string, value interface{}) {
	Config.Set(key, value)
}

// GetString 获取String配置
func GetString(key string) string {
	return Config.GetString(key)
}

// GetInt 获取Int配置
func GetInt(key string) int {
	return Config.GetInt(key)
}

// GetBool 获取Bool配置
func GetBool(key string) bool {
	return Config.GetBool(key)
}
