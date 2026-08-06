package models

type CreatePostRequest struct {
	CategoryID int64  `json:"category_id" example:"1"`
	Title      string `json:"title" example:"My First Post"`
	Slug       string `json:"slug" example:"my-first-post"`
	Content    string `json:"content" example:"This is the content of my post."`
}

type UpdatePostRequest struct {
	CategoryID int64  `json:"category_id" example:"1"`
	Title      string `json:"title" example:"My Updated Post"`
	Slug       string `json:"slug" example:"my-updated-post"`
	Content    string `json:"content" example:"Updated content."`
	Published  bool   `json:"published" example:"true"`
}
