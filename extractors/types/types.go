package types

import "gitee.com/rock_rabbit/goannie/storage"

// Options 解析器&下载器参数
type Options struct {
	Cookie string
	Verify storage.Storage
}

// Extractor 解析器&下载器主要实现接口
type Extractor interface {
	// Name 平台名称
	Name() string
	// Key 存储器使用的键
	Key() string
	// Extract 解析器&下载器运行
	Extract(url string, option Options) error
}
