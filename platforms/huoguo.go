package platforms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 作者视频 https://huoguo.qq.com/m/person.html?userid=18590596
func RunHgUserList(runType RunType, arg map[string]string) error {
	userid, err := hgGetUserID(runType.Url)
	if err != nil {
		return err
	}
	resData, err := hgGetUserWorkList(userid, 0, "")
	if err != nil {
		return err
	}
	count := resData.Data.Count
	page := int(math.Ceil(float64(count) / 20))
	PrintInfo(fmt.Sprintf("总页数：%d  每页个数：%d  总个数：%d", page, 20, count))
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
	pageContext := ""
	onPage := 1
	for {
		if onPage > endInt {
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
		resData, err := hgGetUserWorkList(userid, (onPage-1)*10, pageContext)
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		pageContext = resData.Data.PageContext
		if onPage >= startInt {
			for _, item := range resData.Data.Collections {
				isVID := IsVideoID("tengxun", item.TvBoard.VideoData.Vid, runType.RedisConn)
				if isVID && runType.IsDeWeight {
					continue
				}
				downLoadList = append(downLoadList, map[string]string{
					"vid":   item.TvBoard.VideoData.Vid,
					"title": item.TvBoard.VideoData.Title,
					"url":   fmt.Sprintf("https://v.qq.com/x/page/%s.html", item.TvBoard.VideoData.Vid),
				})
				screenName = item.TvBoard.User.UserInfo.UserName
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
	AnnieDownloadAll(downLoadList, runType, "tengxun")

	PrintInfo("全部下载完成")
	return nil
}

// 看作者作品列表 look https://huoguo.qq.com/m/person.html?userid=18590596
func RunLookHgUserList(runType RunType, arg map[string]string) error {
	runType.Url = strings.Replace(runType.Url, "look ", "", -1)
	userid, err := hgGetUserID(runType.Url)
	if err != nil {
		return err
	}
	resData, err := hgGetUserWorkList(userid, 0, "")
	if err != nil {
		return err
	}
	count := resData.Data.Count
	page := int(math.Ceil(float64(count) / 20))
	PrintInfo(fmt.Sprintf("总页数：%d  每页个数：%d  总个数：%d", page, 20, count))
	startInt := 1
	endInt := page
	var downLoadList []map[string]string
	screenName := "--"
	startTime := time.Now().Unix()
	errorMsg := "--"
	errorCount := 0
	sleepTime := .5
	pageContext := ""
	onPage := 1
	for {
		if onPage > endInt {
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
		resData, err := hgGetUserWorkList(userid, (onPage-1)*10, pageContext)
		if err != nil {
			errorCount++
			errorMsg = err.Error()
			continue
		}
		pageContext = resData.Data.PageContext
		if onPage >= startInt {
			if len(resData.Data.Collections) > 0 {
				downLoadList = append(downLoadList, map[string]string{
					"page":  fmt.Sprintf("%d", onPage),
					"title": resData.Data.Collections[0].TvBoard.VideoData.Title,
				})
				screenName = resData.Data.Collections[0].TvBoard.User.UserInfo.UserName
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
	for _, v := range downLoadList {
		PrintInfo(fmt.Sprintf("第 %s 页 %s", v["page"], v["title"]))
	}
	return nil
}

// 通过 url 获取 UserID
func hgGetUserID(url string) (string, error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`^(http|https)://huoguo\.qq\.com/m/person\.html\?userid=(\d+)`),
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

// 请求作者作品列表
func hgGetUserWorkList(userid string, offset int, pageContext string) (*HuoguoUserVideoList, error) {
	var jsonData HuoguoUserVideoList
	var reqBody *bytes.Reader
	jsonMap := map[string]interface{}{
		"account": map[string]interface{}{
			"id":   userid,
			"type": 50,
		},
		"pageContext": pageContext,
		"countReg":    true,
	}
	bytesData, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	reqBody = bytes.NewReader(bytesData)
	req, err := http.NewRequest("POST", "https://hgaccess.video.qq.com/huoguo/user_work_list?vappid=49109510&vsecret=c1202d7f3ba41f86cdd2d3d1082605b4ed764c21e29520f3&callback=func&raw=1", reqBody)
	if err != nil {
		return &jsonData, err
	}
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("content-type", "application/x-www-form-urlencoded;charset:UTF-8")
	req.Header.Set("referer", "https://hgaccess.video.qq.com/")
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
	if jsonData.Msg != "ok." {
		return &jsonData, errors.New(jsonData.Msg)
	}
	return &jsonData, nil

}
