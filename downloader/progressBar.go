package downloader

import (
	"io"

	"github.com/cheggaaa/pb/v3"
)

// progressBar.go 下载进度条实现

// BarThreadOneReader 单线程下载进度条，读取IO
type BarThreadOneReader struct {
	io.Reader
	bar *pb.ProgressBar
}

func (r *BarThreadOneReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.bar.Add64(int64(n))
	return n, err
}

// BarThreadOne 单线程下载进度条
func BarThreadOne(dl *Downloader, startSize, size int64, r io.Reader) io.Reader {
	reader := io.LimitReader(r, size)
	// 可以读取到文件大小时显示进度条
	if size > 0 {
		bar := dl.GetBarTemplate().Start64(size)
		bar.Add64(startSize)
		barReader := bar.NewProxyReader(reader)
		dl.AddDfer(func(dl *Downloader) {
			bar.Finish()
		})
		return barReader
	}
	tmpl := `{{string . "prefix"}} {{counters . }} {{speed . }} {{string . "suffix"}}`
	bar := pb.ProgressBarTemplate(tmpl).Start64(0)
	bar.Add64(startSize)
	bar.Set(pb.Bytes, true)
	// 读取不到时只显示已下载大小
	dl.AddDfer(func(dl *Downloader) {
		bar.Finish()
	})
	return &BarThreadOneReader{r, bar}
}
