package service

import (
	"blogapi/internal/db"
	"blogapi/internal/models"
	"blogapi/internal/repository"
	"blogapi/internal/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// HandleCreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	Admin only. Creates a new category with a name and slug.
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.CreateCategoryRequest	true	"Category data"
//	@Success		201		{object}	models.CreateCategoryAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		409		{object}	utils.APIResponse
//	@Router			/api/categories [post]
func HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	repo := repository.CategoryRepository{
		DB: db.DB,
	}

	var req models.CreateCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid body request", nil)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Name is required", nil)
		return
	}
	if strings.TrimSpace(req.Slug) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Slug is required", nil)
		return
	}

	category, err := repo.CreateCategory(req.Name, req.Slug)
	if err != nil {
		utils.JSON(w, http.StatusConflict, false, "Could not create category", nil)
		return
	}
	response := models.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
	}
	utils.JSON(w, http.StatusCreated, true, "category created", response)
}

// HandleGetCategories godoc
//
//	@Summary		List all categories
//	@Description	Retrieves all categories.
//	@Tags			Categories
//	@Produce		json
//	@Success		200	{object}	models.ListCategoriesAPIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/categories [get]
func HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categoryRepo := repository.CategoryRepository{
		DB: db.DB,
	}
	categories, err := categoryRepo.GetCategories()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal Server Error", nil)
		return
	}
	responses := make([]models.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		responses = append(responses, models.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			CreatedAt: c.CreatedAt,
		})
	}

	utils.JSON(w, http.StatusOK, true, "Category Retrieved Successfully", responses)

}

// HandleGetCategory godoc
//
//	@Summary		Get a category by ID
//	@Description	Retrieves a single category by its ID.
//	@Tags			Categories
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	models.GetCategoryAPIResponse
//	@Failure		400	{object}	utils.APIResponse
//	@Failure		404	{object}	utils.APIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/categories/{id} [get]
func HandleGetCategory(w http.ResponseWriter, r *http.Request) {
	repo := repository.CategoryRepository{
		DB: db.DB,
	}

	idString := r.PathValue("id")

	categoryID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid Category ID", nil)
		return
	}
	category, err := repo.GetCategoryByID(categoryID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if category == nil {
		utils.JSON(w, http.StatusNotFound, false, "Category not found", nil)
		return
	}
	response := models.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
	}

	utils.JSON(w, http.StatusOK, true, "Category retrieved successfully", response)

}

// HandleUpdateCategory godoc
//
//	@Summary		Update a category
//	@Description	Admin only. Updates an existing category's name and slug.
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"Category ID"
//	@Param			request	body		models.UpdateCategoryRequest	true	"Updated category data"
//	@Success		200		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/categories/{id} [put]
func HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	repo := repository.CategoryRepository{
		DB: db.DB,
	}

	var req models.UpdateCategoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid request body", nil)
		return
	}

	idString := r.PathValue("id")

	categoryID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid category id", nil)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "name is required", nil)
		return
	}
	if strings.TrimSpace(req.Slug) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "slug is required", nil)
		return
	}

	err = repo.UpdateCategory(categoryID, req.Name, req.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusNotFound, false, "Category not found", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "category updated successfully", nil)
}

// HandleDeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Admin only. Deletes an existing category by ID.
//	@Tags			Categories
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	utils.APIResponse
//	@Failure		400	{object}	utils.APIResponse
//	@Failure		404	{object}	utils.APIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/categories/{id} [delete]
func HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	repo := repository.CategoryRepository{
		DB: db.DB,
	}
	idString := r.PathValue("id")
	categoryID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid category id", nil)
		return
	}

	err = repo.DeleteCategory(categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusNotFound, false, "category not found", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	utils.JSON(w, http.StatusOK, true, "category Deleted successfully", nil)

}
