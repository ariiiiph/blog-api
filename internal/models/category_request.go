package models

type CreateCategoryRequest struct {
	Name string `json:"name" example:"Technology"`
	Slug string `json:"slug" example:"technology"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" example:"Technology"`
	Slug string `json:"slug" example:"technology"`
}
