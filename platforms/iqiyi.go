package platforms

import (
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

const (
	IqyAuthCookie   = ""
	IqyDfp          = "a049cca57369584ec09b6d99c8fc03da96c5f68f2cbd09625f04b0ad7c15640026"
	IqyMDeviceID    = "go5e8hbplm0iumdxpsat9upx"
	IqyAgentversion = "10.7.5"
)

// RunIqyOne 单视频
func RunIqyOne(runType RunType, arg map[string]string) error {
	vid := IqyGetVID(runType.URL)
	// 判断是否过滤重复
	if vid != "" {
		isVID := IsVideoID("iqiyi", vid, runType.RedisConn)
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
		AddVideoID("iqiyi", vid, runType.RedisConn)
	}
	return nil
}

// IqyGetVID 通过请求获取vid
func IqyGetVID(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("referer", "https://www.iqiyi.com/")
	req.Header.Set("user-agent", UserAgentPc)
	resP, err := Client.Do(req)
	if err != nil {
		return ""
	}
	defer resP.Body.Close()
	if resP.StatusCode != 200 {
		return ""
	}
	body, err := ioutil.ReadAll(resP.Body)
	if err != nil {
		return ""
	}
	content := string(body)

	getVIDList := regexp.MustCompile(`param\['vid'\] = "(.*?)";`).FindStringSubmatch(content)
	if len(getVIDList) < 2 {
		return ""
	}
	return getVIDList[1]
}

// RunIqyUserList 作者视频
func RunIqyUserList(runType RunType, arg map[string]string) error {
	var err error
	fuid, err := IqyGetFuID(runType.URL)
	if err != nil {
		return err
	}
	page, count, err := IqyGetMaxPage(fuid)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("\r总页数：%d  每页个数：%d  总个数：%d", page, 28, count))
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
		resData, err := IqyGetUserVideos(fuid, fmt.Sprint(startInt))
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		for _, item := range resData.Data {
			isVID := IsVideoID("iqiyi", item.VID, runType.RedisConn)
			if isVID && runType.IsDeWeight {
				continue
			}
			downLoadList = append(downLoadList, map[string]string{
				"vid":   item.VID,
				"title": item.Title,
				"url":   item.PageUrl,
			})
			screenName = item.Nickname
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
	AnnieDownloadAll(downLoadList, runType, "iqiyi")
	PrintInfo("全部下载完成")
	return nil
}

// IqyGetFuID 通过url获取FuID
func IqyGetFuID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://www\.iqiyi\.com/u/(\d+)`),
	}
	for _, regxp := range regexps {
		resFuID := regxp.FindStringSubmatch(url)
		if len(resFuID) < 3 {
			continue
		}
		return resFuID[2], nil
	}
	return "", errors.New("获取FuID失败")

}

// IqyGetMaxPage 获取作者列表页数
func IqyGetMaxPage(fuid string) (int, int, error) {
	data, err := IqyVideosAction(fuid, "1", "1")
	if err != nil {
		return 0, 0, err
	}
	return int(math.Ceil(float64(data.Data.TotalNum / 28))), data.Data.TotalNum, nil
}

// IqyGetUserVideos 获取作者作品
func IqyGetUserVideos(fuid string, page string) (*IqiyEpisodeInfoAction, error) {
	var err error
	videos, err := IqyVideosAction(fuid, page, "28")
	if err != nil {
		return nil, err
	}
	qipuIds := ""
	for _, i := range videos.Data.Sort.Flows {
		qipuIds += fmt.Sprintf("%d,", i.QipuID)
	}
	if qipuIds != "" {
		qipuIds = qipuIds[:len(qipuIds)-1]
	}
	data, err := IqyEpisodeInfoAction(fuid, qipuIds)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// IqyEpisodeInfoAction 请求作品详情API
func IqyEpisodeInfoAction(fuid, qipuIds string) (*IqiyEpisodeInfoAction, error) {
	timestamp := fmt.Sprint(time.Now().Unix() * 1000)
	parameters := [][2]string{
		[2]string{"agenttype", "118"},
		[2]string{"agentversion", IqyAgentversion},
		[2]string{"authcookie", IqyAuthCookie},
		[2]string{"dfp", IqyDfp},
		[2]string{"m_device_id", IqyMDeviceID},
		[2]string{"qipuIds", qipuIds},
		[2]string{"sign", IqyEpisodeInfoActionSign(qipuIds, timestamp)},
		[2]string{"timestamp", timestamp},
	}
	url := fmt.Sprintf("https://iqiyihao.iqiyi.com/iqiyihao/episode_info.action?%s", MapToParameters(parameters))
	var jsonData IqiyEpisodeInfoAction
	if err := RequestGetJSON(url, map[string]string{
		"Referer":    fmt.Sprintf("https://www.iqiyi.com/u/%s/videos", fuid),
		"User-Agent": UserAgentPc,
	}, &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}

// IqyVideosAction 请求作品列表API
func IqyVideosAction(fuid, page, size string) (*IqiyGetVideosAction, error) {
	timestamp := fmt.Sprint(time.Now().Unix() * 1000)
	parameters := [][2]string{
		[2]string{"agenttype", "118"},
		[2]string{"agentversion", IqyAgentversion},
		[2]string{"authcookie", IqyAuthCookie},
		[2]string{"dfp", IqyDfp},
		[2]string{"fuid", fuid},
		[2]string{"m_device_id", IqyMDeviceID},
		[2]string{"page", page},
		[2]string{"size", size},
		[2]string{"sign", IqyVideosActionSign(fuid, timestamp, page, size)},
		[2]string{"timestamp", timestamp},
	}
	url := fmt.Sprintf("https://iqiyihao.iqiyi.com/iqiyihao/entity/get_videos.action?%s", MapToParameters(parameters))
	var jsonData IqiyGetVideosAction
	if err := RequestGetJSON(url, map[string]string{
		"Referer":    fmt.Sprintf("https://www.iqiyi.com/u/%s/videos", fuid),
		"User-Agent": UserAgentPc,
	}, &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}

// IqyVideosActionSign 生成作者作品列表API的Sign
func IqyVideosActionSign(fuid, timestamp, page, size string) string {
	parameters := [][2]string{
		[2]string{"agenttype", "118"},
		[2]string{"agentversion", IqyAgentversion},
		[2]string{"authcookie", IqyAuthCookie},
		[2]string{"dfp", IqyDfp},
		[2]string{"fuid", fuid},
		[2]string{"m_device_id", IqyMDeviceID},
		[2]string{"page", page},
		[2]string{"size", size},
		[2]string{"timestamp", timestamp},
	}
	n := MapToParameters(parameters)
	t := "GET"
	e := "iqiyihao.iqiyi.com/iqiyihao/entity/get_videos.action"
	i := fmt.Sprintf("%s%s?%s%s", t, e, n, "NZrFGv72GYppTUxO")
	return MD5(i)
}

// IqyEpisodeInfoActionSign 生成作者作品详细信息API的Sign
func IqyEpisodeInfoActionSign(qipuIds, timestamp string) string {
	parameters := [][2]string{
		[2]string{"agenttype", "118"},
		[2]string{"agentversion", IqyAgentversion},
		[2]string{"authcookie", IqyAuthCookie},
		[2]string{"dfp", IqyDfp},
		[2]string{"m_device_id", IqyMDeviceID},
		[2]string{"qipuIds", qipuIds},
		[2]string{"timestamp", timestamp},
	}
	n := NewMapToParameters(parameters, false)
	t := "GET"
	e := "iqiyihao.iqiyi.com/iqiyihao/episode_info.action"
	i := fmt.Sprintf("%s%s?%s%s", t, e, n, "NZrFGv72GYppTUxO")
	return MD5(i)
}

// RunIqyDetail 爱奇艺归档页：https://www.iqiyi.com/a_19rrht2ok5.html
func RunIqyDetail(runType RunType, arg map[string]string) error {
	var (
		start, end string
	)
	err := GetCmdDataString("请输入起始（2020）", &start)
	if err != nil {
		return err
	}
	err = GetCmdDataString("请输入结束（2020）", &end)
	if err != nil {
		return err
	}
	var (
		startYear, endYear int
	)
	if startYear, err = strconv.Atoi(start); err != nil {
		return errors.New("[起始]格式错误")
	}
	if endYear, err = strconv.Atoi(end); err != nil {
		return errors.New("[结束]格式错误")
	}
	if startYear > endYear {
		return errors.New("格式错误")
	}
	//var playPageInfo map[string]string

	playPageInfo, err := DetailGetPlayPageInfo(runType.URL, runType.CookieFile)
	if err != nil {
		return err
	}

	var downLoadList []map[string]string
	for {
		if startYear > endYear {
			break
		}
		url := fmt.Sprintf("https://pcw-api.iqiyi.com/album/source/svlistinfo?cid=%s&sourceid=%s&timelist=%d", playPageInfo.Cid, playPageInfo.AlbumId, startYear)

		resData, err := DetailGetSvlistinfo(url, runType.CookieFile)
		if err != nil {
			PrintErrInfo(err.Error())
			startYear++
			continue
		}
		if _, ok := resData.Data[fmt.Sprintf("%d", startYear)]; !ok {
			PrintErrInfo(fmt.Sprintf("%d：获取数据为空", startYear))
			startYear++
			continue
		}
		for _, item := range resData.Data[fmt.Sprintf("%d", startYear)] {
			isVID := IsVideoID("iqiyi", item.Vid, runType.RedisConn)
			if isVID && runType.IsDeWeight {
				continue
			}
			downLoadList = append(downLoadList, map[string]string{
				"vid":   item.Vid,
				"title": item.Name,
				"url":   item.PlayURL,
			})
		}
		startYear++
	}

	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))

	AnnieDownloadAll(downLoadList, runType, "iqiyi")

	PrintInfo("全部下载完成")
	return nil
}

// DetailGetSvlistinfo 爱奇艺归档请求API
func DetailGetSvlistinfo(url, cookiePath string) (*IqiyiSvlistinfo, error) {
	var jsonData IqiyiSvlistinfo
	if err := RequestGetJSON(url, map[string]string{
		"accept":     "*/*",
		"cookie":     GetTxtContent(cookiePath),
		"referer":    "https://www.iqiyi.com/",
		"user-agent": UserAgentPc,
	}, &jsonData); err != nil {
		return nil, err
	}
	if jsonData.Code != "A00000" {
		return &jsonData, errors.New("爱奇艺请求错误：" + jsonData.Code)
	}
	return &jsonData, nil
}

// DetailGetPlayPageInfo 爱奇艺归档请求HTML
func DetailGetPlayPageInfo(url, cookiePath string) (*IqiyiPlayPageInfo, error) {
	var jsonData IqiyiPlayPageInfo
	resP, err := RequestGet(url, map[string]string{
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"cookie":     GetTxtContent(cookiePath),
		"referer":    "https://www.iqiyi.com/",
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
	resDatas := regexp.MustCompile(`albumId: "(\w+)",(\s|[\r\n])+tvId: "\w+",(\s|[\r\n])+sourceId: \w+,(\s|[\r\n])+cid: "(\w+)",`).FindStringSubmatch(content)
	if len(resDatas) < 6 {
		return &jsonData, errors.New("获取 albumId,cid 失败")
	}
	jsonData.AlbumId = resDatas[1]
	jsonData.Cid = resDatas[5]
	return &jsonData, nil
}
