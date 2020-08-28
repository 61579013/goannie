package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RunTxOne(runType RunType, arg map[string]string) error {
	err := AnnieDownload(runType.Url, runType.SavePath, runType.CookieFile, runType.DefaultCookie)
	if err != nil {
		return err
	}
	return nil
}

// 腾讯归档页：https://v.qq.com/detail/5/52852.html
func RunTxDetail(runType RunType, arg map[string]string) error {
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
	var id int

	resGetId := regexp.MustCompile(`^(http|https)://v\.qq\.com/detail/\d+/(\d+)\.html$`).FindStringSubmatch(runType.Url)
	if len(resGetId) < 3 {
		return errors.New("获取ID失败")
	}
	if id, err = strconv.Atoi(resGetId[2]); err != nil {
		return errors.New("获取ID失败")
	}

	for {
		if startYear > endYear || (startYear == endYear && startMonth > endMonth) {
			break
		}
		api := fmt.Sprintf(
			"https://s.video.qq.com/get_playsource?id=%d&plat=2&type=4&data_type=3&video_type=5&year=%d&month=%d&plname=qq&otype=json&_t=%d",
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
			isVID := IsVideoID("tengxun", item.ID)
			if isVID && runType.IsDeWeight {
				continue
			}
			downLoadList = append(downLoadList, map[string]string{
				"vid":item.ID,
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

// 腾讯作者页：https://v.qq.com/s/videoplus/1790091432#uin=42ffd591994e622dd2e414ecc3137397
func RunTxUserList(runType RunType, arg map[string]string) error {
	vuid, err := txGetVuid(runType.Url)
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
				isVID := IsVideoID("tengxun", item.Data.Vid)
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

// 看作者作品列表	look https://v.qq.com/s/videoplus/1790091432
func RunLookTxUserList(runType RunType, arg map[string]string) error {
	runType.Url = strings.Replace(runType.Url, "look ", "", -1)
	vuid, err := txGetVuid(runType.Url)
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

// 通过 url 获取 Vuid
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

// 请求腾讯作者作品列表API
func txGetVideoplusData(vuid string, offset, errorCount int) (*TengxunUserVideoList, error) {
	api := fmt.Sprintf("https://nodeyun.video.qq.com/x/api/videoplus/data?type=all&vuid=%s&last_vid_position=%d&offset=%d&index_context=score&_=%d",
		vuid, offset, offset, (time.Now().Unix() * 1000),
	)
	var jsonData TengxunUserVideoList
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return &jsonData, err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("referer", "https://v.qq.com/")
	req.Header.Set("user-agent", UserAgentPc)
	resP, err := Client.Do(req)
	if err != nil {
		return &jsonData, err
	}
	defer resP.Body.Close()
	if resP.StatusCode != 200 {
		if errorCount < 3 {
			return txGetVideoplusData(vuid, offset, errorCount+1)
		}
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
	if jsonData.ErrorMsg != "" {
		return &jsonData, errors.New(jsonData.ErrorMsg)
	}
	return &jsonData, nil

}

// 腾讯作者作品页数检查
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

// 腾讯归档请求API
func DetailGetPlaysource(url string) (*TengxunPlaysource, error) {
	var jsonData TengxunPlaysource
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &jsonData, err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("referer", "https://v.qq.com/")
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
