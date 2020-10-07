package platforms

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

// RunYkOne 单视频
func RunYkOne(runType RunType, arg map[string]string) error {
	return nil
}

// RunYkUserlist 作者
func RunYkUserlist(runType RunType, arg map[string]string) error {
	return nil
}

// YkGetUserVideolist 请求作者列表
func YkGetUserVideolist(page string) ([]map[string]string, error) {
	var err error
	var retData []map[string]string
	api := "https://i.youku.com/i/UNjMwMTY2MDUyMA==/videos?order=1&page=%s"
	res, err := RequestGet(fmt.Sprintf(api, page), map[string]string{
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"pragma":     "no-cache",
		"user-agent": UserAgentPc,
	})
	if err != nil {
		return retData, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return retData, err
	}
	mustList := regexp.MustCompile(`<div class="v va">(.*?)</span></div></div></div>`).FindAllStringSubmatch(string(body), -1)

	for _, i := range mustList {
		if len(i) == 0 {
			continue
		}
		centent := regexp.MustCompile(`<a href="(.*?)".*?title="(.*?)">.*?</a>`).FindStringSubmatch(i[0])
		if len(centent) < 3 {
			continue
		}
		retData = append(retData, map[string]string{
			"title": centent[2],
			"url":   fmt.Sprintf("https%s", centent[1]),
		})
	}

	return retData, nil
}
