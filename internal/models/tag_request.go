package models

type CreateTagRequest struct {
	Name string `json:"name" example:"Golang"`
	Slug string `json:"slug" example:"golang"`
}

type UpdateTagRequest struct {
	Name string `json:"name" example:"Golang"`
	Slug string `json:"slug" example:"golang"`
}
