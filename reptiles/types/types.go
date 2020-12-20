package types

import "gitee.com/rock_rabbit/goannie/storage"

// Data 采集器返回的数据
type Data struct {
	URL string
}

// Options 采集器参数
type Options struct {
	Cookie string
	Verify storage.Storage
}

// Reptiles 采集器主要实现接口
type Reptiles interface {
	// Extract 采集器运行
	Extract(url string, option Options) ([]*Data, error)
}
