package downloader

import "io"

// IoProxyReader 代理io读
type IoProxyReader struct {
	io.Reader
	info *DownladerInfo
}

// Read 读
func (r *IoProxyReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	// 读完直接加入到进度
	r.info.AddDownloadedSize(int64(n))
	// 外部进度更新回调
	go r.info.OnProgress(r.info.Size, r.info.Speed, r.info.GetDownloadedSize(), r.info.Context)
	// 监听上下文是否结束
	select {
	case <-r.info.Context.Done():
		// 上下文关闭，读取结束
		r.Close()
	default:
	}
	return n, err
}

// Close the wrapped reader when it implements io.Closer
func (r *IoProxyReader) Close() (err error) {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}

// ProxyReader 代理io读
func (info *DownladerInfo) ProxyReader(r io.Reader) io.Reader {
	return &IoProxyReader{r, info}
}
