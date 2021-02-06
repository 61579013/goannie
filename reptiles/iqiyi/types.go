package iqiyi

type videosActionAPI struct {
	Data struct {
		TotalNum int `json:"totalNum"`
		HasMore  int `json:"hasMore"`
		Sort     struct {
			Flows []struct {
				QipuID int64 `json:"qipuId"`
			} `json:"flows"`
		} `json:"sort"`
	} `json:"data"`
}

type episodeInfoAPI struct {
	Data map[string]struct {
		Title   string `json:"title"`
		PageURL string `json:"pageUrl"`
	} `json:"data"`
}
