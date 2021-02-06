package iqiyi

import (
	"errors"
	"regexp"

	"gitee.com/rock_rabbit/goannie/reptiles/types"
)

const (
	referer = "https://www.iqiyi.com"
	host    = "https://www.iqiyi.com"
)

// SubCmdList 子命令集
var SubCmdList types.SubCmdList

func init() {
	SubCmdList = []types.SubCmd{}
	// 视频内容页
	SubCmdList = append(SubCmdList, types.SubCmd{
		URLRegexps: []*regexp.Regexp{
			regexp.MustCompile(`www\.iqiyi\.com/v_\w+\.html`),
		},
		Extract: insidePage,
	})
	// 视频作者
	SubCmdList = append(SubCmdList, types.SubCmd{
		URLRegexps: []*regexp.Regexp{
			regexp.MustCompile(`www\.iqiyi\.com/u/.*?`),
		},
		Extract: homePage,
	})
}

type reptiles struct{}

// Extract 运行采集器
func (e *reptiles) Extract(url string, option types.Options) ([]*types.Data, error) {
	if extract := SubCmdList.Parse(url); extract != nil {
		return extract(url, option)
	}
	return nil, errors.New("暂不支持此链接")
}

// New 创建一个采集器
func New() types.Reptiles {
	return &reptiles{}
}

func insidePage(url string, option types.Options) ([]*types.Data, error) {
	return []*types.Data{&types.Data{URL: url}}, nil
}

func homePage(url string, option types.Options) ([]*types.Data, error) {
	return []*types.Data{&types.Data{URL: url}}, nil
}
