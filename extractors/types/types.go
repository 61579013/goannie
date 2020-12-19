package types

// Options 解析器&下载器参数
type Options struct {
	Cookie string
}

// Extractor 解析器&下载器主要实现接口
type Extractor interface {
	// Extract 解析器&下载器运行
	Extract(url string, option Options) error
}
