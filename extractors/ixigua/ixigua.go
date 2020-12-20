package ixigua

import "gitee.com/rock_rabbit/goannie/extractors/types"

type extractor struct{}

// Name 平台名称
func (e *extractor) Name() string {
	return "西瓜视频"
}

// Key 存储器使用的键
func (e *extractor) Key() string {
	return "xigua"
}

// Extract 运行解析器
func (e *extractor) Extract(url string, option types.Options) error {
	return nil
}

// New 创建一个解析器
func New() types.Extractor {
	return &extractor{}
}
