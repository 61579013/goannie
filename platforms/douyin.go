package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// RunDyShortURL 短链接
func RunDyShortURL(runType RunType, arg map[string]string) error {
	url, err := GetRealURL(runType.URL)
	if err != nil {
		return errors.New("获取SecUID失败")
	}
	runType.URL = url
	// 作者视频匹配
	regexpsUserList := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://www\.iesdouyin\.com/share/user/\d+\?sec_uid=(.*?)(&|$)`),
	}
	// 单视频匹配
	regexpsOne := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://www\.iesdouyin\.com/share/video/(\d+)`),
	}

	for _, regxp := range regexpsUserList {
		resItemID := regxp.FindStringSubmatch(url)
		if len(resItemID) < 3 {
			continue
		}
		return RunDyUserList(runType, arg)
	}
	for _, regxp := range regexpsOne {
		resItemID := regxp.FindStringSubmatch(url)
		if len(resItemID) < 3 {
			continue
		}
		return RunDyOne(runType, arg)
	}
	return nil
}

// RunDyUserList 作者视频
func RunDyUserList(runType RunType, arg map[string]string) error {
	// 获取 sec_uid
	secUID, err := DyGetSecUID(runType.URL)
	if err != nil {
		return err
	}
	// 获取视频作者信息
	userInfo, _ := DyGetUserInfo(secUID)
	PrintInfof(fmt.Sprintf("作者名称：%s 作品数量：%d 收藏数量：%d 粉丝数量：%d\n",
		userInfo.UserInfo.Nickname,
		userInfo.UserInfo.AwemeCount,
		userInfo.UserInfo.FavoritingCount,
		userInfo.UserInfo.FollowerCount,
	))
	var (
		start, end string
	)
	if err = GetCmdDataString("请输入起始数", &start); err != nil {
		return err
	}
	startInt, err := strconv.Atoi(start)
	if err != nil || startInt > userInfo.UserInfo.AwemeCount || startInt <= 0 {
		return errors.New("起始数格式错误")
	}
	if err = GetCmdDataString("请输入结束数", &end); err != nil {
		return err
	}
	endInt, err := strconv.Atoi(end)
	if err != nil || endInt > userInfo.UserInfo.AwemeCount || endInt < startInt || endInt <= 0 {
		return errors.New("结束数格式错误")
	}
	var downLoadList []map[string]string
	maxCursor := 0
	videoCount := 0
	for {
		if videoCount > endInt {
			break
		}
		PrintInfof(fmt.Sprintf(
			"\r作者名称：%s  起始：%d  结束：%d  获取到：%d  待下载：%d",
			userInfo.UserInfo.Nickname, startInt, endInt, videoCount, len(downLoadList),
		))
		res, err := DyGetUserlistAPI(secUID, maxCursor, 0)
		if err != nil {
			break
		}
		if res.HasMore == 0 {
			break
		}
		maxCursor = res.MaxCursor
		for _, item := range res.AwemeList {
			videoCount++
			VID := item.AwemeID
			isVID := IsVideoID("douyin", VID, runType.RedisConn)
			if (isVID && runType.IsDeWeight) || (videoCount > endInt || videoCount < startInt) {
				continue
			}
			downLoadList = append(downLoadList, map[string]string{
				"vid":   VID,
				"title": item.Desc,
				"url":   item.ShareURL,
			})
		}
	}
	fmt.Println("")
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))
	for _, video := range downLoadList {
		oneRunType := runType
		oneRunType.URL = video["url"]
		err := RunDyOne(oneRunType, map[string]string{})
		if err != nil {
			PrintErrInfo(err.Error())
		}
	}
	PrintInfo("全部下载完成")
	return nil
}

// DyGetSecUID 获取 sec_uid
func DyGetSecUID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://www\.iesdouyin\.com/share/user/\d+\?sec_uid=(.*?)&`),
	}
	for _, regxp := range regexps {
		resItemID := regxp.FindStringSubmatch(url)
		if len(resItemID) < 3 {
			continue
		}
		return resItemID[2], nil
	}
	// 没有在 url 中获取到 sec_uid，尝试通过请求获取
	url, err := GetRealURL(url)
	if err != nil {
		return "", errors.New("获取SecUID失败")
	}
	for _, regxp := range regexps {
		resItemID := regxp.FindStringSubmatch(url)
		if len(resItemID) < 3 {
			continue
		}
		return resItemID[2], nil
	}
	return "", errors.New("获取SecUID失败")
}

// DyGetUserlistAPI 请求抖音作者
func DyGetUserlistAPI(secUserID string, maxCursor, errCount int) (*DouyinVideoList, error) {
	api := "https://api3-core-c-lf.amemv.com/aweme/v1/aweme/post/?source=0&publish_video_strategy_type=2&max_cursor=%d&sec_user_id=%s&count=20&ts=%d&_rticket=%s&"
	ts := time.Now().Unix()
	rticket := fmt.Sprintf("%d", time.Now().UnixNano())[:13]

	var jsonData DouyinVideoList
	if err := RequestGetJSON(fmt.Sprintf(api, maxCursor, secUserID, ts, rticket), map[string]string{
		"Host":            "api3-core-c-lf.amemv.com",
		"Connection":      "keep-alive",
		"X-SS-REQ-TICKET": fmt.Sprintf("%d", time.Now().UnixNano())[:13],
		"sdk-version":     "1",
		"X-SS-DP":         "1128",
		"User-Agent":      "com.ss.android.ugc.aweme/100401 (Linux; U; Android 6.0.1; zh_CN; MI 5s; Build/V417IR; Cronet/TTNetVersion:3154e555 2020-03-04 QuicVersion:8fc8a2f3 2020-03-02)",
		"X-Khronos":       fmt.Sprintf("%d", ts),
	}, &jsonData); err != nil {
		if errCount < 1000 {
			return DyGetUserlistAPI(secUserID, maxCursor, errCount+1)
		}
		return nil, err
	}

	if len(jsonData.AwemeList) == 0 {
		if errCount < 1000 {
			time.Sleep(time.Microsecond * 500)
			return DyGetUserlistAPI(secUserID, maxCursor, errCount+1)
		}
		return &jsonData, nil
	}
	return &jsonData, nil
}

// DyGetUserInfo 获取作者信息
func DyGetUserInfo(secUserID string) (*DouyinUserInfo, error) {
	api := "https://www.iesdouyin.com/web/api/v2/user/info/?sec_uid=%s"
	request, err := http.NewRequest("GET", fmt.Sprintf(api, secUserID), nil)
	if err != nil {
		return &DouyinUserInfo{}, err
	}
	request.Header.Set("User-Agent", "com.ss.android.ugc.aweme/100401 (Linux; U; Android 6.0.1; zh_CN; MI 5s; Build/V417IR; Cronet/TTNetVersion:3154e555 2020-03-04 QuicVersion:8fc8a2f3 2020-03-02)")
	res, err := Client.Do(request)
	if err != nil {
		return &DouyinUserInfo{}, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &DouyinUserInfo{}, err
	}
	var jsonData DouyinUserInfo
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return &DouyinUserInfo{}, err
	}
	return &jsonData, nil
}

// RunDyOne 单视频
func RunDyOne(runType RunType, arg map[string]string) error {
	VideoInfo, err := DyGetOne(runType.URL)
	if err != nil {
		return err
	}
	// 判断是否过滤重复
	if VideoInfo.AwemeID != "" {
		isVID := IsVideoID("douyin", VideoInfo.AwemeID, runType.RedisConn)
		if isVID && runType.IsDeWeight {
			// 判断到了重复
			PrintErrInfo(RepetitionMsg)
			return nil
		}
	}
	downloadURL, err := DyGetDownloadURL(VideoInfo.URI, VideoInfo.Ratio)
	if err != nil {
		return err
	}
	dlpt := &DownloadPrint{
		"抖音视频 iesdouyin.com",
		VideoInfo.Title,
		"video",
		fmt.Sprintf("%dx%d", VideoInfo.Widht, VideoInfo.Height),
		"",
		0,
	}
	dlpt.Init(downloadURL)
	dlpt.Print()
	err = Aria2Download(downloadURL, runType.SavePath, fmt.Sprintf("%s.mp4", VideoInfo.Title), runType.CookieFile, 5)
	if err != nil {
		return err
	}
	// 下载完成记录
	AddVideoID("douyin", VideoInfo.AwemeID, runType.RedisConn)
	return nil

}

// DyGetOne 获取单视频信息
func DyGetOne(url string) (*DouyinVideoInfo, error) {
	itemID, err := DyGetItemID(url)
	if err != nil {
		return nil, err
	}
	videoItem, err := DyGetIteminfo(itemID)
	if err != nil {
		return nil, err
	}
	if len(videoItem.ItemList) == 0 {
		return nil, errors.New("请求获取itemData失败")
	}
	ItemData := videoItem.ItemList[0]
	var resData DouyinVideoInfo
	resData = DouyinVideoInfo{
		Title:    ItemData.Desc,
		Widht:    ItemData.Video.Widht,
		Height:   ItemData.Video.Height,
		Duration: ItemData.Video.Duration,
		Ratio:    ItemData.Video.Ratio,
		URI:      ItemData.Video.URI,
		AwemeID:  ItemData.AwemeID,
	}
	return &resData, nil
}

// DyGetItemID 获取ItemID
func DyGetItemID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://www\.iesdouyin\.com/share/video/(\d+)`),
	}
	for _, regxp := range regexps {
		resItemID := regxp.FindStringSubmatch(url)
		if len(resItemID) < 3 {
			continue
		}
		return resItemID[2], nil
	}
	// 没有在 url 中获取到 ItemID，尝试通过请求获取
	url, err := GetRealURL(url)
	if err != nil {
		return "", errors.New("获取ItemID失败")
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

// DyGetIteminfo 获取抖音视频信息
func DyGetIteminfo(itemID string) (*DouyinVideoItem, error) {
	url := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=%s", itemID)
	var jsonData DouyinVideoItem
	if err := RequestGetJSON(url, nil, &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}

// DyGetDownloadURL 通过uri获取真实下载地址
func DyGetDownloadURL(uri, ratio string) (string, error) {
	url := fmt.Sprintf("https://aweme.snssdk.com/aweme/v1/play/?video_id=%s&ratio=%s&line=0", uri, ratio)
	url, err := DyGetRealURL(url, 0)
	if err != nil {
		return "", err
	}
	return url, nil
}

// DyGetRealURL 获取跳转真实地址
func DyGetRealURL(url string, errorCount int) (string, error) {
	newClient := Client
	newClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Mobile Safari/537.36 Edg/85.0.564.63")
	// req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	// req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	// req.Header.Set("upgrade-insecure-requests", "1")
	// req.Header.Set("sec-fetch-dest", "document")
	// req.Header.Set("sec-fetch-mode", "navigate")
	// req.Header.Set("sec-fetch-site", "none")
	// req.Header.Set("sec-fetch-user", "?1")
	// req.Header.Set("accept-encoding", "gzip, deflate, br")
	resP, err := newClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resP.Body.Close()

	if resP.Header.Get("location") != "" {
		return resP.Header.Get("location"), nil
	}

	if errorCount < 1000 {
		time.Sleep(time.Microsecond * 1000)
		return DyGetRealURL(url, errorCount+1)
	}
	return "", errors.New("获取抖音下载地址失败")
}
