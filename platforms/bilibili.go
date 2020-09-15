package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// RunBliOne 单视频
func RunBliOne(runType RunType, arg map[string]string) error {
	err := AnnieDownload(runType.URL, runType.SavePath, runType.CookieFile, runType.DefaultCookie)
	if err != nil {
		return err
	}
	return nil
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
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return &jsonData, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("referer", "https://space.bilibili.com/")
	req.Header.Set("user-agent", UserAgentPc)
	resP, err := Client.Do(req)
	if err != nil {
		return &jsonData, err
	}
	defer resP.Body.Close()
	if resP.StatusCode != 200 {
		return &jsonData, errors.New("请求失败")
	}
	body, err := ioutil.ReadAll(resP.Body)
	if err != nil {
		return &jsonData, err
	}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return &jsonData, err
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
