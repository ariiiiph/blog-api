package models

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

type CreateCategoryAPIResponse struct {
	Status  int              `json:"status"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    CategoryResponse `json:"data"`
}

type GetCategoryAPIResponse struct {
	Status  int              `json:"status"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    CategoryResponse `json:"data"`
}

type ListCategoriesAPIResponse struct {
	Status  int                `json:"status"`
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    []CategoryResponse `json:"data"`
}
