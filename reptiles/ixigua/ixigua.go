package ixigua

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gitee.com/rock_rabbit/goannie/config"

	"gitee.com/rock_rabbit/goannie/utils"

	"gitee.com/rock_rabbit/goannie/reptiles/types"
	"gitee.com/rock_rabbit/goannie/request"
)

const (
	referer = "https://www.ixigua.com"
	host    = "https://www.ixigua.com"
)

// SubCmdList 子命令集
var SubCmdList types.SubCmdList

func init() {
	SubCmdList = []types.SubCmd{}
	// 视频内容页
	SubCmdList = append(SubCmdList, types.SubCmd{
		URLRegexps: []*regexp.Regexp{
			regexp.MustCompile(`ixigua\.com/(\d+$|\d+/$|\d+\?|\d+/\?)`),
		},
		Extract: insidePage,
	})
	// 作者作品页
	SubCmdList = append(SubCmdList, types.SubCmd{
		URLRegexps: []*regexp.Regexp{
			regexp.MustCompile(`ixigua\.com/home/(\d+$|\d+/$|\d+\?|\d+/\?)`),
			regexp.MustCompile(`ixigua\.com/home/\d+/video($|/$|\?|/\?)`),
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
	start, end, userid, err := homePageGetInput(url)
	if err != nil {
		return nil, err
	}
	retData := []*types.Data{}
	onCount := 0
	onMaxTime := 0
	for {
		if onCount > end {
			break
		}
		authorVideo, err := getHomeAuthorVideo(userid, onMaxTime)
		if err != nil {
			utils.ErrInfo(err.Error())
			break
		}
		for _, v := range authorVideo.Data.Data {
			if onCount > end {
				break
			}
			regexpTitle := regexp.MustCompile(config.GetString("reptiles.regexpTitle")).MatchString(v.Title)
			if onCount >= start && onCount <= end && regexpTitle {
				retData = append(retData, &types.Data{URL: fmt.Sprintf("https://www.ixigua.com/%s", v.Gid)})
			}
			onMaxTime = v.PublishTime
			onCount++
		}
	}
	utils.InfoKv("数量：", fmt.Sprintf("%d", len(retData)))
	return retData, nil
}

func homePageGetInput(url string) (int, int, string, error) {
	HydratedData, err := getHomeHydratedData(url)
	if err != nil {
		return 0, 0, "", err
	}
	//打印作者信息
	fmt.Println("")
	utils.InfoKv("作者：", HydratedData.AuthorDetailInfo.Name)
	utils.InfoKv("认证：", HydratedData.AuthorDetailInfo.VerifiedContent)
	utils.InfoKv("简介：", HydratedData.AuthorDetailInfo.Introduce)
	utils.InfoKv("作品数量：", fmt.Sprintf("%d", HydratedData.AuthorTabsCount.VideoCnt))
	fmt.Println("")
	videoCnt := HydratedData.AuthorTabsCount.VideoCnt
	userid := HydratedData.AuthorDetailInfo.UserID
	var (
		start int
		end   int
	)
	if err := utils.GetIntInput("起始数", &start); err != nil {
		return 0, 0, "", err
	}
	if err := utils.GetIntInput("结束数", &end); err != nil {
		return 0, 0, "", err
	}
	if start > videoCnt || start <= 0 || end <= 0 || end > videoCnt || start > end {
		return 0, 0, "", errors.New("值超出范围")
	}
	return start, end, userid, nil
}

func getHomeHydratedData(url string) (*homeHydratedData, error) {
	html, err := request.GetByte(url, referer, nil)
	if err != nil {
		return nil, err
	}
	jsonStrFind := regexp.MustCompile(`window\._SSR_HYDRATED_DATA=(.*?)</script>`).FindSubmatch(html)
	if len(jsonStrFind) < 2 {
		return nil, errors.New("解析数据失败")
	}
	jsonStr := strings.ReplaceAll(string(jsonStrFind[1]), ":undefined", ":\"undefined\"")
	var jsonData homeHydratedData
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}

func getHomeAuthorVideo(authorID string, maxTime int) (*homeAuthorVideo, error) {
	api := fmt.Sprintf("https://www.ixigua.com/api/videov2/author/video?author_id=%s&type=video&max_time=%d", authorID, maxTime)
	var jsonData homeAuthorVideo
	if err := request.GetJSON(api, referer, nil, &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}
