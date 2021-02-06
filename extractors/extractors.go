package extractors

import (
	"errors"
	"net/url"
	"strings"

	"gitee.com/rock_rabbit/goannie/extractors/iqiyi"
	"gitee.com/rock_rabbit/goannie/extractors/ixigua"
	"gitee.com/rock_rabbit/goannie/extractors/types"
	"gitee.com/rock_rabbit/goannie/utils"
)

// ExtractorMap 解析器
var ExtractorMap map[string]types.Extractor

func init() {
	ExtractorMap = map[string]types.Extractor{
		"ixigua": ixigua.New(),
		"iqiyi":  iqiyi.New(),
	}
}

// Extract 运行解析器
func Extract(u string, option types.Options) error {
	var (
		err    error
		domain string
	)
	u = strings.TrimSpace(u)
	parseU, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}
	if parseU.Host == "haokan.baidu.com" {
		domain = "haokan"
	} else {
		domain = utils.Domain(parseU.Host)
	}
	if _, ok := ExtractorMap[domain]; !ok {
		return errors.New("暂不支持此链接")
	}
	extractor := ExtractorMap[domain]
	if err = extractor.Extract(u, option); err != nil {
		return err
	}
	return nil
}
