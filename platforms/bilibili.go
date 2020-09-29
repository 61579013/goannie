package platforms

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// RunBliOne 单视频
func RunBliOne(runType RunType, arg map[string]string) error {
	bvid := BliGetBvid(runType.URL)
	// 判断是否需要多p交易
	if BliIsP(bvid) {
		return BliGetVideoP(runType, bvid)
	}
	// 判断是否过滤重复
	if bvid != "" {
		isVID := IsVideoID("bilibili", bvid, runType.RedisConn)
		if isVID && runType.IsDeWeight {
			// 判断到了重复
			PrintErrInfo(RepetitionMsg)
			return nil
		}
	}
	err := AnnieDownload(runType.URL, runType.SavePath, runType.CookieFile, runType.DefaultCookie)
	if err != nil {
		return err
	}
	// 存储已下载
	if bvid != "" {
		AddVideoID("bilibili", bvid, runType.RedisConn)
	}
	return nil
}

// BliIsP 判断是否需要多P交易
func BliIsP(bvid string) bool {
	url := "https://api.bilibili.com/x/player/pagelist?bvid=%s&jsonp=json"
	var jsonData BliVideoP
	if err := RequestGetJSON(fmt.Sprintf(url, bvid), map[string]string{
		"accept":     "*/*",
		"referer":    "https://www.bilibili.com/",
		"user-agent": UserAgentPc,
	}, &jsonData); err != nil {
		return false
	}
	if len(jsonData.Data) > 1 {
		return true
	}
	return false
}

// BliGetVideoP 多p交易
func BliGetVideoP(runType RunType, bvid string) error {
	url := "https://api.bilibili.com/x/player/pagelist?bvid=%s&jsonp=json"
	var jsonData BliVideoP
	if err := RequestGetJSON(fmt.Sprintf(url, bvid), map[string]string{
		"accept":     "*/*",
		"referer":    "https://www.bilibili.com/",
		"user-agent": UserAgentPc,
	}, &jsonData); err != nil {
		return err
	}
	maxP := len(jsonData.Data)
	PrintInfo(fmt.Sprintf("视频选集数量：%d", maxP))
	var (
		start, end string
	)
	if err := GetCmdDataString("请输入起始P", &start); err != nil {
		return err
	}
	startInt, err := strconv.Atoi(start)
	if err != nil || startInt > maxP || startInt <= 0 {
		return errors.New("起始P格式错误")
	}
	if err := GetCmdDataString("请输入结束P", &end); err != nil {
		return err
	}
	endInt, err := strconv.Atoi(end)
	if err != nil || endInt > maxP || endInt < startInt || endInt <= 0 {
		return errors.New("结束P格式错误")
	}
	var downLoadList []map[string]string

	for idx, item := range jsonData.Data {
		if (idx+1) >= startInt && (idx+1) <= endInt {
			downLoadList = append(downLoadList, map[string]string{
				"vid":   fmt.Sprintf("bvid_%d", item.Page),
				"title": item.Part,
				"url":   fmt.Sprintf("https://www.bilibili.com/video/%s?p=%d", bvid, item.Page),
			})
		}
	}
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))
	AnnieDownloadAll(downLoadList, runType, "bilibili")
	PrintInfo("全部下载完成")
	return nil
}

// BliGetBvid 通过url获取bvid
func BliGetBvid(url string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://www\.bilibili\.com/video/(\w+)($|\?.*?$)`),
		regexp.MustCompile(`^(http|https)://www\.bilibili\.com/bangumi/play/(\w+)($|\?.*?$)`),
	}
	for _, regxp := range regexps {
		Bvid := regxp.FindStringSubmatch(url)
		if len(Bvid) < 3 {
			continue
		}
		return Bvid[2]
	}
	return ""
}

// RunBliUserList 作者列表
func RunBliUserList(runType RunType, arg map[string]string) error {
	userID, err := bliGetUserID(runType.URL)
	if err != nil {
		return err
	}
	tid := bliGetTID(runType.URL)
	order := bliGetOrder(runType.URL)
	keyword := bliGetKeyword(runType.URL)
	page, count, err := bliGetMaxPage(userID, tid, order, keyword)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("\r总页数：%d  每页个数：%d  总个数：%d", page, 30, count))
	var (
		start, end string
	)
	if err = GetCmdDataString("请输入起始页", &start); err != nil {
		return err
	}
	startInt, err := strconv.Atoi(start)
	if err != nil || startInt > page || startInt <= 0 {
		return errors.New("起始页格式错误")
	}
	if err = GetCmdDataString("请输入结束页", &end); err != nil {
		return err
	}
	endInt, err := strconv.Atoi(end)
	if err != nil || endInt > page || endInt < startInt || endInt <= 0 {
		return errors.New("结束页格式错误")
	}
	var downLoadList []map[string]string
	screenName := "--"
	startTime := time.Now().Unix()
	errorMsg := "--"
	errorCount := 0
	sleepTime := .5

	for {
		if startInt > endInt {
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
		resData, err := bliGetSpaceSearch(userID, tid, order, keyword, startInt)
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		for _, item := range resData.Data.List.Vlist {
			isVID := IsVideoID("bilibili", item.Bvid, runType.RedisConn)
			if isVID && runType.IsDeWeight {
				continue
			}
			downLoadList = append(downLoadList, map[string]string{
				"vid":   item.Bvid,
				"title": item.Title,
				"url":   fmt.Sprintf("https://www.bilibili.com/video/%s", item.Bvid),
			})
			screenName = item.Author
		}
		PrintInfof(fmt.Sprintf(
			"\rcurrent: %d gather: %d author: %s duration: %ds sleep：%.2fs",
			startInt, len(downLoadList), screenName, (time.Now().Unix() - startTime), sleepTime,
		))
		if errorMsg != "--" {
			color.Set(color.FgRed, color.Bold)
			fmt.Printf(" errCout：%d errMsg：%s", errorCount, errorMsg)
			color.Unset()
		}
		startInt++
	}
	fmt.Println("")
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))
	AnnieDownloadAll(downLoadList, runType, "bilibili")
	PrintInfo("全部下载完成")
	return nil
}

func bliGetUserID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://space\.bilibili\.com/(\d+)`),
	}
	for _, regxp := range regexps {
		resUserID := regxp.FindStringSubmatch(url)
		if len(resUserID) < 3 {
			continue
		}
		return resUserID[2], nil
	}
	return "", errors.New("获取UserID失败")
}

func bliGetTID(url string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://space\.bilibili\.com/\d+/video\?.*?tid=(\d+)`),
	}
	for _, regxp := range regexps {
		resTID := regxp.FindStringSubmatch(url)
		if len(resTID) < 3 {
			continue
		}
		return resTID[2]
	}
	return "0"
}
func bliGetOrder(url string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://space\.bilibili\.com/\d+/video\?.*?order=(.*?)(^|&.*?)`),
	}
	for _, regxp := range regexps {
		resOrder := regxp.FindStringSubmatch(url)
		if len(resOrder) < 4 {
			continue
		}
		return resOrder[2]
	}
	return "pubdate"
}
func bliGetKeyword(url string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://space\.bilibili\.com/\d+/video\?.*?keyword=(.*?)(^|&.*?)`),
	}
	for _, regxp := range regexps {
		resKeyword := regxp.FindStringSubmatch(url)
		if len(resKeyword) < 4 {
			continue
		}
		return resKeyword[2]
	}
	return ""
}

func bliGetSpaceSearch(userID, tid, order, keyword string, page int) (*BliUserVideoList, error) {
	api := fmt.Sprintf("https://api.bilibili.com/x/space/arc/search?mid=%s&ps=30&tid=%s&pn=%d&keyword=%s&order=%s&jsonp=jsonp", userID, tid, page, keyword, order)
	var jsonData BliUserVideoList
	if err := RequestGetJSON(api, map[string]string{
		"accept":     "application/json, text/plain, */*",
		"referer":    "https://space.bilibili.com/",
		"user-agent": UserAgentPc,
	}, &jsonData); err != nil {
		return nil, err
	}
	if jsonData.Code != 0 {
		return &jsonData, errors.New(jsonData.Message)
	}
	return &jsonData, nil
}

func bliGetMaxPage(userID, tid, order, keyword string) (int, int, error) {
	resData, err := bliGetSpaceSearch(userID, tid, order, keyword, 1)
	if err != nil {
		return 0, 0, err
	}
	return int(math.Ceil(float64(resData.Data.Page.Count) / 30)), resData.Data.Page.Count, nil
}
