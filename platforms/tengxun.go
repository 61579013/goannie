package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 腾讯归档页：https://v.qq.com/detail/5/52852.html
func RunTxDetail(runType RunType) error {
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
			downLoadList = append(downLoadList, map[string]string{
				"title": item.Title,
				"url":   item.PlayURL,
			})
		}
	}
	PrintInfo(fmt.Sprintf("采集到 %d 个视频", len(downLoadList)))

	AnnieDownloadAll(downLoadList,runType)

	PrintInfo("全部下载完成")
	return nil
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
