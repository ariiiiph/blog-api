package models

type TagListAPIResponse struct {
	Status  int           `json:"status"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []TagResponse `json:"data"`
}

type PostListByTagAPIResponse struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []PostResponse `json:"data"`
}
