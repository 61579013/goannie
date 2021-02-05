package iqiyi

import (
	"errors"
	"fmt"

	"gitee.com/rock_rabbit/goannie/config"
	"gitee.com/rock_rabbit/goannie/extractors/types"
	"gitee.com/rock_rabbit/goannie/request"
	"gitee.com/rock_rabbit/goannie/utils"
	"gitee.com/rock_rabbit/goannie/utils/annie"
)

const (
	referer       = "https://www.iqiyi.com"
	defaultCookie = ""
)

type extractor struct{}

// Name 平台名称
func (e *extractor) Name() string {
	return "爱奇艺视频"
}

// Key 存储器使用的键
func (e *extractor) Key() string {
	return "iqiyi"
}

// DefCookie 默认cookie
func (e *extractor) DefCookie() string {
	return defaultCookie
}

// New 创建一个解析器
func New() types.Extractor {
	return &extractor{}
}

// Extract 运行解析器
func (e *extractor) Extract(url string, option types.Options) error {
	html, err := request.GetByte(url, referer, nil)
	if err != nil {
		return err
	}
	regexpVid := utils.MatchOneOfByte(html, `param['vid'] = "(.*?)"`, `"vid":"(.*?)"`)
	if len(regexpVid) < 2 {
		return errors.New("获取vid失败")
	}
	vid := string(regexpVid[1])
	if config.GetBool("app.isFiltrationID") && vid != "" && option.Verify.Check(e.Key(), vid) {
		return fmt.Errorf("%s 在过滤库中，修改配置文件isFiltrationID为false不再过滤。", vid)
	}

	if err := annie.Download(url, option.SavePath, option.Cookie); err != nil {
		return err
	}

	option.Verify.Add(e.Key(), vid)
	return nil
}
