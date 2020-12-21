package ixigua

type homeHydratedData struct {
	AuthorDetailInfo struct {
		UserID          string `json:"user_id"`
		MediaID         string `json:"media_id"`
		Name            string `json:"name"`
		Introduce       string `json:"introduce"`
		FansNum         string `json:"fansNum"`
		FollowNum       string `json:"followNum"`
		DiggNum         int    `json:"diggNum"`
		VerifiedContent string `json:"verified_content"`
	} `json:"AuthorDetailInfo"`
	AuthorTabsCount struct {
		VideoCnt int `json:"videoCnt"`
	} `json:"AuthorTabsCount"`
}

type homeAuthorVideo struct {
	Code int `json:"code"`
	Data struct {
		Message string `json:"message"`
		Data    []struct {
			Gid            string `json:"gid"`
			PublishTime    int    `json:"publish_time"`
			Title          string `json:"title"`
			VideoID        string `json:"video_id"`
			VideoLikeCount int    `json:"video_like_count"`
		} `json:"data"`
		HasMore          bool `json:"has_more"`
		HasMoreToRefresh bool `json:"has_more_to_refresh"`
		TotalNumber      int  `json:"total_number"`
	} `json:"data"`
}
