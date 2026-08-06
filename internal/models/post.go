package models

type Post struct {
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
