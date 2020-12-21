package ixigua

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gitee.com/rock_rabbit/goannie/config"

	"gitee.com/rock_rabbit/goannie/downloader"
	"gitee.com/rock_rabbit/goannie/extractors/types"
	"gitee.com/rock_rabbit/goannie/request"
)

const (
	referer   = "https://www.ixigua.com"
	host      = "https://www.ixigua.com"
	defCookie = "xiguavideopcwebid=6872983880459503118; xiguavideopcwebid.sig=B4DvNwwGGQ-hDxYcJo5FfbMEIn4; _ga=GA1.2.572711536.1600241266; MONITOR_WEB_ID=bfe0e43a-e004-400e-8040-81677a199a22; ttwid=1%7CPWHvUSGTtsxK0WUzkuq7SxJtT7L3WHRvJeSGG5WZjiw%7C1604995289%7Cec6a591ac986362929a9be173d65df8f3551269fff0694d34a5e57a33cd287eb; ixigua-a-s=1; Hm_lvt_db8ae92f7b33b6596893cdf8c004a1a2=1608261601; _gid=GA1.2.1203395873.1608261601; Hm_lpvt_db8ae92f7b33b6596893cdf8c004a1a2=1608262109"
)

type extractor struct{}

// Name 平台名称
func (e *extractor) Name() string {
	return "西瓜视频"
}

// Key 存储器使用的键
func (e *extractor) Key() string {
	return "xigua"
}

// DefCookie 默认cookie
func (e *extractor) DefCookie() string {
	return defCookie
}

// Extract 运行解析器
func (e *extractor) Extract(url string, option types.Options) error {
	html, err := request.GetByte(url, referer, nil)
	if err != nil {
		return err
	}
	ratedData, err := getSsrHydratedData(html)
	if err != nil {
		return err
	}
	title := ratedData.AnyVideo.GidInformation.PackerData.Video.Title
	vid := ratedData.AnyVideo.GidInformation.PackerData.Video.VID
	if config.GetBool("app.debug") {
		fmt.Printf("title: %s\n", title)
		fmt.Printf("vid: %s\n", vid)
	}
	if config.GetBool("app.isFiltrationID") && vid != "" && option.Verify.Check(e.Key(), vid) {
		return fmt.Errorf("%s 在过滤库中，修改配置文件isFiltrationID为false不再过滤。", vid)
	}
	downloadURL, quality := getNewDownloadURL(ratedData)
	if downloadURL == "" {
		return errors.New("获取下载链接失败")
	}
	if err := download(e.Name(), downloadURL, title, quality, option); err != nil {
		return err
	}
	option.Verify.Add(e.Key(), vid)
	return nil
}

// New 创建一个解析器
func New() types.Extractor {
	return &extractor{}
}

func download(name, url, title, quality string, option types.Options) error {
	if config.GetBool("app.debug") {
		fmt.Printf("downloadURL: %s\n", url)
	}
	types.NewDownloadPrint(fmt.Sprintf("%s ixigua.com", name), title, quality, url).Print()
	if err := downloader.New(url, option.SavePath).SetOutputName(fmt.Sprintf("%s.mp4", title)).Run(); err != nil {
		return err
	}
	fmt.Println("")
	return nil
}

func getNewDownloadURL(d *SsrHydratedData) (string, string) {
	videoList := d.AnyVideo.GidInformation.PackerData.Video.VideoResource.Dash120Fps.DynamicVideo.DynamicVideoList
	if len(videoList) != 0 && videoList[len(videoList)-1].MainURL != "" {
		decoded, _ := base64.StdEncoding.DecodeString(videoList[len(videoList)-1].MainURL)
		return string(decoded), videoList[len(videoList)-1].Definition
	}
	vl := d.AnyVideo.GidInformation.PackerData.Video.VideoResource.Normal.VideoList
	var murl string
	var quality string
	if vl.Video4.MainURL != "" {
		murl = vl.Video4.MainURL
		quality = vl.Video4.Definition
	} else if vl.Video3.MainURL != "" {
		murl = vl.Video3.MainURL
		quality = vl.Video3.Definition
	} else if vl.Video2.MainURL != "" {
		murl = vl.Video2.MainURL
		quality = vl.Video2.Definition
	} else if vl.Video1.MainURL != "" {
		murl = vl.Video1.MainURL
		quality = vl.Video1.Definition
	}
	r, _ := base64.StdEncoding.DecodeString(murl)
	return string(r), quality
}

func getSsrHydratedData(html []byte) (*SsrHydratedData, error) {
	jsonStrFind := regexp.MustCompile(`window\._SSR_HYDRATED_DATA=(.*?)</script>`).FindSubmatch(html)
	if len(jsonStrFind) < 2 {
		return nil, errors.New("解析数据失败")
	}
	jsonStr := strings.ReplaceAll(string(jsonStrFind[1]), ":undefined", ":\"undefined\"")
	var jsonData SsrHydratedData
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}
