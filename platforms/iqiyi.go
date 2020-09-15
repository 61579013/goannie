package platforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

// RunIqyOne 单视频
func RunIqyOne(runType RunType, arg map[string]string) error {
	err := AnnieDownload(runType.URL, runType.SavePath, runType.CookieFile, runType.DefaultCookie)
	if err != nil {
		return err
	}
	return nil
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &jsonData, err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("cookie", GetTxtContent(cookiePath))
	req.Header.Set("referer", "https://www.iqiyi.com/")
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
	err = json.Unmarshal([]byte(content), &jsonData)
	if err != nil {
		return &jsonData, err
	}
	if jsonData.Code != "A00000" {
		return &jsonData, errors.New("爱奇艺请求错误：" + jsonData.Code)
	}
	return &jsonData, nil
}

// DetailGetPlayPageInfo 爱奇艺归档请求HTML
func DetailGetPlayPageInfo(url, cookiePath string) (*IqiyiPlayPageInfo, error) {
	var jsonData IqiyiPlayPageInfo
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &jsonData, err
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("cookie", GetTxtContent(cookiePath))
	req.Header.Set("referer", "https://www.iqiyi.com/")
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
	resDatas := regexp.MustCompile(`albumId: "(\w+)",(\s|[\r\n])+tvId: "\w+",(\s|[\r\n])+sourceId: \w+,(\s|[\r\n])+cid: "(\w+)",`).FindStringSubmatch(content)
	if len(resDatas) < 6 {
		return &jsonData, errors.New("获取 albumId,cid 失败")
	}
	jsonData.AlbumId = resDatas[1]
	jsonData.Cid = resDatas[5]
	return &jsonData, nil
}
