package platforms

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// RunYkOne 单视频
func RunYkOne(runType RunType, arg map[string]string) error {
	var err error
	vid, _ := YkGetVID(runType.URL)
	// 判断是否过滤重复
	if vid != "" {
		isVID := IsVideoID("youku", vid, runType.RedisConn)
		if isVID && runType.IsDeWeight {
			// 判断到了重复
			PrintErrInfo(RepetitionMsg)
			return nil
		}
	}
	if err = AnnieDownload(runType.URL, runType.SavePath, runType.CookieFile, runType.DefaultCookie); err != nil {
		return err
	}
	// 存储已下载
	if vid != "" {
		AddVideoID("youku", vid, runType.RedisConn)
	}
	return nil
}

// RunYkUserlist 作者
func RunYkUserlist(runType RunType, arg map[string]string) error {
	var err error
	uid, err := YkGetUID(runType.URL)
	if err != nil {
		return err
	}
	page, err := YkGetMaxPage(uid)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("\r总页数：%d  每页个数：%d ", page, 50))
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
	for {
		if startInt > endInt {
			break
		}
		resData, err := YkGetUserVideolist(uid, fmt.Sprint(startInt))
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		for _, item := range resData {
			VID, _ := YkGetVID(item["url"])
			isVID := IsVideoID("youku", VID, runType.RedisConn)
			if isVID && runType.IsDeWeight {
				continue
			}
			downLoadList = append(downLoadList, map[string]string{
				"vid":   VID,
				"title": item["title"],
				"url":   item["url"],
			})
			screenName = "--"
		}
		PrintInfof(fmt.Sprintf(
			"\rcurrent: %d gather: %d author: %s duration: %ds",
			startInt, len(downLoadList), screenName, (time.Now().Unix() - startTime),
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
	AnnieDownloadAll(downLoadList, runType, "youku")
	PrintInfo("全部下载完成")
	return nil
}

// YkGetVID 获取VID
func YkGetVID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://v\.youku\.com/v_show/id_(.*?).html`),
	}
	for _, regxp := range regexps {
		resUserID := regxp.FindStringSubmatch(url)
		if len(resUserID) < 3 {
			continue
		}
		return resUserID[2], nil
	}
	return "", nil
}

// YkGetUID 获取UID
func YkGetUID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://i\.youku\.com/(i|u)/(.*?)/`),
		regexp.MustCompile(`^(http|https)://i\.youku\.com/(i|u)/(.*?)\?`),
		regexp.MustCompile(`^(http|https)://i\.youku\.com/(i|u)/(.*?)$`),
	}
	for _, regxp := range regexps {
		resUserID := regxp.FindStringSubmatch(url)
		if len(resUserID) < 4 {
			continue
		}
		return resUserID[3], nil
	}
	// 有可能是独立链接 https://i.youku.com/yijialianmeng
	body, err := RequestGetHTML(url, nil)
	if err != nil {
		return "", err
	}
	getUID := regexp.MustCompile(`<a href="/i/(.*?)" title="主页">主页</a>`).FindSubmatch(body)
	if len(getUID) < 2 {
		return "", errors.New("获取UserID失败")
	}
	return string(getUID[1]), nil
}

// YkGetMaxPage 获取页数
func YkGetMaxPage(uid string) (int, error) {
	var err error
	body, err := YkGetUserVideolistHTML(uid, "1")
	if err != nil {
		return 0, err
	}
	getPage := regexp.MustCompile(`">(\d+)</a></li><li class="next"><a href=".*?">下一页</a>`).FindStringSubmatch(string(body))
	if len(getPage) > 2 {
		return 0, errors.New("获取页数失败")
	}
	page, err := strconv.Atoi(getPage[1])
	if err != nil {
		return 0, err
	}
	return page, nil
}

// YkGetUserVideolist 获取作者视频
func YkGetUserVideolist(uid, page string) ([]map[string]string, error) {
	var err error
	var retData []map[string]string
	body, err := YkGetUserVideolistHTML(uid, page)
	if err != nil {
		return retData, err
	}
	mustList := regexp.MustCompile(`<div class="v va">(.*?)</span></div></div></div>`).FindAllStringSubmatch(string(body), -1)
	for _, i := range mustList {
		if len(i) == 0 {
			continue
		}
		centent := regexp.MustCompile(`<a href="(.*?)".*?title="(.*?)">.*?</a>`).FindStringSubmatch(i[0])
		if len(centent) < 3 {
			continue
		}
		retData = append(retData, map[string]string{
			"title": centent[2],
			"url":   fmt.Sprintf("https:%s", centent[1]),
		})
	}
	return retData, nil
}

// YkGetUserVideolistHTML 请求作者列表返回HTML
func YkGetUserVideolistHTML(uid, page string) ([]byte, error) {
	var err error
	api := "https://i.youku.com/i/%s/videos?order=1&page=%s"
	body, err := RequestGetHTML(fmt.Sprintf(api, uid, page), map[string]string{
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"pragma":     "no-cache",
		"user-agent": UserAgentPc,
	})
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
