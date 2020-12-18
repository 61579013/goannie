package platforms

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// RunXgUserList TA的视频 https://www.ixigua.com/home/85383446500/video/
func RunXgUserList(runType RunType, arg map[string]string) error {
	page, count, err := xgGetUserListPage(runType.URL)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("总页数：%d  每页个数：%d  总个数：%d", page, 30, count))
	var (
		start, end, maxConnection string
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
	if err = GetCmdDataString("请输入最大连接数", &maxConnection); err != nil {
		return err
	}
	if _, err = strconv.Atoi(maxConnection); err != nil {
		return errors.New("最大连接数格式错误")
	}
	sleepTime := .5
	resUserID := regexp.MustCompile(`^(http|https)://www\.ixigua\.com/home/(\d+)`).FindStringSubmatch(runType.URL)
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
	breakCount := 0
	for {
		if onPage > endInt {
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
		PrintInfof(fmt.Sprintf(
			"\rcurrent: %d gather: %d author: %s duration: %ds sleep：%.2fs",
			onPage, len(downLoadList), screenName, (time.Now().Unix() - startTime), sleepTime,
		))
		if errorMsg != "--" {
			color.Set(color.FgRed, color.Bold)
			fmt.Printf(" errCout：%d errMsg：%s", errorCount, errorMsg)
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
				isVID := IsVideoID("xigua", item.GroupIdStr, runType.RedisConn)
				if isVID && runType.IsDeWeight {
					continue
				}
				downLoadList = append(downLoadList, map[string]string{
					"vid":   item.GroupIdStr,
					"title": item.Title,
					"url":   item.ArticleURL,
				})
			}
		}
		if !jsonData.HasMore && len(jsonData.Data) == 0 {
			// break 机会
			if breakCount < 10 {
				breakCount++
				continue
			}
			break
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
		oneRunType.URL = video["url"]
		err := XgNewOne(oneRunType, map[string]string{})
		if err != nil {
			PrintErrInfo(err.Error())
		} else {
			AddVideoID("xigua", video["vid"], runType.RedisConn)
		}
	}

	PrintInfo("全部下载完成")

	return nil
}

// RunLookXgUserList 看作者作品列表 look https://www.ixigua.com/home/85383446500/video/
func RunLookXgUserList(runType RunType, arg map[string]string) error {
	runType.URL = strings.Replace(runType.URL, "look ", "", -1)
	page, count, err := xgGetUserListPage(runType.URL)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("总页数：%d  每页个数：%d  总个数：%d", page, 30, count))
	sleepTime := .5
	resUserID := regexp.MustCompile(`^(http|https)://www\.ixigua\.com/home/(\d+)`).FindStringSubmatch(runType.URL)
	if len(resUserID) < 3 {
		return errors.New("西瓜获取UserID失败")
	}
	userID := resUserID[2]
	startInt := 1
	endInt := page
	var downLoadList []map[string]string
	onPage := 1
	maxBehotTime := 0
	screenName := "--"
	startTime := time.Now().Unix()
	errorMsg := "--"
	errorCount := 0
	breakCount := 0
	for {
		if onPage > endInt {
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
		PrintInfof(fmt.Sprintf(
			"\rcurrent: %d gather: %d author: %s duration: %ds sleep：%.2fs",
			onPage, len(downLoadList), screenName, (time.Now().Unix() - startTime), sleepTime,
		))
		if errorMsg != "--" {
			color.Set(color.FgRed, color.Bold)
			fmt.Printf(" errCout：%d errMsg：%s", errorCount, errorMsg)
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
			if len(jsonData.Data) > 0 {
				downLoadList = append(downLoadList, map[string]string{
					"title": jsonData.Data[0].Title,
					"page":  fmt.Sprintf("%d", onPage),
				})
			}
		}
		// 判断连续致命错误
		if !jsonData.HasMore && len(jsonData.Data) == 0 {
			// break 机会
			if breakCount < 10 {
				breakCount++
				continue
			}
			break
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
	for _, v := range downLoadList {
		PrintInfo(fmt.Sprintf("第 %s 页 %s", v["page"], v["title"]))
	}
	return nil
}

// xgGetUserListHome 请求作者作品列表api
func xgGetUserListHome(url, cookiePath string, errorCount int) (*XiguaUserList, error) {
	var jsonData XiguaUserList
	if err := RequestGetJSON(url, map[string]string{
		"accept":     "application/json",
		"cookie":     GetTxtContent(cookiePath),
		"referer":    url,
		"user-agent": UserAgentWap,
	}, &jsonData); err != nil {
		if errorCount < 3 {
			return xgGetUserListHome(url, cookiePath, errorCount+1)
		}
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
	resP, err := RequestGet(url, map[string]string{
		"accept":     "*/*",
		"referer":    url,
		"user-agent": UserAgentPc,
	})
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

// RunXgOne 西瓜单视频 https://www.ixigua.com/6832194590221533707
func RunXgOne(runType RunType, arg map[string]string) error {
	itemID, err := xgGetItemID(runType.URL)
	if err != nil {
		return err
	}
	title, videoID, err := xgGetVideoID(itemID, runType.CookieFile)
	if err != nil {
		return err
	}
	downloadURL, err := xgDownloadUrl(itemID, videoID)
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
	dlpt.Init(downloadURL)
	dlpt.Print()
	maxConnectionPerServer := 1
	if _, ok := arg["maxConnectionPerServer"]; ok {
		maxConnectionPerServer, err = strconv.Atoi(arg["maxConnectionPerServer"])
		if err != nil {
			return err
		}
	}
	err = Aria2Download(downloadURL, runType.SavePath, fmt.Sprintf("%s.mp4", title), runType.CookieFile, maxConnectionPerServer)
	if err != nil {
		return err
	}
	return nil
}

// XgNewOne 新的获取单视频
func XgNewOne(runType RunType, arg map[string]string) error {
	data, err := xgGetSsrHydratedData(runType.URL)
	if err != nil {
		return err
	}
	downloadURL := xgGetNewDownloadUrl(data)
	title := data.AnyVideo.GidInformation.PackerData.Video.Title
	dlpt := &DownloadPrint{
		"西瓜视频 ixigua.com",
		title,
		"video",
		"normal",
		"",
		0,
	}
	dlpt.Init(downloadURL)
	dlpt.Print()
	maxConnectionPerServer := 1
	if _, ok := arg["maxConnectionPerServer"]; ok {
		maxConnectionPerServer, err = strconv.Atoi(arg["maxConnectionPerServer"])
		if err != nil {
			return err
		}
	}
	err = Aria2Download(downloadURL, runType.SavePath, fmt.Sprintf("%s.mp4", title), runType.CookieFile, maxConnectionPerServer)
	if err != nil {
		return err
	}
	return nil
}

func xgGetNewDownloadUrl(d *XiguaSsrHydratedData) string {
	videoList := d.AnyVideo.GidInformation.PackerData.Video.VideoResource.Dash120Fps.DynamicVideo.DynamicVideoList
	if len(videoList) != 0 && videoList[len(videoList)-1].MainURL != "" {
		decoded, _ := base64.StdEncoding.DecodeString(videoList[len(videoList)-1].MainURL)
		return string(decoded)
	}
	vl := d.AnyVideo.GidInformation.PackerData.Video.VideoResource.Normal.VideoList
	var murl string
	if vl.Video4.MainURL != "" {
		murl = vl.Video4.MainURL
	} else if vl.Video3.MainURL != "" {
		murl = vl.Video3.MainURL
	} else if vl.Video2.MainURL != "" {
		murl = vl.Video2.MainURL
	} else if vl.Video1.MainURL != "" {
		murl = vl.Video1.MainURL
	}
	r, _ := base64.StdEncoding.DecodeString(murl)
	return string(r)
}

func xgGetSsrHydratedData(url string) (*XiguaSsrHydratedData, error) {
	data, err := RequestGetHTML(url, map[string]string{
		"Connection":                "keep-alive",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"User-Agent":                UserAgentPc,
		"Referer":                   "https://www.ixigua.com/",
		"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"Cookie":                    "xiguavideopcwebid=6872983880459503118; xiguavideopcwebid.sig=B4DvNwwGGQ-hDxYcJo5FfbMEIn4; _ga=GA1.2.572711536.1600241266; MONITOR_WEB_ID=bfe0e43a-e004-400e-8040-81677a199a22; ttwid=1%7CPWHvUSGTtsxK0WUzkuq7SxJtT7L3WHRvJeSGG5WZjiw%7C1604995289%7Cec6a591ac986362929a9be173d65df8f3551269fff0694d34a5e57a33cd287eb; ixigua-a-s=1; Hm_lvt_db8ae92f7b33b6596893cdf8c004a1a2=1608261601; _gid=GA1.2.1203395873.1608261601; Hm_lpvt_db8ae92f7b33b6596893cdf8c004a1a2=1608262109",
		"Upgrade-Insecure-Requests": "1",
	})
	if err != nil {
		return nil, err
	}
	jsonStrFind := regexp.MustCompile(`window\._SSR_HYDRATED_DATA=(.*?)</script>`).FindSubmatch(data)
	if len(jsonStrFind) < 2 {
		return nil, errors.New("解析数据失败")
	}
	jsonStr := strings.ReplaceAll(string(jsonStrFind[1]), ":undefined", ":\"undefined\"")
	var jsonData XiguaSsrHydratedData
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
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
	if err := RequestGetJSON(url, map[string]string{
		"accept":     "*/*",
		"referer":    url,
		"user-agent": UserAgentPc,
	}, &jsonData); err != nil {
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
