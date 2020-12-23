package ixigua

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	var (
		title    string
		vid      string
		videoURL string
		audioURL string
		quality  string
	)
	if regexp.MustCompile(`"albumId"`).Match(html) {
		ratedData, err := getSsrHydratedDataEpisode(html)
		if err != nil {
			return err
		}
		episodeInfo := ratedData.AnyVideo.GidInformation.PackerData.EpisodeInfo
		title = fmt.Sprintf("%s %s", episodeInfo.Title, episodeInfo.Name)
		vid = episodeInfo.EpisodeID
		videoURL, quality = getDownloadURLEpisode(ratedData)
	} else {
		ratedData, err := getSsrHydratedData(html)
		if err != nil {
			return err
		}
		title = ratedData.AnyVideo.GidInformation.PackerData.Video.Title
		vid = ratedData.AnyVideo.GidInformation.PackerData.Video.VID
		videoURL, audioURL, quality = getDownloadURL(ratedData)
	}
	if videoURL == "" {
		return errors.New("获取下载链接失败")
	}
	if config.GetBool("app.debug") {
		fmt.Printf("title: %s\n", title)
		fmt.Printf("vid: %s\n", vid)
	}
	if config.GetBool("app.isFiltrationID") && vid != "" && option.Verify.Check(e.Key(), vid) {
		return fmt.Errorf("%s 在过滤库中，修改配置文件isFiltrationID为false不再过滤。", vid)
	}
	if err := download(e.Name(), videoURL, audioURL, title, quality, option); err != nil {
		return err
	}
	option.Verify.Add(e.Key(), vid)
	return nil
}

// New 创建一个解析器
func New() types.Extractor {
	return &extractor{}
}

func download(name, videoURL, audioURL, title, quality string, option types.Options) error {
	if config.GetBool("app.debug") {
		fmt.Printf("videoURL: %s\n", videoURL)
		fmt.Printf("audioURL: %s\n", audioURL)
	}
	types.NewDownloadPrint(fmt.Sprintf("%s ixigua.com", name), title, quality, videoURL).Print()
	if audioURL != "" {
		videoFile := fmt.Sprintf("%s_video.mp4", title)
		if err := downloader.New(videoURL, option.SavePath).SetOutputName(videoFile).Run(); err != nil {
			return err
		}
		defer os.Remove(filepath.Join(option.SavePath, videoFile))
		audioFile := fmt.Sprintf("%s_audio.mp3", title)
		if err := downloader.New(audioURL, option.SavePath).SetOutputName(audioFile).Run(); err != nil {
			return err
		}
		defer os.Remove(filepath.Join(option.SavePath, audioFile))
		if err := exec.Command("ffmpeg", "-i", filepath.Join(option.SavePath, audioFile), "-i", filepath.Join(option.SavePath, videoFile), "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", filepath.Join(option.SavePath, fmt.Sprintf("%s.mp4", title))).Run(); err != nil {
			return err
		}
	} else {
		if err := downloader.New(videoURL, option.SavePath).SetOutputName(fmt.Sprintf("%s.mp4", title)).Run(); err != nil {
			return err
		}
	}
	fmt.Println("")
	return nil
}

func getDownloadURL(d *ssrHydratedData) (string, string, string) {
	videoList := d.AnyVideo.GidInformation.PackerData.Video.VideoResource.Dash120Fps.DynamicVideo.DynamicVideoList
	audioList := d.AnyVideo.GidInformation.PackerData.Video.VideoResource.Dash120Fps.DynamicVideo.DynamicAudioList
	if len(videoList) != 0 && videoList[len(videoList)-1].MainURL != "" && len(audioList) != 0 && videoList[len(audioList)-1].MainURL != "" {
		videoURL, _ := base64.StdEncoding.DecodeString(videoList[len(videoList)-1].MainURL)
		audioURL, _ := base64.StdEncoding.DecodeString(audioList[len(audioList)-1].MainURL)
		return string(videoURL), string(audioURL), videoList[len(videoList)-1].Definition
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
	return string(r), "", quality
}

func getDownloadURLEpisode(d *ssrHydratedDataEpisode) (string, string) {
	vl := d.AnyVideo.GidInformation.PackerData.VideoResource.Normal.VideoList
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

func getSsrHydratedData(html []byte) (*ssrHydratedData, error) {
	jsonStrFind := regexp.MustCompile(`window\._SSR_HYDRATED_DATA=(.*?)</script>`).FindSubmatch(html)
	if len(jsonStrFind) < 2 {
		return nil, errors.New("解析数据失败")
	}
	jsonStr := strings.ReplaceAll(string(jsonStrFind[1]), ":undefined", ":\"undefined\"")
	var jsonData ssrHydratedData
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}

func getSsrHydratedDataEpisode(html []byte) (*ssrHydratedDataEpisode, error) {
	jsonStrFind := regexp.MustCompile(`window\._SSR_HYDRATED_DATA=(.*?)</script>`).FindSubmatch(html)
	if len(jsonStrFind) < 2 {
		return nil, errors.New("解析数据失败")
	}
	jsonStr := strings.ReplaceAll(string(jsonStrFind[1]), ":undefined", ":\"undefined\"")
	var jsonData ssrHydratedDataEpisode
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, err
	}
	return &jsonData, nil
}
