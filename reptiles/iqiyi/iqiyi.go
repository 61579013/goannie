package iqiyi

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"gitee.com/rock_rabbit/goannie/config"
	"gitee.com/rock_rabbit/goannie/reptiles/types"
	"gitee.com/rock_rabbit/goannie/request"
	"gitee.com/rock_rabbit/goannie/utils"
)

const (
	referer = "https://www.iqiyi.com"
	host    = "https://www.iqiyi.com"

	authCookie   = ""
	dfp          = "a0b9b7387466a14cb98c18230de05e42b0eecd0926d67e4a13f54286c435fc9660"
	mDeviceID    = "go5e8hbplm0iumdxpsat9upx"
	agentversion = "10.7.5"
	salt         = "NZrFGv72GYppTUxO"
	defaultSize  = 28
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
	return []*types.Data{{URL: url}}, nil
}

func homePage(url string, option types.Options) ([]*types.Data, error) {
	fuidMatch := utils.MatchOneOf(url, `/u/(\d+)`)
	if len(fuidMatch) < 2 {
		return nil, errors.New("获取fuid失败")
	}
	d, e := videosAction(fuidMatch[1], "1", fmt.Sprint(defaultSize))
	if e != nil {
		return nil, e
	}
	page := int(math.Ceil(float64(d.Data.TotalNum) / float64(defaultSize)))
	fmt.Println("")
	utils.InfoKv("fuid：", fuidMatch[1])
	utils.InfoKv("总页数：", fmt.Sprint(page))
	utils.InfoKv("每页个数：", fmt.Sprint(defaultSize))
	utils.InfoKv("作品数量：", fmt.Sprint(d.Data.TotalNum))
	fmt.Println("")
	var (
		start int
		end   int
	)
	if err := utils.GetIntInput("起始页", &start); err != nil {
		return nil, err
	}
	if err := utils.GetIntInput("结束页", &end); err != nil {
		return nil, err
	}
	if start > page || start <= 0 || end <= 0 || end > page || start > end {
		return nil, errors.New("值超出范围")
	}
	retData := []*types.Data{}
	template := "\r当前页数：%s 采集范围：" + fmt.Sprintf("%d-%d", start, end) + " 已采集数：%d 错误数：%d"
	var errCount int
	render := func() {
		utils.Infof(template, fmt.Sprint(start), len(retData), errCount)
	}
	for {
		if start > end {
			break
		}
		l, err := videosAction(fuidMatch[1], fmt.Sprint(start), fmt.Sprint(defaultSize))
		if err != nil {
			if errCount <= 10 {
				time.Sleep(time.Second)
				errCount++
				render()
				continue
			}
			return retData, nil
		}
		var qipuIds bytes.Buffer
		if len(l.Data.Sort.Flows) == 0 {
			start++
			render()
			continue
		}
		for _, i := range l.Data.Sort.Flows {
			qipuIds.WriteString(fmt.Sprintf("%d,", i.QipuID))
		}
		qipuIdsStr := qipuIds.String()
		v, err := episodeInfo(qipuIdsStr[:len(qipuIdsStr)-1])
		if err != nil {
			if errCount <= 10 {
				time.Sleep(time.Second)
				errCount++
				render()
				continue
			}
			return retData, nil
		}
		for _, i := range v.Data {
			regexpTitle := regexp.MustCompile(config.GetString("reptiles.regexpTitle")).MatchString(i.Title)
			if regexpTitle {
				retData = append(retData, &types.Data{URL: i.PageURL})
				render()
			}
		}
		start++
	}
	fmt.Println("")
	utils.InfoKv("采集数量：", fmt.Sprintf("%d", len(retData)))
	return retData, nil
}

func videosAction(fuid, page, size string) (*videosActionAPI, error) {
	timestamp := fmt.Sprint(time.Now().Unix() * 1000)
	parameters := strings.Join([]string{
		"agenttype=118", "agentversion=" + agentversion, "authcookie=" + authCookie, "dfp=" + dfp, "fuid=" + fuid, "m_device_id=" + mDeviceID,
		"page=" + page, "size=" + size, "timestamp=" + fmt.Sprint(timestamp),
	}, "&")
	sign := utils.Md5(fmt.Sprintf("GET%s?%s%s", "iqiyihao.iqiyi.com/iqiyihao/entity/get_videos.action", parameters, salt))
	api := fmt.Sprintf("https://iqiyihao.iqiyi.com/iqiyihao/entity/get_videos.action?%s", parameters+"&sign="+sign)
	var jsonData videosActionAPI
	if err := request.GetJSON(api, referer, nil, &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}

func episodeInfo(qipuIds string) (*episodeInfoAPI, error) {
	timestamp := fmt.Sprint(time.Now().Unix() * 1000)
	parameters := strings.Join([]string{
		"agenttype=118", "agentversion=" + agentversion, "authcookie=" + authCookie, "dfp=" + dfp, "m_device_id=" + mDeviceID,
		"qipuIds=" + qipuIds, "timestamp=" + fmt.Sprint(timestamp),
	}, "&")
	sign := utils.Md5(fmt.Sprintf("GET%s?%s%s", "iqiyihao.iqiyi.com/iqiyihao/episode_info.action", parameters, salt))
	api := fmt.Sprintf("https://iqiyihao.iqiyi.com/iqiyihao/episode_info.action?%s", parameters+"&sign="+sign)
	var jsonData episodeInfoAPI
	if err := request.GetJSON(api, referer, nil, &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}
