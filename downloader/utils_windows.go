package downloader

import (
	"runtime"
	"syscall"
	"time"
)

// GetFileCreateTime 获取文件创建时间
func GetFileCreateTime(sys interface{}) int64 {
	osType := runtime.GOOS
	if osType == "windows" {
		wFileSys := sys.(*syscall.Win32FileAttributeData)
		tNanSeconds := wFileSys.CreationTime.Nanoseconds() /// 返回的是纳秒
		tSec := tNanSeconds / 1e9                          ///秒
		return tSec
	}
	return time.Now().Unix()
}
