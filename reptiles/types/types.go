package types

import (
	"regexp"

	"gitee.com/rock_rabbit/goannie/storage"
)

// Extract 运行
type Extract = func(url string, option Options) ([]*Data, error)

// SubCmd 子命令
type SubCmd struct {
	URLRegexps []*regexp.Regexp
	Extract    Extract
}

// SubCmdList 子命令集
type SubCmdList []SubCmd

// Parse 解析子命令
func (s SubCmdList) Parse(url string) Extract {
	for _, sub := range s {
		for _, cmd := range sub.URLRegexps {
			if cmd.MatchString(url) {
				return sub.Extract
			}
		}
	}
	return nil
}

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
