package models

type PostResponse struct {
	ID         int64   `json:"id"`
	AuthorID   int64   `json:"author_id"`
	CategoryID int64   `json:"category_id"`
	Title      string  `json:"title"`
	Slug       string  `json:"slug"`
	Content    string  `json:"content"`
	Published  bool    `json:"published"`
	ImageURL   *string `json:"image_url"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type PostListResponse struct {
	Posts       []PostResponse `json:"posts"`
	Search      string         `json:"search"`
	Page        int            `json:"page"`
	Limit       int            `json:"limit"`
	Count       int            `json:"count"`
	TotalPosts  int            `json:"total_posts"`
	TotalPages  int            `json:"total_pages"`
	HasNext     bool           `json:"has_next"`
	HasPrevious bool           `json:"has_previous"`
}

type CreatePostAPIResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    PostResponse `json:"data"`
}

type GetPostAPIResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    PostResponse `json:"data"`
}

type ListPostsAPIResponse struct {
	Status  int              `json:"status"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    PostListResponse `json:"data"`
}

type UploadPostImageAPIResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    PostResponse `json:"data"`
}
