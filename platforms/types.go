package platforms

import (
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"time"
)

var Client = http.Client{Timeout: time.Second * 30}

var UserAgentPc = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36 Edg/84.0.522.61"
var UserAgentWap = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/84.0.4147.135"
var AppPath = fmt.Sprintf("%s\\goannie", os.Getenv("APPDATA"))
var AppBinPath = fmt.Sprintf("%s\\bin", AppPath)
var AnnieFile = fmt.Sprintf("%s\\annie.exe", AppBinPath)
var FfmpegFile = fmt.Sprintf("%s\\ffmpeg.exe", AppBinPath)

type DownloadPrint struct {
	Site      string
	Title     string
	Type      string
	Quality   string
	Size      string
	SizeBytes int64
}

// 初始化
func (d *DownloadPrint) Init(url string) {
	d.SetSize(url)
	d.FormatSize()
}

// 获取文件大小
func (d *DownloadPrint) SetSize(url string) {
	reqHead, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}
	resData, err := Client.Do(reqHead)
	if err != nil {
		return
	}
	defer resData.Body.Close()
	ranges := resData.Header.Get("Accept-Ranges")
	if ranges != "bytes" {
		d.SizeBytes = resData.ContentLength
		return
	}
	d.SizeBytes = resData.ContentLength
}

// 格式化字节
func (d *DownloadPrint) FormatSize() {
	fileSize := d.SizeBytes
	if fileSize < 1024 {
		d.Size = fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		d.Size = fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		d.Size = fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// 打印
func (d DownloadPrint) Print() {
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Site:      ")
	color.Unset()
	fmt.Println(d.Site)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Title:     ")
	color.Unset()
	fmt.Println(d.Title)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Type:      ")
	color.Unset()
	fmt.Println(d.Type)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf(" Streams:   ")
	color.Unset()
	fmt.Println("# All available quality")
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("     [default]  -------------------\n")
	color.Unset()
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("     Quality:         ")
	color.Unset()
	fmt.Println(d.Quality)
	color.Set(color.FgBlue, color.Bold)
	fmt.Printf("     Size:            ")
	color.Unset()
	fmt.Printf("%s (%d Bytes)\n", d.Size, d.SizeBytes)
}

type RunType struct {
	Url        string
	SavePath   string
	CookieFile string
}

// 腾讯归档API
type TengxunPlaysource struct {
	PlaylistItem struct {
		AsyncParam    string        `json:"asyncParam"`
		BtnList       []interface{} `json:"btnList"`
		BtnPlayURL    string        `json:"btnPlayUrl"`
		BtnTitle      string        `json:"btnTitle"`
		DisplayType   int           `json:"displayType"`
		FirstResList  []string      `json:"firstResList"`
		IndexList     []string      `json:"indexList"`
		IndexList2    []string      `json:"indexList2"`
		Name          string        `json:"name"`
		NeedAsync     bool          `json:"needAsync"`
		PayType       int           `json:"payType"`
		PlVideoType   int           `json:"pl_video_type"`
		RealName      string        `json:"realName"`
		StrIconURL    string        `json:"strIconUrl"`
		Title         string        `json:"title"`
		TotalEpisode  int           `json:"totalEpisode"`
		VideoPlayList []struct {
			EpisodeNumber string        `json:"episode_number"`
			ID            string        `json:"id"`
			MarkLabelList []interface{} `json:"markLabelList"`
			PayType       int           `json:"payType"`
			Pic           string        `json:"pic"`
			PlayURL       string        `json:"playUrl"`
			ThirdLine     string        `json:"thirdLine"`
			Title         string        `json:"title"`
			Type          string        `json:"type"`
		} `json:"videoPlayList"`
	} `json:"PlaylistItem"`
	Error int    `json:"error"`
	Msg   string `json:"msg"`
}

// 爱奇艺归档API
type IqiyiSvlistinfo struct {
	Code string `json:"code"`
	Data map[string][]struct {
		TvID            int64  `json:"tvId"`
		Description     string `json:"description"`
		Subtitle        string `json:"subtitle"`
		Vid             string `json:"vid"`
		Name            string `json:"name"`
		PlayURL         string `json:"playUrl"`
		IssueTime       int64  `json:"issueTime"`
		ContentType     int    `json:"contentType"`
		PayMark         int    `json:"payMark"`
		PayMarkURL      string `json:"payMarkUrl"`
		ImageURL        string `json:"imageUrl"`
		Duration        string `json:"duration"`
		AlbumImageURL   string `json:"albumImageUrl"`
		Period          string `json:"period"`
		Exclusive       bool   `json:"exclusive"`
		Order           int    `json:"order"`
		QiyiProduced    bool   `json:"qiyiProduced"`
		Focus           string `json:"focus"`
		ShortTitle      string `json:"shortTitle"`
		DownloadAllowed bool   `json:"downloadAllowed"`
		Is1080P         int    `json:"is1080p"`
		IP              struct {
			ID         string        `json:"id"`
			Deleted    string        `json:"deleted"`
			Books      []interface{} `json:"books"`
			Games      []interface{} `json:"games"`
			Tickets    []interface{} `json:"tickets"`
			Comicbooks []interface{} `json:"comicbooks"`
		} `json:"ip"`
	} `json:"data"`
}

// 爱奇艺归档html数据
type IqiyiPlayPageInfo struct {
	AlbumId string
	Cid     string
}

// 西瓜视频信息API
type XiguaInfo struct {
	Ck struct {
	} `json:"_ck"`
	Data struct {
		IsOriginal      bool          `json:"is_original"`
		Title                string   `json:"title"`
		URL                  string   `json:"url"`
		VideoID              string   `json:"video_id"`
	} `json:"data"`
	Success bool `json:"success"`
}

// 西瓜TA的视频列表API
type XiguaUserList struct {
	UserInfo    struct {
		Name              string `json:"name"`
	} `json:"user_info"`
	Message          string `json:"message"`
	Data             []struct {
		MediaName     string `json:"media_name"`
		Title         string `json:"title"`
		ArticleURL    string `json:"article_url"`
		BehotTime     int    `json:"behot_time"`
		UserInfo      struct {
			Name            string `json:"name"`
		} `json:"user_info"`
	} `json:"data"`
}
