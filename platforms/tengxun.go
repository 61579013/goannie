package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// RunTxOne 单视频
func RunTxOne(runType RunType, arg map[string]string) error {
	vid := TxGetVID(runType.URL)
	// 判断是否过滤重复
	if vid != "" {
		isVID := IsVideoID("tengxun", vid, runType.RedisConn)
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
	if vid != "" {
		AddVideoID("tengxun", vid, runType.RedisConn)
	}
	return nil
}

// TxGetVID 通过url获取vid
func TxGetVID(url string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://v\.qq\.com/x/cover/.*?/(.*?)\.html($|\?.*?$)`),
		regexp.MustCompile(`^(http|https)://v\.qq\.com/x/cover/(.*?)\.html($|\?.*?$)`),
		regexp.MustCompile(`^(http|https)://v\.qq\.com/x/page/.*?/(.*?)\.html($|\?.*?$)`),
		regexp.MustCompile(`^(http|https)://v\.qq\.com/x/page/(.*?)\.html($|\?.*?$)`),
	}
	for _, regxp := range regexps {
		resVID := regxp.FindStringSubmatch(url)
		if len(resVID) < 3 {
			continue
		}
		return resVID[2]
	}
	return ""
}

// RunTxDetailTow 腾讯归档页 改版后
func RunTxDetailTow(runType RunType, arg map[string]string) error {
	urls, err := TxGetDetailTowURLS(runType.URL)
	if err != nil {
		return err
	}
	var downLoadList []map[string]string
	for _, url := range urls {
		vid := TxGetVID(url)
		isVID := IsVideoID("tengxun", vid, runType.RedisConn)
		if isVID && runType.IsDeWeight {
			continue
		}
		downLoadList = append(downLoadList, map[string]string{
			"vid":   vid,
			"title": "",
			"url":   url,
		})
	}
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))

	AnnieDownloadAll(downLoadList, runType, "tengxun")
	PrintInfo("全部下载完成")
	return nil
}

// TxGetDetailTowURLS 获取剧集链接
func TxGetDetailTowURLS(url string) ([]string, error) {
	content, err := RequestGetHTML(url, map[string]string{
		"referer":    "https://v.qq.com/",
		"user-agent": UserAgentPc,
	})
	if err != nil {
		return nil, err
	}
	pageStr := regexp.MustCompile(`mod_row_episode((.|\n)*?)<div r-component`).FindStringSubmatch(string(content))
	if len(pageStr) < 2 {
		return nil, errors.New("获取集数失败")
	}
	urls := regexp.MustCompile(`<a href="(https://v\.qq\.com/x/cover/.*?/.*?\.html)"`).FindAllStringSubmatch(pageStr[1], -1)
	var resData []string
	for _, i := range urls {
		if len(i) < 2 {
			continue
		}
		resData = append(resData, i[1])
	}
	return resData, nil
}

// RunTxDetail 腾讯归档页：https://v.qq.com/detail/5/52852.html
func RunTxDetail(runType RunType, arg map[string]string) error {
	data, err := RequestGetHTML(runType.URL, map[string]string{
		"accept":     "*/*",
		"referer":    "https://v.qq.com/",
		"user-agent": UserAgentPc,
	})
	if err != nil {
		return err
	}
	props := regexp.MustCompile(`(?s)r-component="p-index" r-props=".*?displayType: (\d+).*?" itemscope="`).FindStringSubmatch(string(data))
	if len(props) < 2 {
		return errors.New("基础信息解析失败")
	}
	displayType := props[1]

	if displayType == "8" {
		return TxDetailOne(runType)
	}

	return TxDetailTow(runType)
}

// TxDetailOne 腾讯剧集页第一种 https://v.qq.com/detail/8/87895.html
func TxDetailOne(runType RunType) error {
	var (
		start, end string
	)
	err := GetCmdDataString("请输入起始（2020-1）", &start)
	if err != nil {
		return err
	}
	err = GetCmdDataString("请输入结束（2020-1）", &end)
	if err != nil {
		return err
	}
	var (
		startYear, startMonth, endYear, endMonth int
	)

	startList := strings.Split(start, "-")
	endList := strings.Split(end, "-")
	if len(startList) < 2 {
		return errors.New("[起始]格式错误")
	}
	if len(endList) < 2 {
		return errors.New("[结束]格式错误")
	}
	if startYear, err = strconv.Atoi(startList[0]); err != nil {
		return errors.New("[起始]格式错误")
	}
	if startMonth, err = strconv.Atoi(startList[1]); err != nil {
		return errors.New("[起始]格式错误")
	}
	if endYear, err = strconv.Atoi(endList[0]); err != nil {
		return errors.New("[结束]格式错误")
	}
	if endMonth, err = strconv.Atoi(endList[1]); err != nil {
		return errors.New("[结束]格式错误")
	}
	if startYear > endYear || (startYear == endYear && startMonth > endMonth) {
		return errors.New("格式错误")
	}
	var apiList []string
	var id string

	resGetID := regexp.MustCompile(`^(http|https)://v\.qq\.com/detail/.*?/(.*?)\.html$`).FindStringSubmatch(runType.URL)
	if len(resGetID) < 3 {
		return errors.New("获取ID失败")
	}
	id = resGetID[2]
	for {
		if startYear > endYear || (startYear == endYear && startMonth > endMonth) {
			break
		}
		api := fmt.Sprintf(
			"https://s.video.qq.com/get_playsource?id=%s&plat=2&type=4&data_type=3&video_type=5&year=%d&month=%d&plname=qq&otype=json&_t=%d",
			id, startYear, startMonth, time.Now().Unix()*1000,
		)
		apiList = append(apiList, api)
		if startMonth >= 12 {
			startYear++
			startMonth = 1
		} else {
			startMonth++
		}
	}

	var downLoadList []map[string]string
	for _, url := range apiList {
		resData, err := DetailGetPlaysource(url)
		if err != nil {
			PrintErrInfo(err.Error())
			continue
		}
		for _, item := range resData.PlaylistItem.VideoPlayList {
			isVID := IsVideoID("tengxun", item.ID, runType.RedisConn)
			if isVID && runType.IsDeWeight {
				continue
			}
			downLoadList = append(downLoadList, map[string]string{
				"vid":   item.ID,
				"title": item.Title,
				"url":   item.PlayURL,
			})
		}
	}
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))

	AnnieDownloadAll(downLoadList, runType, "tengxun")

	PrintInfo("全部下载完成")
	return nil
}

// TxDetailTow 腾讯剧集页第二种 https://v.qq.com/detail/m/m441e3rjq9kwpsc.html
func TxDetailTow(runType RunType) error {
	var id string
	resGetID := regexp.MustCompile(`^(http|https)://v\.qq\.com/detail/.*?/(.*?)\.html$`).FindStringSubmatch(runType.URL)
	if len(resGetID) < 3 {
		return errors.New("获取ID失败")
	}
	id = resGetID[2]
	jsonData, err := txGetPlaysource(id)
	if err != nil {
		return err
	}
	videoCount := len(jsonData.PlaylistItem.VideoPlayList)
	PrintInfo(fmt.Sprintf("\r总个数：%d ", videoCount))

	var (
		start, end string
	)
	if err = GetCmdDataString("请输入起始数", &start); err != nil {
		return err
	}
	startInt, err := strconv.Atoi(start)
	if err != nil || startInt > videoCount || startInt <= 0 {
		return errors.New("起始数格式错误")
	}
	if err = GetCmdDataString("请输入结束数", &end); err != nil {
		return err
	}
	endInt, err := strconv.Atoi(end)
	if err != nil || endInt > videoCount || endInt < startInt || endInt <= 0 {
		return errors.New("结束数格式错误")
	}
	var downLoadList []map[string]string
	for {
		if startInt > endInt {
			break
		}
		item := jsonData.PlaylistItem.VideoPlayList[startInt-1]
		downLoadList = append(downLoadList, map[string]string{
			"vid":   item.ID,
			"title": item.Title,
			"url":   fmt.Sprintf("https://v.qq.com/x/page/%s.html", item.ID),
		})
		startInt++
	}
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))
	AnnieDownloadAll(downLoadList, runType, "tengxun")

	PrintInfo("全部下载完成")
	return nil
}

// RunTxUserList 腾讯作者页：https://v.qq.com/s/videoplus/1790091432#uin=42ffd591994e622dd2e414ecc3137397
func RunTxUserList(runType RunType, arg map[string]string) error {
	vuid, err := txGetVuid(runType.URL)
	if err != nil {
		return err
	}
	page, count, err := txGetUserListPage(vuid, 0, 20)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("\r总页数：%d  每页个数：%d  总个数：%d", page, 10, count))
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
		resData, err := txGetVideoplusData(vuid, (startInt-1)*10, 0)
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		if len(resData.Body.Modules) > 0 && len(resData.Body.Modules[0].Sections) > 0 {
			for _, item := range resData.Body.Modules[0].Sections[0].BlockList.Blocks {
				isVID := IsVideoID("tengxun", item.Data.Vid, runType.RedisConn)
				if isVID && runType.IsDeWeight {
					continue
				}
				downLoadList = append(downLoadList, map[string]string{
					"vid":   item.Data.Vid,
					"title": item.Data.Title,
					"url":   fmt.Sprintf("https://v.qq.com/x/page/%s.html", item.Data.Vid),
				})
				screenName = item.Data.CardInfo.UserInfo.UserName
			}
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
	AnnieDownloadAll(downLoadList, runType, "tengxun")

	PrintInfo("全部下载完成")
	return nil
}

// RunLookTxUserList 看作者作品列表	look https://v.qq.com/s/videoplus/1790091432
func RunLookTxUserList(runType RunType, arg map[string]string) error {
	runType.URL = strings.Replace(runType.URL, "look ", "", -1)
	vuid, err := txGetVuid(runType.URL)
	if err != nil {
		return err
	}
	page, count, err := txGetUserListPage(vuid, 0, 20)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("\r总页数：%d  每页个数：%d  总个数：%d", page, 10, count))
	startInt := 1
	endInt := page
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
		resData, err := txGetVideoplusData(vuid, (startInt-1)*10, 0)
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			break
		}
		if len(resData.Body.Modules) > 0 && len(resData.Body.Modules[0].Sections) > 0 {
			if len(resData.Body.Modules[0].Sections[0].BlockList.Blocks) > 0 {
				downLoadList = append(downLoadList, map[string]string{
					"page":  fmt.Sprintf("%d", startInt),
					"title": resData.Body.Modules[0].Sections[0].BlockList.Blocks[0].Data.Title,
				})
				screenName = resData.Body.Modules[0].Sections[0].BlockList.Blocks[0].Data.CardInfo.UserInfo.UserName
			}
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
	for _, v := range downLoadList {
		PrintInfo(fmt.Sprintf("第 %s 页 %s", v["page"], v["title"]))
	}
	return nil
}

// txGetVuid 通过 url 获取 Vuid
func txGetVuid(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://haokan\.baidu\.com/author/(\d+)`),
		regexp.MustCompile(`^(http|https)://v\.qq\.com/s/videoplus/(\d+)`),
		regexp.MustCompile(`^(http|https)://v\.qq\.com/x/bu/h5_user_center\?vuid=(\d+)`),
	}
	for _, regxp := range regexps {
		resVuid := regxp.FindStringSubmatch(url)
		if len(resVuid) < 3 {
			continue
		}
		return resVuid[2], nil
	}
	return "", errors.New("获取Vuid失败")
}

// txGetPlaysource 获取视频播放列表
func txGetPlaysource(id string) (*TxGetPlaysource, error) {
	api := fmt.Sprintf("https://s.video.qq.com/get_playsource?id=%s&data_type=2&type=4&range=1-100000&otype=json&num_mod_cnt=20&callback=Replacejsonp", id)
	body, err := RequestGetHTML(api, map[string]string{
		"accept":     "*/*",
		"referer":    "https://v.qq.com/",
		"user-agent": UserAgentPc,
	})
	if err != nil {
		return nil, err
	}
	newBody := strings.ReplaceAll(string(body), "Replacejsonp(", "")
	if len(newBody) > 0 {
		newBody = newBody[:len(newBody)-1]
	}
	var jsonData TxGetPlaysource
	err = json.Unmarshal([]byte(newBody), &jsonData)
	if err != nil {
		return nil, err
	}
	if jsonData.Error != 0 {
		return nil, errors.New(jsonData.Msg)
	}
	return &jsonData, nil
}

// txGetVideoplusData 请求腾讯作者作品列表API
func txGetVideoplusData(vuid string, offset, errorCount int) (*TengxunUserVideoList, error) {
	api := fmt.Sprintf("https://nodeyun.video.qq.com/x/api/videoplus/data?type=all&vuid=%s&last_vid_position=%d&offset=%d&index_context=score&_=%d",
		vuid, offset, offset, (time.Now().Unix() * 1000),
	)
	var jsonData TengxunUserVideoList
	if err := RequestGetJSON(api, map[string]string{
		"accept":     "*/*",
		"referer":    "https://v.qq.com/",
		"user-agent": UserAgentPc,
	}, &jsonData); err != nil {
		if errorCount < 3 {
			return txGetVideoplusData(vuid, offset, errorCount+1)
		}
		return nil, err
	}
	if jsonData.ErrorMsg != "" {
		return &jsonData, errors.New(jsonData.ErrorMsg)
	}
	return &jsonData, nil
}

// txGetUserListPage 腾讯作者作品页数检查
func txGetUserListPage(vuid string, s, b int) (int, int, error) {
	resData, err := txGetVideoplusData(vuid, s*10, 0)
	if err != nil {
		count := s * 10
		return s + 1, count, nil
	}
	PrintInfof(fmt.Sprintf("\r 检查页数----> 当前页 %d", s+1))
	if resData.Body.HasNextPage {
		return txGetUserListPage(vuid, s+b, b)
	}
	if b == 1 {
		endPageCount := 0
		if len(resData.Body.Modules) > 0 && len(resData.Body.Modules[0].Sections) > 0 {
			endPageCount = len(resData.Body.Modules[0].Sections[0].BlockList.Blocks)
		}
		count := (s * 10) + endPageCount
		return s + 1, count, nil
	}
	return txGetUserListPage(vuid, s-b, 1)
}

// DetailGetPlaysource 腾讯归档请求API
func DetailGetPlaysource(url string) (*TengxunPlaysource, error) {
	var jsonData TengxunPlaysource
	resP, err := RequestGet(url, map[string]string{
		"accept":     "*/*",
		"referer":    "https://v.qq.com/",
		"user-agent": UserAgentPc,
	})
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
	content := string(body)
	content = strings.Replace(content, "QZOutputJson=", "", -1)
	content = content[:len(content)-1]
	err = json.Unmarshal([]byte(content), &jsonData)
	if err != nil {
		return &jsonData, err
	}
	if jsonData.Error != 0 {
		return &jsonData, errors.New(jsonData.Msg)
	}
	return &jsonData, nil
}
