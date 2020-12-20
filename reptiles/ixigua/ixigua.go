package ixigua

import "gitee.com/rock_rabbit/goannie/reptiles/types"

type reptiles struct{}

// Extract 运行采集器
func (e *reptiles) Extract(url string, option types.Options) ([]*types.Data, error) {
	return nil, nil
}

// New 创建一个采集器
func New() types.Reptiles {
	return &reptiles{}
}
