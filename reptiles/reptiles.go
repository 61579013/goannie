package reptiles

import (
	"errors"
	"net/url"
	"strings"

	"gitee.com/rock_rabbit/goannie/reptiles/ixigua"
	"gitee.com/rock_rabbit/goannie/reptiles/types"
	"gitee.com/rock_rabbit/goannie/utils"
)

// ReptilesMap 采集器
var ReptilesMap map[string]types.Reptiles

func init() {
	ReptilesMap = map[string]types.Reptiles{
		"ixigua": ixigua.New(),
	}
}

// Extract 运行采集器
func Extract(u string, option types.Options) ([]*types.Data, error) {
	var (
		err    error
		domain string
	)
	u = strings.TrimSpace(u)
	parseU, err := url.ParseRequestURI(u)
	if err != nil {
		return nil, err
	}
	if parseU.Host == "haokan.baidu.com" {
		domain = "haokan"
	} else {
		domain = utils.Domain(parseU.Host)
	}
	if _, ok := ReptilesMap[domain]; !ok {
		return nil, errors.New("暂不支持此链接")
	}
	reptiles := ReptilesMap[domain]
	return reptiles.Extract(u, option)
}
