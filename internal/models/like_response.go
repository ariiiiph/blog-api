package models

type LikeCountResponse struct {
	PostID int64 `json:"post_id"`
	Likes  int64 `json:"likes"`
}

type LikeCountAPIResponse struct {
	Status  int               `json:"status"`
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    LikeCountResponse `json:"data"`
}
