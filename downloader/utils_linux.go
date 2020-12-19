package downloader

import (
	"runtime"
	"syscall"
	"time"
)

// GetFileCreateTime 获取文件创建时间
func GetFileCreateTime(sys interface{}) int64 {
	osType := runtime.GOOS
	if osType == "linux" {
		stat_t := sys.(*syscall.Stat_t)
		tCreate := int64(stat_t.Ctim.Sec)
		return tCreate
	}
	return time.Now().Unix()
}
