package types

import (
	"fmt"

	"gitee.com/rock_rabbit/goannie/request"
	"gitee.com/rock_rabbit/goannie/storage"
	"github.com/fatih/color"
)

// Options 解析器&下载器参数
type Options struct {
	Cookie   string
	SavePath string
	Verify   storage.Storage
	Stream   string
}

// Extractor 解析器&下载器主要实现接口
type Extractor interface {
	// Name 平台名称
	Name() string
	// Key 存储器使用的键
	Key() string

	DefCookie() string
	// Extract 解析器&下载器运行
	Extract(url string, option Options) error
}

// DownloadPrint 程序下载时打印结构体
type DownloadPrint struct {
	Site      string
	Title     string
	Type      string
	Quality   string
	Size      string
	SizeBytes int64
}

// NewDownloadPrint 创建
func NewDownloadPrint(site, title, quality, url string) *DownloadPrint {
	d := &DownloadPrint{
		Site:      site,
		Title:     title,
		Type:      "video",
		Quality:   quality,
		Size:      "",
		SizeBytes: 0,
	}

	d.SizeBytes, _ = request.Size(url, url)
	d.FormatSize()
	return d
}

// FormatSize 格式化字节
func (d *DownloadPrint) FormatSize() {
	fileSize := d.SizeBytes
	if fileSize < 1024 {
		d.Size = fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		d.Size = fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// Print 打印
func (d DownloadPrint) Print() {
	fmt.Println("")
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Site:      ")
	color.Unset()
	fmt.Println(d.Site)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Title:     ")
	color.Unset()
	fmt.Println(d.Title)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Type:      ")
	color.Unset()
	fmt.Println(d.Type)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Streams:   ")
	color.Unset()
	fmt.Println("# All available quality")
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("     [default]  -------------------\n")
	color.Unset()
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("     Quality:         ")
	color.Unset()
	fmt.Println(d.Quality)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("     Size:            ")
	color.Unset()
	fmt.Printf("%s (%d Bytes)\n", d.Size, d.SizeBytes)
	fmt.Println("")
}
