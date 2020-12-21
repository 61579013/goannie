package ixigua

type SsrHydratedData struct {
	AnyVideo struct {
		GidInformation struct {
			Gid        string `json:"gid"`
			PackerData struct {
				Video struct {
					Title         string `json:"title"`
					VID           string `json:"vid"`
					VideoResource struct {
						Normal struct {
							Status         int     `json:"status"`
							Message        string  `json:"message"`
							EnableSsl      bool    `json:"enable_ssl"`
							EnableAdaptive bool    `json:"enable_adaptive"`
							VideoID        string  `json:"video_id"`
							VideoDuration  float64 `json:"video_duration"`
							MediaType      string  `json:"media_type"`
							BigThumbs      []struct {
								ImgNum   int      `json:"img_num"`
								URI      string   `json:"uri"`
								ImgURL   string   `json:"img_url"`
								ImgUrls  []string `json:"img_urls"`
								ImgXSize int      `json:"img_x_size"`
								ImgYSize int      `json:"img_y_size"`
								ImgXLen  int      `json:"img_x_len"`
								ImgYLen  int      `json:"img_y_len"`
								Duration float64  `json:"duration"`
								Interval int      `json:"interval"`
								Fext     string   `json:"fext"`
							} `json:"big_thumbs"`
							VideoList struct {
								Video1 struct {
									Definition      string `json:"definition"`
									Quality         string `json:"quality"`
									Vtype           string `json:"vtype"`
									Vwidth          int    `json:"vwidth"`
									Vheight         int    `json:"vheight"`
									Bitrate         int    `json:"bitrate"`
									Fps             int    `json:"fps"`
									CodecType       string `json:"codec_type"`
									Size            int    `json:"size"`
									MainURL         string `json:"main_url"`
									BackupURL1      string `json:"backup_url_1"`
									URLExpire       int    `json:"url_expire"`
									PreloadSize     int    `json:"preload_size"`
									PreloadInterval int    `json:"preload_interval"`
									PreloadMinStep  int    `json:"preload_min_step"`
									PreloadMaxStep  int    `json:"preload_max_step"`
									FileHash        string `json:"file_hash"`
									FileID          string `json:"file_id"`
									P2PVerifyURL    string `json:"p2p_verify_url"`
								} `json:"video_1"`
								Video2 struct {
									Definition      string `json:"definition"`
									Quality         string `json:"quality"`
									Vtype           string `json:"vtype"`
									Vwidth          int    `json:"vwidth"`
									Vheight         int    `json:"vheight"`
									Bitrate         int    `json:"bitrate"`
									Fps             int    `json:"fps"`
									CodecType       string `json:"codec_type"`
									Size            int    `json:"size"`
									MainURL         string `json:"main_url"`
									BackupURL1      string `json:"backup_url_1"`
									URLExpire       int    `json:"url_expire"`
									PreloadSize     int    `json:"preload_size"`
									PreloadInterval int    `json:"preload_interval"`
									PreloadMinStep  int    `json:"preload_min_step"`
									PreloadMaxStep  int    `json:"preload_max_step"`
									FileHash        string `json:"file_hash"`
									FileID          string `json:"file_id"`
									P2PVerifyURL    string `json:"p2p_verify_url"`
								} `json:"video_2"`
								Video3 struct {
									Definition      string `json:"definition"`
									Quality         string `json:"quality"`
									Vtype           string `json:"vtype"`
									Vwidth          int    `json:"vwidth"`
									Vheight         int    `json:"vheight"`
									Bitrate         int    `json:"bitrate"`
									Fps             int    `json:"fps"`
									CodecType       string `json:"codec_type"`
									Size            int    `json:"size"`
									MainURL         string `json:"main_url"`
									BackupURL1      string `json:"backup_url_1"`
									URLExpire       int    `json:"url_expire"`
									PreloadSize     int    `json:"preload_size"`
									PreloadInterval int    `json:"preload_interval"`
									PreloadMinStep  int    `json:"preload_min_step"`
									PreloadMaxStep  int    `json:"preload_max_step"`
									FileHash        string `json:"file_hash"`
									FileID          string `json:"file_id"`
									P2PVerifyURL    string `json:"p2p_verify_url"`
								} `json:"video_3"`
								Video4 struct {
									Definition      string `json:"definition"`
									Quality         string `json:"quality"`
									Vtype           string `json:"vtype"`
									Vwidth          int    `json:"vwidth"`
									Vheight         int    `json:"vheight"`
									Bitrate         int    `json:"bitrate"`
									Fps             int    `json:"fps"`
									CodecType       string `json:"codec_type"`
									Size            int    `json:"size"`
									MainURL         string `json:"main_url"`
									BackupURL1      string `json:"backup_url_1"`
									URLExpire       int    `json:"url_expire"`
									PreloadSize     int    `json:"preload_size"`
									PreloadInterval int    `json:"preload_interval"`
									PreloadMinStep  int    `json:"preload_min_step"`
									PreloadMaxStep  int    `json:"preload_max_step"`
									FileHash        string `json:"file_hash"`
									FileID          string `json:"file_id"`
									P2PVerifyURL    string `json:"p2p_verify_url"`
								} `json:"video_4"`
							} `json:"video_list"`
						} `json:"normal"`

						Dash120Fps struct {
							Status         int     `json:"status"`
							Message        string  `json:"message"`
							EnableSsl      bool    `json:"enable_ssl"`
							EnableAdaptive bool    `json:"enable_adaptive"`
							VideoID        string  `json:"video_id"`
							VideoDuration  float64 `json:"video_duration"`
							MediaType      string  `json:"media_type"`
							BigThumbs      []struct {
								ImgNum   int      `json:"img_num"`
								URI      string   `json:"uri"`
								ImgURL   string   `json:"img_url"`
								ImgUrls  []string `json:"img_urls"`
								ImgXSize int      `json:"img_x_size"`
								ImgYSize int      `json:"img_y_size"`
								ImgXLen  int      `json:"img_x_len"`
								ImgYLen  int      `json:"img_y_len"`
								Duration float64  `json:"duration"`
								Interval int      `json:"interval"`
								Fext     string   `json:"fext"`
							} `json:"big_thumbs"`
							DynamicVideo struct {
								DynamicType      string `json:"dynamic_type"`
								DynamicVideoList []struct {
									Definition string `json:"definition"`
									Quality    string `json:"quality"`
									Vtype      string `json:"vtype"`
									Vwidth     int    `json:"vwidth"`
									Vheight    int    `json:"vheight"`
									Bitrate    int    `json:"bitrate"`
									Size       int    `json:"size"`
									CodecType  string `json:"codec_type"`
									FileHash   string `json:"file_hash"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
									URLExpire  int    `json:"url_expire"`
									FileID     string `json:"file_id"`
									InitRange  string `json:"init_range"`
									IndexRange string `json:"index_range"`
									CheckInfo  string `json:"check_info"`
								} `json:"dynamic_video_list"`
								DynamicAudioList []struct {
									Quality    string `json:"quality"`
									Vtype      string `json:"vtype"`
									Bitrate    int    `json:"bitrate"`
									CodecType  string `json:"codec_type"`
									FileHash   string `json:"file_hash"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
									URLExpire  int    `json:"url_expire"`
									InitRange  string `json:"init_range"`
									IndexRange string `json:"index_range"`
									CheckInfo  string `json:"check_info"`
								} `json:"dynamic_audio_list"`
								MainURL    string `json:"main_url"`
								BackupURL1 string `json:"backup_url_1"`
							} `json:"dynamic_video"`
							PopularityLevel int `json:"popularity_level"`
							ExtraInfos      struct {
								Status            string `json:"Status"`
								Message           string `json:"Message"`
								LogoType          string `json:"LogoType"`
								VideoModelVersion int    `json:"VideoModelVersion"`
							} `json:"extraInfos"`
							InterfaceInfo struct {
								Code       int    `json:"code"`
								Message    string `json:"message"`
								Logid      string `json:"logid"`
								APIStr     string `json:"api_str"`
								Timestamep int64  `json:"timestamep"`
							} `json:"interfaceInfo"`
						} `json:"dash_120fps"`
					} `json:"videoResource"`
				} `json:"video"`
			} `json:"packerData"`
		} `json:"gidInformation"`
	} `json:"anyVideo"`
}
