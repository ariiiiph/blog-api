package models

type CommentResponse struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	AuthorID  int64  `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentListResponse struct {
	Comments      []CommentResponse `json:"comments"`
	Page          int               `json:"page"`
	Limit         int               `json:"limit"`
	Count         int               `json:"count"`
	TotalComments int               `json:"total_comments"`
	TotalPages    int               `json:"total_pages"`
	HasNext       bool              `json:"has_next"`
	HasPrevious   bool              `json:"has_previous"`
}

type CreateCommentAPIResponse struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    CommentResponse `json:"data"`
}

type ListCommentsAPIResponse struct {
	Status  int                 `json:"status"`
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    CommentListResponse `json:"data"`
}
