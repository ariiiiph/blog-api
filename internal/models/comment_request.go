package models

type CreateCommentRequest struct {
	Content string `json:"content" example:"Great post, thanks for sharing!"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" example:"Updated comment text."`
}
