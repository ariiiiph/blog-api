package models

type TagResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

type CreateTagAPIResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    TagResponse `json:"data"`
}

type GetTagAPIResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    TagResponse `json:"data"`
}

type ListTagsAPIResponse struct {
	Status  int           `json:"status"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []TagResponse `json:"data"`
}
