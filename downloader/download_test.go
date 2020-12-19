package downloader_test

import (
	"context"
	"fmt"
	"testing"

	"gitee.com/rock_rabbit/downloader"
)

func TestDownload(t *testing.T) {
	var err error
	// 测试下载
	url := "http://speed.hetzner.de/100MB.bin"
	outPath := `D:\测试下载`
	dl := downloader.New(url, outPath)
	// 设置不显示进度条
	dl.SetIsBar(false)
	// 设置进度回调
	dl.SetOnProgress(func(size, speed, downloadedSize int64, context context.Context) {
		fmt.Printf("\r总大小：%d byte  已下载：%d byte  下载速度：%d byte/s", size, downloadedSize, speed)
	})
	// 添加 Start() 结束后的回调
	dl.AddDfer(func(dl *downloader.Downloader) {
		fmt.Printf("  下载结束\n")
	})
	// 设置文件名，不设置自动获取文件名
	dl.SetOutputName("100MB.bin")
	// 启动下载器
	err = dl.Run()
	if err != nil {
		t.Log(err)
	}
}
