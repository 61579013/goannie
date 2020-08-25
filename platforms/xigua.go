package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/rock_rabbit/goannie/godler"
	"github.com/fatih/color"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// TA的视频 https://www.ixigua.com/home/85383446500/video/
func RunXgUserList(runType RunType) error {
	page, count, err := xgGetUserListPage(runType.Url)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("总页数：%d  每页个数：%d  总个数：%d", page, 30, count))
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
	//if err = GetCmdDataString("请输入间隔时间（s）", &sleep); err != nil {
	//	return err
	//}
	//sleepTime, err := strconv.ParseFloat(sleep, 2)
	//if err != nil {
	//	return errors.New("间隔时间格式错误")
	//}
	//startInt := 1
	//endInt := 1
	sleepTime := .5
	resUserID := regexp.MustCompile("^(http|https)://www\\.ixigua\\.com/home/(\\d+)/").FindStringSubmatch(runType.Url)
	if len(resUserID) < 3 {
		return errors.New("西瓜获取UserID失败")
	}
	userID := resUserID[2]
	var downLoadList []map[string]string
	onPage := 1
	maxBehotTime := 0
	screenName := "--"
	startTime := time.Now().Unix()
	errorMsg := "--"
	errorCount := 0
	for {
		if onPage > endInt {
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
		PrintInfof(fmt.Sprintf(
			"\r起始页：%d 结束页：%d 当前页：%d 已采集：%d个 作者名称：%s 耗时：%d秒 时间间隔：%.2f秒",
			startInt, endInt, onPage, len(downLoadList), screenName, (time.Now().Unix() - startTime), sleepTime,
		))
		if errorMsg != "--" {
			color.Set(color.FgRed, color.Bold)
			fmt.Printf(" 错误次数：%d 错误信息：%s", errorCount, errorMsg)
			color.Unset()
		}

		apiUrl := fmt.Sprintf("https://m.ixigua.com/video/app/user/home/?to_user_id=%s&format=json&preActiveKey=home&max_behot_time=%d", userID, maxBehotTime)
		jsonData, err := xgGetUserListHome(apiUrl, runType.CookieFile, 0)
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		screenName = jsonData.UserInfo.Name
		if onPage >= startInt {
			for _, item := range jsonData.Data {
				downLoadList = append(downLoadList, map[string]string{
					"title": item.Title,
					"url":   item.ArticleURL,
				})
			}
		}
		if len(jsonData.Data) == 0 {
			errorCount++
			errorMsg = "获取信息错误"
			continue
		}
		maxBehotTime = jsonData.Data[len(jsonData.Data)-1 : len(jsonData.Data)][0].BehotTime
		onPage++
	}
	fmt.Println("")
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))

	for _, video := range downLoadList {
		oneRunType := runType
		oneRunType.Url = video["url"]
		err := RunXgOne(oneRunType)
		if err !=nil{
			PrintErrInfo(err.Error())
		}
	}

	PrintInfo("全部下载完成")

	return nil
}

func xgGetUserListHome(url, cookiePath string, errorCount int) (*XiguaUserList, error) {
	var jsonData XiguaUserList
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if errorCount < 3 {
			return xgGetUserListHome(url, cookiePath, errorCount+1)
		}
		return &jsonData, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("cookie", GetTxtContent(cookiePath))
	req.Header.Set("referer", url)
	req.Header.Set("user-agent", UserAgentWap)
	resP, err := Client.Do(req)
	if err != nil {
		if errorCount < 3 {
			return xgGetUserListHome(url, cookiePath, errorCount+1)
		}
		return &jsonData, err
	}
	defer resP.Body.Close()
	if resP.StatusCode != 200 {
		if errorCount < 3 {
			return xgGetUserListHome(url, cookiePath, errorCount+1)
		}
		return &jsonData, errors.New("请求失败")
	}
	body, err := ioutil.ReadAll(resP.Body)
	if err != nil {
		if errorCount < 3 {
			return xgGetUserListHome(url, cookiePath, errorCount+1)
		}
		return &jsonData, err
	}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return &jsonData, errors.New("大概是速度太快")
	}
	if jsonData.Message != "success" {
		if errorCount < 3 {
			return xgGetUserListHome(url, cookiePath, errorCount+1)
		}
		return &jsonData, errors.New("西瓜视频请求失败：" + jsonData.Message)
	}
	return &jsonData, nil
}

// 获取TA的视频总页数
func xgGetUserListPage(url string) (int, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("accept", "*/*")
	//req.Header.Set("cookie", GetTxtContent(cookiePath))
	req.Header.Set("referer", url)
	req.Header.Set("user-agent", UserAgentPc)
	resP, err := Client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resP.Body.Close()
	if resP.StatusCode != 200 {
		return 0, 0, errors.New("请求失败")
	}
	body, err := ioutil.ReadAll(resP.Body)
	resPage := regexp.MustCompile("<div class=\"count\">(\\d+)</div>").FindStringSubmatch(string(body))
	if len(resPage) < 2 {
		return 0, 0, errors.New("西瓜获取页数失败")
	}
	count, err := strconv.Atoi(resPage[1])
	if err != nil {
		return 0, 0, errors.New("西瓜获取页数失败")
	}
	return int(math.Ceil(float64(count) / 30)), count, nil
}

// 西瓜单视频 https://www.ixigua.com/6832194590221533707
func RunXgOne(runType RunType) error {
	itemID, err := xgGetItemID(runType.Url)
	if err != nil {
		return err
	}
	title, videoID, err := xgGetVideoID(itemID, runType.CookieFile)
	if err != nil {
		return err
	}
	downloadUrl, err := xgDownloadUrl(itemID, videoID)
	if err != nil {
		return err
	}
	dlpt := &DownloadPrint{
		"西瓜视频 ixigua.com",
		title,
		"video",
		"normal",
		"",
		0,
	}
	dlpt.Init(downloadUrl)
	dlpt.Print()
	dler := &godler.DlerUrl{
		Url:      downloadUrl,
		SavePath: fmt.Sprintf("%s\\%s.mp4", runType.SavePath, title),
		IsBar:    true,
		Client:   http.Client{Timeout: time.Second * 180},
		TimeOut:  time.Second * 180,
		OneThreading: true,
	}
	err = dler.Download()
	if err != nil {
		return err
	}
	return nil
}

// 通过 url 获取 ItemID
func xgGetItemID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://www\.ixigua\.com/(\d+)`),
		regexp.MustCompile(`^(http|https)://m\.ixigua\.com/(\d+)`),
		regexp.MustCompile(`^(http|https)://www\.ixigua\.com/.*?\?id=(\d+)`),
		regexp.MustCompile(`^(http|https)://m\.ixigua\.com/video/(\d+)`),
		regexp.MustCompile(`^(http|https)://toutiao\.com/group/(\d+)/`),
	}
	for _, regxp := range regexps {
		resItemID := regxp.FindStringSubmatch(url)
		if len(resItemID) < 3 {
			continue
		}
		return resItemID[2], nil
	}
	return "", errors.New("获取ItemID失败")
}

// 请求获取 标题 和 VideoID
func xgGetVideoID(itemID, cookiePath string) (string, string, error) {
	var jsonData XiguaInfo
	url := fmt.Sprintf("https://m.365yg.com/i%s/info/", itemID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("accept", "*/*")
	//req.Header.Set("cookie", GetTxtContent(cookiePath))
	req.Header.Set("referer", url)
	req.Header.Set("user-agent", UserAgentPc)
	resP, err := Client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resP.Body.Close()
	if resP.StatusCode != 200 {
		return "", "", errors.New("请求失败")
	}
	body, err := ioutil.ReadAll(resP.Body)
	if err != nil {
		return "", "", err
	}
	content := string(body)
	err = json.Unmarshal([]byte(content), &jsonData)
	if err != nil {
		return "", "", err
	}
	if !jsonData.Success {
		return "", "", errors.New("西瓜视频info请求错误")
	}
	return jsonData.Data.Title, jsonData.Data.VideoID, nil
}

// 获取下载地址
func xgDownloadUrl(itemID, videoID string) (string, error) {
	newClient := Client
	newClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	url := fmt.Sprintf("https://api.huoshan.com/hotsoon/item/video/_source/?video_id=%s&line=0&app_id=0&vquality=normal&watermark=0&sf=5&item_id=%s",
		videoID, itemID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", UserAgentPc)
	resP, err := newClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resP.Body.Close()
	if resP.StatusCode == 301 || resP.StatusCode == 302 {
		return resP.Header.Get("location"), nil
	}
	return "", errors.New("获取西瓜视频下载路径失败")
}
