package ixigua

type ssrHydratedData struct {
	AnyVideo struct {
		GidInformation struct {
			Gid        string `json:"gid"`
			PackerData struct {
				Video struct {
					Title         string `json:"title"`
					VID           string `json:"vid"`
					VideoResource struct {
						Normal struct {
							VideoID       string  `json:"video_id"`
							VideoDuration float64 `json:"video_duration"`
							MediaType     string  `json:"media_type"`
							VideoList     struct {
								Video1 struct {
									Definition string `json:"definition"`
									Quality    string `json:"quality"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
								} `json:"video_1"`
								Video2 struct {
									Definition string `json:"definition"`
									Quality    string `json:"quality"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
								} `json:"video_2"`
								Video3 struct {
									Definition string `json:"definition"`
									Quality    string `json:"quality"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
								} `json:"video_3"`
								Video4 struct {
									Definition string `json:"definition"`
									Quality    string `json:"quality"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
								} `json:"video_4"`
							} `json:"video_list"`
						} `json:"normal"`
						Dash120Fps struct {
							VideoID      string `json:"video_id"`
							MediaType    string `json:"media_type"`
							DynamicVideo struct {
								DynamicVideoList []struct {
									Definition string `json:"definition"`
									Quality    string `json:"quality"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
								} `json:"dynamic_video_list"`
								DynamicAudioList []struct {
									Quality    string `json:"quality"`
									MainURL    string `json:"main_url"`
									BackupURL1 string `json:"backup_url_1"`
								} `json:"dynamic_audio_list"`
								MainURL    string `json:"main_url"`
								BackupURL1 string `json:"backup_url_1"`
							} `json:"dynamic_video"`
						} `json:"dash_120fps"`
					} `json:"videoResource"`
				} `json:"video"`
			} `json:"packerData"`
		} `json:"gidInformation"`
	} `json:"anyVideo"`
}

type ssrHydratedDataEpisode struct {
	AnyVideo struct {
		GidInformation struct {
			AlbumID    string `json:"albumId"`
			PackerData struct {
				EpisodeInfo struct {
					EpisodeID   string `json:"episodeId"`
					AlbumID     string `json:"albumId"`
					Rank        int    `json:"rank"`
					SeqOld      string `json:"seqOld"`
					Title       string `json:"title"`
					Name        string `json:"name"`
					BottomLabel string `json:"bottomLabel"`
				} `json:"episodeInfo"`
				VideoResource struct {
					Vid    string `json:"vid"`
					Normal struct {
						VideoList struct {
							Video1 struct {
								Definition string `json:"definition"`
								Quality    string `json:"quality"`
								Vtype      string `json:"vtype"`
								MainURL    string `json:"main_url"`
								BackupURL1 string `json:"backup_url_1"`
							} `json:"video_1"`
							Video2 struct {
								Definition string `json:"definition"`
								Quality    string `json:"quality"`
								Vtype      string `json:"vtype"`
								MainURL    string `json:"main_url"`
								BackupURL1 string `json:"backup_url_1"`
							} `json:"video_2"`
							Video3 struct {
								Definition string `json:"definition"`
								Quality    string `json:"quality"`
								Vtype      string `json:"vtype"`
								MainURL    string `json:"main_url"`
								BackupURL1 string `json:"backup_url_1"`
							} `json:"video_3"`
							Video4 struct {
								Definition string `json:"definition"`
								Quality    string `json:"quality"`
								Vtype      string `json:"vtype"`
								MainURL    string `json:"main_url"`
								BackupURL1 string `json:"backup_url_1"`
							} `json:"video_4"`
						} `json:"video_list"`
					} `json:"normal"`
				} `json:"videoResource"`
			} `json:"packerData"`
		} `json:"gidInformation"`
	} `json:"anyVideo"`
}

// Stream is the data structure for each video stream, eg: 720P, 1080P.
type Stream struct {
	ID      string
	Quality string
	URL     string
}
