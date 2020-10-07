package platforms

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/garyburd/redigo/redis"
)

// RepetitionMsg 视频重复下载时的提示信息
const RepetitionMsg = "检测到当前视频重复下载，因之前设置了过滤重复，直接跳过当前下载。"

// Client 默认Cliend
var Client = http.Client{Timeout: time.Second * 30}

// UserAgentPc 电脑端user-agent
var UserAgentPc = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36 Edg/84.0.522.61"

// UserAgentWap 手机端端user-agent
var UserAgentWap = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/84.0.4147.135"

// AppPath 程序app目录
var AppPath = fmt.Sprintf("%s\\goannie", os.Getenv("APPDATA"))

// AppBinPath 程序bin目录
var AppBinPath = fmt.Sprintf("%s\\bin", AppPath)

// AppDataPath 程序data目录
var AppDataPath = fmt.Sprintf("%s\\data", AppPath)

// AnnieFile 程序annie存放位置
var AnnieFile = fmt.Sprintf("%s\\annie.exe", AppBinPath)

// FfmpegFile 程序ffmpeg存放位置
var FfmpegFile = fmt.Sprintf("%s\\ffmpeg.exe", AppBinPath)

// Aria2File 程序aria2存放位置
var Aria2File = fmt.Sprintf("%s\\aria2c.exe", AppBinPath)

// RedisFile 程序redis存放位置
var RedisFile = fmt.Sprintf("%s\\redis-server.exe", AppBinPath)

// RedisConfFile 程序redisconf存放位置
var RedisConfFile = fmt.Sprintf("%s\\redis.windows-service.conf", AppBinPath)

// DownloadPrint 程序下载时打印结构体
type DownloadPrint struct {
	Site      string
	Title     string
	Type      string
	Quality   string
	Size      string
	SizeBytes int64
}

// Init 初始化
func (d *DownloadPrint) Init(url string) {
	d.SetSize(url)
	d.FormatSize()
}

// SetSize 获取文件大小
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

// FormatSize 格式化字节
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

// Print 打印
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

// RunType 程序传递信息的主要结构体
type RunType struct {
	URL           string
	SavePath      string
	CookieFile    string
	DefaultCookie string
	IsDeWeight    bool
	RedisConn     redis.Conn
}

// TengxunPlaysource 腾讯归档API
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

// TengxunUserVideoList 腾讯作者作品列表API
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

// IqiyiSvlistinfo 爱奇艺归档API
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

// IqiyiPlayPageInfo 爱奇艺归档html数据
type IqiyiPlayPageInfo struct {
	AlbumId string
	Cid     string
}

// IqiyGetVideosAction 爱奇艺作品列表api
type IqiyGetVideosAction struct {
	Code string `json:"code"`
	Data struct {
		HasMore  int `json:"hasMore"`
		TotalNum int `json:"totalNum"`
		Sort     struct {
			ReprentativeWork interface{} `json:"reprentativeWork"`
			Flows            []struct {
				QipuID int64 `json:"qipuId"`
			} `json:"flows"`
		} `json:"sort"`
	} `json:"data"`
	Msg interface{} `json:"msg"`
}

// IqiyEpisodeInfoAction 爱奇艺作品详情api
type IqiyEpisodeInfoAction struct {
	Code string `json:"code"`
	Data map[string]struct {
		QipuID   int    `json:qipuId`
		Title    string `json:title`
		PageUrl  string `json:pageUrl`
		Nickname string `json:nickname`
		VID      string `json:vid`
	} `json:"data"`
	Msg interface{} `json:"msg"`
}

// XiguaInfo 西瓜视频信息API
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

// XiguaUserList 西瓜TA的视频列表API
type XiguaUserList struct {
	UserInfo struct {
		Name string `json:"name"`
	} `json:"user_info"`
	Message string `json:"message"`
	HasMore bool   `json:"has_more"`
	Data    []struct {
		MediaName  string `json:"media_name"`
		Title      string `json:"title"`
		ArticleURL string `json:"article_url"`
		GroupIdStr string `json:"group_id_str"`
		BehotTime  int    `json:"behot_time"`
		UserInfo   struct {
			Name string `json:"name"`
		} `json:"user_info"`
	} `json:"data"`
}

// HuoguoUserVideoList 火锅视频作者视频列表API
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

// HaokanUserVideoList 好看视频作者视频列表API
type HaokanUserVideoList struct {
	Errno int    `json:"errno"`
	Error string `json:"error"`
	Data  struct {
		RequestParam []interface{} `json:"requestParam"`
		Response     struct {
			ResponseCount int `json:"response_count"`
			HasMore       int `json:"has_more"`
			Results       []struct {
				TplName string `json:"tplName"`
				Type    string `json:"type"`
				Content struct {
					Vid              string        `json:"vid"`
					SortTime         int           `json:"sort_time"`
					PublishTime      string        `json:"publish_time"`
					Title            string        `json:"title"`
					CoverSrc         string        `json:"cover_src"`
					CoverSrcPc       string        `json:"cover_src_pc"`
					VideoSrc         string        `json:"video_src"`
					Duration         string        `json:"duration"`
					Authorid         string        `json:"authorid"`
					Poster           string        `json:"poster"`
					Thumbnails       string        `json:"thumbnails"`
					URL              string        `json:"url"`
					LocURL           string        `json:"loc_url"`
					FeedID           string        `json:"feed_id"`
					Author           string        `json:"author"`
					AuthorIcon       string        `json:"author_icon"`
					LikeNum          int           `json:"like_num"`
					RecType          int           `json:"rec_type"`
					ReadNum          int           `json:"read_num"`
					Playcnt          int           `json:"playcnt"`
					PlaycntText      string        `json:"playcntText"`
					HaokanSourceFrom string        `json:"haokan_source_from"`
					Ext              []interface{} `json:"ext"`
					VideoList        struct {
						Sd string `json:"sd"`
						Hd string `json:"hd"`
						Sc string `json:"sc"`
					} `json:"video_list"`
					Size struct {
						Sd int `json:"sd"`
						Hd int `json:"hd"`
						Sc int `json:"sc"`
					} `json:"size"`
					VideoType         string      `json:"videoType"`
					Ctk               string      `json:"ctk"`
					Pvid              string      `json:"pvid"`
					VideoShortURL     string      `json:"video_short_url"`
					AuthorPassportID  string      `json:"author_passport_id"`
					Dtime             int64       `json:"dtime"`
					VoteDisableCtrl   string      `json:"vote_disable_ctrl"`
					SensitiveFlags    interface{} `json:"sensitive_flags"`
					DisplaytypeExinfo struct {
						GoodsType []interface{} `json:"goods_type"`
					} `json:"displaytype_exinfo"`
					IsShowFeature bool   `json:"is_show_feature"`
					CommentNum    string `json:"commentNum"`
					PraiseNum     int    `json:"praiseNum"`
				} `json:"content"`
			} `json:"results"`
			Ctime string `json:"ctime"`
		} `json:"response"`
	} `json:"data"`
}

// BliUserVideoList 哔哩哔哩视频作者视频列表API
type BliUserVideoList struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		List struct {
			Tlist struct {
				Num3 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"3"`
				Num4 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"4"`
				Num160 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"160"`
			} `json:"tlist"`
			Vlist []struct {
				Comment      int    `json:"comment"`
				Typeid       int    `json:"typeid"`
				Play         int    `json:"play"`
				Pic          string `json:"pic"`
				Subtitle     string `json:"subtitle"`
				Description  string `json:"description"`
				Copyright    string `json:"copyright"`
				Title        string `json:"title"`
				Review       int    `json:"review"`
				Author       string `json:"author"`
				Mid          int    `json:"mid"`
				Created      int    `json:"created"`
				Length       string `json:"length"`
				VideoReview  int    `json:"video_review"`
				Aid          int    `json:"aid"`
				Bvid         string `json:"bvid"`
				HideClick    bool   `json:"hide_click"`
				IsPay        int    `json:"is_pay"`
				IsUnionVideo int    `json:"is_union_video"`
			} `json:"vlist"`
		} `json:"list"`
		Page struct {
			Count int `json:"count"`
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
		} `json:"page"`
	} `json:"data"`
}

// BliVideoP 多P api
type BliVideoP struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Page int    `json:"page"`
		Part string `json:"part"`
	} `json:"data"`
}

// DouyinVideoList 抖音视频作者列表API
type DouyinVideoList struct {
	AwemeList []struct {
		AwemeID   string `json:"aweme_id"`
		Desc      string `json:"desc"`
		ShareInfo struct {
			ShareTitle string `json:"share_title"`
		} `json:"share_info"`
		ShareURL string `json:"share_url"`
		Video    struct {
			DownloadAddr struct {
				URI      string   `json:"uri"`
				DataSize int      `json:"data_size"`
				Height   int      `json:"height"`
				Width    int      `json:"width"`
				URLList  []string `json:"url_list"`
			} `json:"download_addr"`
		} `json:"video"`
	} `json:"aweme_list"`
	HasMore    int `json:"has_more"`
	MaxCursor  int `json:"max_cursor"`
	MinCursor  int `json:"min_cursor"`
	StatusCode int `json:"status_code"`
}

// DouyinUserInfo 抖音作者信息
type DouyinUserInfo struct {
	StatusCode int `json:"status_code"`
	UserInfo   struct {
		Nickname        string `json:"nickname"`
		AwemeCount      int    `json:"aweme_count"`
		FavoritingCount int    `json:"favoriting_count"`
		FollowerCount   int    `json:"follower_count"`
		UID             string `json:"uid"`
	} `json:"user_info"`
}

// DouyinVideoInfo 抖音视频信息
type DouyinVideoInfo struct {
	Title    string
	Widht    int
	Height   int
	Duration int
	Ratio    string
	URI      string
	AwemeID  string
}

// DouyinVideoItem 视频信息接口
type DouyinVideoItem struct {
	StatusCode int `json:"status_code"`
	ItemList   []struct {
		Desc    string `json:"desc"`
		AwemeID string `json:"aweme_id"`
		Video   struct {
			Widht    int    `json:"width"`
			Height   int    `json:"height"`
			Ratio    string `json:"ratio"`
			Duration int    `json:"duration"`
			URI      string `json:"vid"`
		} `json:"video"`
	} `json:"item_list"`
}
