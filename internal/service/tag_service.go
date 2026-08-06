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

// HandleCreateTag godoc
//
//	@Summary		Create a new tag
//	@Description	Admin only. Creates a new tag with a name and slug.
//	@Tags			Tags
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.CreateTagRequest	true	"Tag data"
//	@Success		201		{object}	models.CreateTagAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		409		{object}	utils.APIResponse
//	@Router			/api/tags [post]
func HandleCreateTag(w http.ResponseWriter, r *http.Request) {
	repo := repository.TagRepository{
		DB: db.DB,
	}

	var req models.CreateTagRequest

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

	tag, err := repo.CreateTag(req.Name, req.Slug)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Could not create Tag", nil)
		return
	}
	response := models.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	}

	utils.JSON(w, http.StatusCreated, true, "Tag created", response)
}

// HandleGetTags godoc
//
//	@Summary		List all tags
//	@Description	Retrieves all tags.
//	@Tags			Tags
//	@Produce		json
//	@Success		200	{object}	models.ListTagsAPIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/tags [get]
func HandleGetTags(w http.ResponseWriter, r *http.Request) {
	tagRepo := repository.TagRepository{
		DB: db.DB,
	}
	tags, err := tagRepo.GetTags()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal Server Error", nil)
		return
	}
	responses := make([]models.TagResponse, 0, len(tags))
	for _, t := range tags {
		responses = append(responses, models.TagResponse{
			ID:        t.ID,
			Name:      t.Name,
			Slug:      t.Slug,
			CreatedAt: t.CreatedAt,
		})
	}

	utils.JSON(w, http.StatusOK, true, "Tag Retrieved Successfully", responses)

}

// HandleGetTag godoc
//
//	@Summary		Get a tag by ID
//	@Description	Retrieves a single tag by its ID.
//	@Tags			Tags
//	@Produce		json
//	@Param			id	path		int	true	"Tag ID"
//	@Success		200	{object}	models.GetTagAPIResponse
//	@Failure		400	{object}	utils.APIResponse
//	@Failure		404	{object}	utils.APIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/tags/{id} [get]
func HandleGetTag(w http.ResponseWriter, r *http.Request) {
	repo := repository.TagRepository{
		DB: db.DB,
	}

	idString := r.PathValue("id")

	tagID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid Tag ID", nil)
		return
	}
	tag, err := repo.GetTagByID(tagID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if tag == nil {
		utils.JSON(w, http.StatusNotFound, false, "Tag not found", nil)
		return
	}
	response := models.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	}

	utils.JSON(w, http.StatusOK, true, "Tag retrieved successfully", response)

}

// HandleUpdateTag godoc
//
//	@Summary		Update a tag
//	@Description	Admin only. Updates an existing tag's name and slug.
//	@Tags			Tags
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Tag ID"
//	@Param			request	body		models.UpdateTagRequest	true	"Updated tag data"
//	@Success		200		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/tags/{id} [put]
func HandleUpdateTag(w http.ResponseWriter, r *http.Request) {
	repo := repository.TagRepository{
		DB: db.DB,
	}

	var req models.UpdateTagRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid request body", nil)
		return
	}

	idString := r.PathValue("id")

	tagID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid Tag id", nil)
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

	err = repo.UpdateTag(tagID, req.Name, req.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusNotFound, false, "Tag not found", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Tag updated successfully", nil)
}

// HandleDeleteTag godoc
//
//	@Summary		Delete a tag
//	@Description	Admin only. Deletes an existing tag by ID.
//	@Tags			Tags
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Tag ID"
//	@Success		200	{object}	utils.APIResponse
//	@Failure		400	{object}	utils.APIResponse
//	@Failure		404	{object}	utils.APIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/tags/{id} [delete]
func HandleDeleteTag(w http.ResponseWriter, r *http.Request) {
	repo := repository.TagRepository{
		DB: db.DB,
	}
	idString := r.PathValue("id")
	tagID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid tag id", nil)
		return
	}

	err = repo.DeleteTag(tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusNotFound, false, "tag not found", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	utils.JSON(w, http.StatusOK, true, "tag Deleted successfully", nil)

}
