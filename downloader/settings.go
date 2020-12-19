package downloader

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"sync/atomic"

	"github.com/cheggaaa/pb/v3"
)

// AddDownloadedSize 增加已下载文件大小
func (info *DownladerInfo) AddDownloadedSize(v int64) {
	atomic.AddInt64(info.DownloadedSize, v)
}

// GetDownloadedSize 获取已下载文件大小
func (info *DownladerInfo) GetDownloadedSize() int64 {
	return atomic.LoadInt64(info.DownloadedSize)
}

// SetOnProgress 设置进度事件
func (dl *Downloader) SetOnProgress(f OnProgress) *Downloader {
	dl.OnProgress = f
	return dl
}

// SetClient 设置请求Client
func (dl *Downloader) SetClient(Client *http.Client) *Downloader {
	dl.Options.Rquest.Client = Client
	return dl
}

// SetMethod 设置请求Method
func (dl *Downloader) SetMethod(method string) *Downloader {
	dl.Options.Rquest.Method = method
	return dl
}

// SetBody 设置请求Body
func (dl *Downloader) SetBody(body io.Reader) *Downloader {
	dl.Options.Rquest.Body = body
	return dl
}

// SetHeader 设置请求时的Header
func (dl *Downloader) SetHeader(name, value string) *Downloader {
	dl.Options.Rquest.Header.Set(name, value)
	return dl
}

// AddHeader 添加请求时的Header
func (dl *Downloader) AddHeader(name, value string) *Downloader {
	dl.Options.Rquest.Header.Add(name, value)
	return dl
}

// AddDfer 添加Defer 下载关闭事件
func (dl *Downloader) AddDfer(d OnDefer) *Downloader {
	dl.Defer = append(dl.Defer, d)
	return dl
}

// SetOutputPath 设置输出目录
func (dl *Downloader) SetOutputPath(outputPath string) *Downloader {
	if outputPath == "" {
		outputPath = "./"
	}
	dl.Options.OutputPath = outputPath
	return dl
}

// SetOutputName 设置输出文件名称
func (dl *Downloader) SetOutputName(outputName string) *Downloader {
	dl.Options.OutputName = outputName
	return dl
}

// SetIsBar 设置是否打印进度条
func (dl *Downloader) SetIsBar(i bool) *Downloader {
	dl.ProgressBar.IsBar = i
	return dl
}

// GetIsBar 获取是否打印进度条
func (dl *Downloader) GetIsBar() bool {
	return dl.ProgressBar.IsBar
}

// SetBarTemplate 设置进度条模板
func (dl *Downloader) SetBarTemplate(BarTemplate pb.ProgressBarTemplate) *Downloader {
	dl.ProgressBar.BarTemplate = BarTemplate
	return dl
}

// GetBarTemplate 获取进度条模板
func (dl *Downloader) GetBarTemplate() pb.ProgressBarTemplate {
	return dl.ProgressBar.BarTemplate
}

// SetThreadNum 设置下载线程数
func (dl *Downloader) SetThreadNum(n int) *Downloader {
	if n <= 1 {
		dl.Options.ThreadNum = 1
	} else {
		dl.Options.ThreadNum = n
	}
	return dl
}

// GetThreadNum 获取下载线程数
func (dl *Downloader) GetThreadNum() int {
	return dl.Options.ThreadNum
}

// GetPath 获取文件完整路径
func (dl *Downloader) GetPath() string {
	return path.Join(dl.Options.OutputPath, dl.Options.OutputName)
}

// GetTempName 获取临时文件名
func (dl *Downloader) GetTempName() string {
	return fmt.Sprintf("%s.download", dl.Options.OutputName)
}

// GetTempPath 获取临时文件完整路径
func (dl *Downloader) GetTempPath() string {
	return path.Join(dl.Options.OutputPath, dl.GetTempName())
}
