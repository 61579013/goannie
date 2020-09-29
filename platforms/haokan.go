package platforms

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// RunHkOne 单视频
func RunHkOne(runType RunType, arg map[string]string) error {
	vid := hkGetVID(runType.URL)
	// 判断是否过滤重复
	if vid != "" {
		isVID := IsVideoID("haokan", vid, runType.RedisConn)
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
		AddVideoID("haokan", vid, runType.RedisConn)
	}
	return nil
}

// hkGetVID 通过url获取vid
func hkGetVID(url string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://haokan\.baidu\.com/v\?vid=(\d+)($|\?.*?$)`),
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

// RunHkUserList 作者列表
func RunHkUserList(runType RunType, arg map[string]string) error {
	userID, err := hkGetUserID(runType.URL)
	if err != nil {
		return err
	}
	page, count, err := hkGetMaxPage(userID)
	if err != nil {
		return err
	}
	PrintInfo(fmt.Sprintf("\r总页数：%d  每页个数：%d  总个数：%d", page, 16, count))
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
	onPage := 1
	cTime := "0"
	for {
		if onPage > endInt {
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))

		resData, err := hkGetUserVideoList(userID, cTime)
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		cTime = resData.Data.Response.Ctime
		if onPage >= startInt {
			for _, item := range resData.Data.Response.Results {
				isVID := IsVideoID("haokan", item.Content.Vid, runType.RedisConn)
				if isVID && runType.IsDeWeight {
					continue
				}
				downLoadList = append(downLoadList, map[string]string{
					"vid":   item.Content.Vid,
					"title": item.Content.Title,
					"url":   item.Content.VideoShortURL,
				})
				screenName = item.Content.Author
			}
		}
		PrintInfof(fmt.Sprintf(
			"\rcurrent: %d gather: %d author: %s duration: %ds sleep：%.2fs",
			onPage, len(downLoadList), screenName, (time.Now().Unix() - startTime), sleepTime,
		))
		if errorMsg != "--" {
			color.Set(color.FgRed, color.Bold)
			fmt.Printf(" errCout：%d errMsg：%s", errorCount, errorMsg)
			color.Unset()
		}
		onPage++
	}
	fmt.Println("")
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))
	AnnieDownloadAll(downLoadList, runType, "haokan")
	PrintInfo("全部下载完成")
	return nil
}

func hkGetUserID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://haokan\.baidu\.com/author/(\d+)`),
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

func hkGetUserVideoList(userID, ctime string) (*HaokanUserVideoList, error) {
	api := fmt.Sprintf("https://haokan.baidu.com/author/%s?_format=json&rn=16&ctime=%s&_api=1", userID, ctime)
	var jsonData HaokanUserVideoList
	if err := RequestGetJSON(api, map[string]string{
		"accept":     "*/*",
		"referer":    "https://haokan.baidu.com/",
		"user-agent": UserAgentPc,
	}, &jsonData); err != nil {
		return nil, err
	}
	if jsonData.Errno != 0 {
		return &jsonData, errors.New(jsonData.Error)
	}
	return &jsonData, nil
}

func hkGetMaxPage(userID string) (int, int, error) {
	api := fmt.Sprintf("https://haokan.baidu.com/author/%s", userID)
	resP, err := RequestGet(api, map[string]string{
		"accept":     "*/*",
		"referer":    "https://haokan.baidu.com/",
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
	resCount := regexp.MustCompile(`<h3>视频</h3><p>(\d+)</p>`).FindStringSubmatch(string(body))
	if len(resCount) < 2 {
		return 0, 0, errors.New("获取分页失败")
	}
	count, err := strconv.Atoi(resCount[1])
	if err != nil {
		return 0, 0, err
	}
	return int(math.Ceil(float64(count) / 16)), count, nil
}
