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
var Aria2File = fmt.Sprintf("%s\\aria2c.exe", AppBinPath)

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

// 爱奇艺作者作品列表API
type TengxunUserVideoList struct {
	RequestID int    `json:"requestId"`
	Ret       int    `json:"ret"`
	ErrorMsg  string `json:"errorMsg"`
	FuncRet   int    `json:"funcRet"`
	TransInfo struct {
	} `json:"transInfo"`
	Body struct {
		Modules []struct {
			Sections []struct {
				ReportDict struct {
				} `json:"report_dict"`
				SpecialBlocks struct {
				} `json:"special_blocks"`
				OperationMap struct {
				} `json:"operation_map"`
				SectionType       string      `json:"section_type"`
				SectionLayoutType int         `json:"section_layout_type"`
				CSSStruct         interface{} `json:"css_struct"`
				SectionID         string      `json:"section_id"`
				BlockList         struct {
					Blocks []struct {
						OperationMap struct {
							Num0 struct {
								ReportDict struct {
								} `json:"report_dict"`
								OperationType string `json:"operation_type"`
								Operation     struct {
									URL string `json:"url"`
								} `json:"operation"`
								ReportID string `json:"report_id"`
							} `json:"0"`
							Num300 struct {
								ReportDict struct {
								} `json:"report_dict"`
								OperationType string `json:"operation_type"`
								Operation     struct {
									PraiseData struct {
										PraiseType     string `json:"praise_type"`
										PraiseDataKey  string `json:"praise_data_key"`
										PraiseMatchKey string `json:"praise_match_key"`
									} `json:"praise_data"`
									PraiseStatus string `json:"praise_status"`
									PraiseUIInfo struct {
										PraiseCount string `json:"praise_count"`
									} `json:"praise_ui_info"`
									DislikeStatus string `json:"dislike_status"`
								} `json:"operation"`
								ReportID string `json:"report_id"`
							} `json:"300"`
						} `json:"operation_map"`
						MarkLabelListMap struct {
							Num0 struct {
								MarkLabelList []struct {
									MarkLabelType string `json:"mark_label_type"`
									Position      int    `json:"position"`
									PrimeText     string `json:"prime_text"`
									MarkImageURL  string `json:"mark_image_url"`
								} `json:"mark_label_list"`
							} `json:"0"`
						} `json:"mark_label_list_map"`
						ReportDict struct {
							ModID            string `json:"mod_id"`
							Rtype            string `json:"rtype"`
							ItemIdx          string `json:"item_idx"`
							PosterType       string `json:"poster_type"`
							ModIdx           string `json:"mod_idx"`
							Vid              string `json:"vid"`
							TabPersonalValue string `json:"tab_personal_value"`
						} `json:"report_dict"`
						BlockType      string      `json:"block_type"`
						BlockStyleType int         `json:"block_style_type"`
						CSSStruct      interface{} `json:"css_struct"`
						Data           struct {
							TagText  []interface{} `json:"tag_text"`
							CardType string        `json:"card_type"`
							CardInfo struct {
								ImageList []struct {
									ImageURL       string      `json:"image_url"`
									ThumbURL       string      `json:"thumb_url"`
									ImageType      string      `json:"image_type"`
									AspectRatio    float64     `json:"aspect_ratio"`
									ImageFacePoint interface{} `json:"image_face_point"`
									ExtraData      interface{} `json:"extra_data"`
								} `json:"image_list"`
								UserInfo struct {
									AccountInfo struct {
										AccountType int    `json:"account_type"`
										AccountID   string `json:"account_id"`
									} `json:"account_info"`
									UserType     string      `json:"user_type"`
									UserName     string      `json:"user_name"`
									UserImageURL string      `json:"user_image_url"`
									UserLabelURL string      `json:"user_label_url"`
									ExtraData    interface{} `json:"extra_data"`
								} `json:"user_info"`
							} `json:"card_info"`
							Title   string `json:"title"`
							Content string `json:"content"`
							Vid     string `json:"vid"`
						} `json:"data"`
						BlockID   string      `json:"block_id"`
						VnViewID  string      `json:"vn_view_id"`
						ExtraData interface{} `json:"extra_data"`
					} `json:"blocks"`
					OptionalBlocks []interface{} `json:"optional_blocks"`
				} `json:"block_list"`
				ExtraAnyData interface{} `json:"extra_any_data"`
				MergeID      string      `json:"merge_id"`
			} `json:"sections"`
			ReportDict struct {
			} `json:"report_dict"`
			ID           string      `json:"id"`
			ExtraAnyData interface{} `json:"extra_any_data"`
			MergeID      string      `json:"merge_id"`
			UniqueID     string      `json:"unique_id"`
		} `json:"modules"`
		PageContext struct {
			LastVidPosition string `json:"last_vid_position"`
			Offset          string `json:"offset"`
			IndexContext    string `json:"index_context"`
		} `json:"page_context"`
		ReportDict struct {
		} `json:"report_dict"`
		RequestContext struct {
		} `json:"request_context"`
		PrePageContext struct {
		} `json:"pre_page_context"`
		HasNextPage                bool        `json:"has_next_page"`
		StyleCollectionCheckResult interface{} `json:"style_collection_check_result"`
		ReportPageID               string      `json:"report_page_id"`
		ExtraData                  interface{} `json:"extra_data"`
		HasPrePage                 bool        `json:"has_pre_page"`
	} `json:"body"`
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
		IsOriginal bool   `json:"is_original"`
		Title      string `json:"title"`
		URL        string `json:"url"`
		VideoID    string `json:"video_id"`
	} `json:"data"`
	Success bool `json:"success"`
}

// 西瓜TA的视频列表API
type XiguaUserList struct {
	UserInfo struct {
		Name string `json:"name"`
	} `json:"user_info"`
	Message          string `json:"message"`
	HasMore bool   `json:"has_more"`
	Data             []struct {
		MediaName  string `json:"media_name"`
		Title      string `json:"title"`
		ArticleURL string `json:"article_url"`
		BehotTime  int    `json:"behot_time"`
		UserInfo   struct {
			Name string `json:"name"`
		} `json:"user_info"`
	} `json:"data"`
}

// 火锅视频作者视频列表API
type HuoguoUserVideoList struct {
	Data struct {
		ErrCode     int    `json:"errCode"`
		PageContext string `json:"pageContext"`
		Collections []struct {
			TvBoard struct {
				Poster struct {
					FirstLine string `json:"firstLine"`
					ImageURL  string `json:"imageUrl"`
					Action    struct {
						URL          string `json:"url"`
						CacheType    int    `json:"cacheType"`
						PreReadType  int    `json:"preReadType"`
						ReportParams string `json:"reportParams"`
					} `json:"action"`
					PlayCountL   int     `json:"playCountL"`
					DisplayRatio float64 `json:"displayRatio"`
					FaceArea     struct {
						XPoint float64 `json:"xPoint"`
						YPoint float64 `json:"yPoint"`
					} `json:"faceArea"`
					ImageURLRatio float64 `json:"imageUrlRatio"`
					GifURLRatio   float64 `json:"gifUrlRatio"`
				} `json:"poster"`
				VideoData struct {
					Vid       string `json:"vid"`
					PayStatus int    `json:"payStatus"`
					Poster    struct {
						Action struct {
							CacheType   int `json:"cacheType"`
							PreReadType int `json:"preReadType"`
						} `json:"action"`
						PlayCountL   int     `json:"playCountL"`
						DisplayRatio float64 `json:"displayRatio"`
						FaceArea     struct {
							XPoint float64 `json:"xPoint"`
							YPoint float64 `json:"yPoint"`
						} `json:"faceArea"`
						ImageURLRatio float64 `json:"imageUrlRatio"`
						GifURLRatio   float64 `json:"gifUrlRatio"`
					} `json:"poster"`
					SkipStart               int    `json:"skipStart"`
					Title                   string `json:"title"`
					IsNoStroeWatchedHistory bool   `json:"isNoStroeWatchedHistory"`
					WatchRecordPoster       struct {
						Action struct {
							CacheType   int `json:"cacheType"`
							PreReadType int `json:"preReadType"`
						} `json:"action"`
						PlayCountL   int     `json:"playCountL"`
						DisplayRatio float64 `json:"displayRatio"`
						FaceArea     struct {
							XPoint float64 `json:"xPoint"`
							YPoint float64 `json:"yPoint"`
						} `json:"faceArea"`
						ImageURLRatio float64 `json:"imageUrlRatio"`
						GifURLRatio   float64 `json:"gifUrlRatio"`
					} `json:"watchRecordPoster"`
					ShareItem struct {
						ShareStyle int `json:"shareStyle"`
						ShareCount int `json:"shareCount"`
					} `json:"shareItem"`
					StreamRatio float64 `json:"streamRatio"`
				} `json:"videoData"`
				IsAutoPlayer   bool `json:"isAutoPlayer"`
				IsLoopPlayBack bool `json:"isLoopPlayBack"`
				AttentInfo     struct {
					AttentKey   string `json:"attentKey"`
					AttentState int    `json:"attentState"`
					Count       int    `json:"count"`
				} `json:"attentInfo"`
				ShareItem struct {
					ShareURL      string `json:"shareUrl"`
					ShareTitle    string `json:"shareTitle"`
					ShareSubtitle string `json:"shareSubtitle"`
					ShareImgURL   string `json:"shareImgUrl"`
					ShareStyle    int    `json:"shareStyle"`
					ShareCount    int    `json:"shareCount"`
				} `json:"shareItem"`
				TimeStamp int `json:"timeStamp"`
				User      struct {
					UserInfo struct {
						Account struct {
							Type int    `json:"type"`
							ID   string `json:"id"`
						} `json:"account"`
						UserName     string `json:"userName"`
						FaceImageURL string `json:"faceImageUrl"`
						DetailInfo   []struct {
							ItemKey   string `json:"itemKey"`
							ItemValue string `json:"itemValue"`
							ItemID    string `json:"itemId"`
						} `json:"detailInfo"`
						LegalizeInfo struct {
							Type int `json:"type"`
						} `json:"legalizeInfo"`
					} `json:"userInfo"`
					RelationItem struct {
						RelationKey string `json:"relationKey"`
						FromMe      int    `json:"fromMe"`
						ToMe        int    `json:"toMe"`
					} `json:"relationItem"`
					Action struct {
						URL         string `json:"url"`
						CacheType   int    `json:"cacheType"`
						PreReadType int    `json:"preReadType"`
					} `json:"action"`
					UserShareItem struct {
						ShareStyle int `json:"shareStyle"`
						ShareCount int `json:"shareCount"`
					} `json:"userShareItem"`
					PickInfo struct {
						Count      int `json:"count"`
						Trend      int `json:"trend"`
						PickScence struct {
							Scence int `json:"scence"`
						} `json:"pickScence"`
						Rank          int `json:"rank"`
						AllowPick     int `json:"allowPick"`
						ActionBarInfo struct {
							Action struct {
								CacheType   int `json:"cacheType"`
								PreReadType int `json:"preReadType"`
							} `json:"action"`
						} `json:"actionBarInfo"`
						TrendCount int `json:"trendCount"`
					} `json:"pickInfo"`
					OwnPickInfo struct {
						PickScence struct {
							Scence int `json:"scence"`
						} `json:"pickScence"`
						LeftPicks int `json:"leftPicks"`
					} `json:"ownPickInfo"`
				} `json:"user"`
				CommentInfo struct {
					Action struct {
						URL         string `json:"url"`
						CacheType   int    `json:"cacheType"`
						PreReadType int    `json:"preReadType"`
					} `json:"action"`
					CommentCount  int    `json:"commentCount"`
					HotCommentKey string `json:"hotCommentKey"`
				} `json:"commentInfo"`
				AuditStatus   int `json:"auditStatus"`
				PrivacyStatus int `json:"privacyStatus"`
				Duration      int `json:"duration"`
			} `json:"tvBoard"`
		} `json:"collections"`
		HasNextPage bool `json:"hasNextPage"`
		Count       int  `json:"count"`
	} `json:"data"`
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}