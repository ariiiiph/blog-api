package service

import (
	"blogapi/internal/db"
	"blogapi/internal/middlewares"
	"blogapi/internal/models"
	"blogapi/internal/repository"
	"blogapi/internal/utils"
	"database/sql"
	"net/http"
	"strconv"
)

// HandleAddTagToPost godoc
//
//	@Summary		Attach a tag to a post
//	@Description	Adds a tag to the specified post. Only the post's author can modify its tags.
//	@Tags			Post Tags
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postID	path		int	true	"Post ID"
//	@Param			tagID	path		int	true	"Tag ID"
//	@Success		201		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		403		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		409		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{postID}/tags/{tagID} [post]
func HandleAddTagToPost(w http.ResponseWriter, r *http.Request) {
	postRepo := repository.PostRepository{DB: db.DB}
	tagRepo := repository.TagRepository{DB: db.DB}
	postTagRepo := repository.PostTagRepository{DB: db.DB}

	postID, err := strconv.ParseInt(r.PathValue("postID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid post ID", nil)
		return
	}

	tagID, err := strconv.ParseInt(r.PathValue("tagID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid tag ID", nil)
		return
	}
	post, err := postRepo.GetPostByID(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if post == nil {
		utils.JSON(w, http.StatusNotFound, false, "could not found post", nil)
		return
	}

	tag, err := tagRepo.GetTagByID(tagID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if tag == nil {
		utils.JSON(w, http.StatusNotFound, false, "could not found tag", nil)
		return
	}
	ctx := r.Context()
	userIdValue := ctx.Value(middlewares.CtxUserID)

	authorId, ok := userIdValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
		return
	}

	if post.AuthorID != authorId {
		utils.JSON(w, http.StatusForbidden, false, "You cannot modify this post", nil)
		return
	}
	err = postTagRepo.AddTagToPost(postID, tagID)
	if err != nil {
		utils.JSON(w, http.StatusConflict, false, "Tag is already attached to this post", nil)
		return
	}

	utils.JSON(w, http.StatusCreated, true, "Tag added to post successfully", nil)
}

// HandleRemoveTagFromPost godoc
//
//	@Summary		Remove a tag from a post
//	@Description	Removes a tag from the specified post. Only the post's author can modify its tags.
//	@Tags			Post Tags
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postID	path		int	true	"Post ID"
//	@Param			tagID	path		int	true	"Tag ID"
//	@Success		200		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		403		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{postID}/tags/{tagID} [delete]
func HandleRemoveTagFromPost(w http.ResponseWriter, r *http.Request) {
	postRepo := repository.PostRepository{DB: db.DB}
	tagRepo := repository.TagRepository{DB: db.DB}
	postTagRepo := repository.PostTagRepository{DB: db.DB}

	postID, err := strconv.ParseInt(r.PathValue("postID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid post ID", nil)
		return
	}

	tagID, err := strconv.ParseInt(r.PathValue("tagID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid tag ID", nil)
		return
	}
	post, err := postRepo.GetPostByID(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if post == nil {
		utils.JSON(w, http.StatusNotFound, false, "could not found post", nil)
		return
	}

	tag, err := tagRepo.GetTagByID(tagID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if tag == nil {
		utils.JSON(w, http.StatusNotFound, false, "could not found tag", nil)
		return
	}
	ctx := r.Context()
	userIdValue := ctx.Value(middlewares.CtxUserID)

	authorId, ok := userIdValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
		return
	}

	if post.AuthorID != authorId {
		utils.JSON(w, http.StatusForbidden, false, "You cannot modify this post", nil)
		return
	}
	err = postTagRepo.RemoveTagFromPost(postID, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusNotFound, false, "Tag is not attached to this post", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Tag Deleted successfully", nil)

}

// HandleGetTagsByPost godoc
//
//	@Summary		Get tags for a post
//	@Description	Retrieves all tags attached to the specified post.
//	@Tags			Post Tags
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	models.TagListAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{postID}/tags [get]
func HandleGetTagsByPost(w http.ResponseWriter, r *http.Request) {
	postRepo := repository.PostRepository{DB: db.DB}
	postTagRepo := repository.PostTagRepository{DB: db.DB}

	postID, err := strconv.ParseInt(r.PathValue("postID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid post ID", nil)
		return
	}

	post, err := postRepo.GetPostByID(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if post == nil {
		utils.JSON(w, http.StatusNotFound, false, "Post not found", nil)
		return
	}

	tags, err := postTagRepo.GetTagsByPost(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
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
	utils.JSON(w, http.StatusOK, true, "Tags retrieved successfully", responses)
}

// HandleGetPostsByTag godoc
//
//	@Summary		Get posts for a tag
//	@Description	Retrieves all posts attached to the specified tag.
//	@Tags			Post Tags
//	@Produce		json
//	@Param			tagID	path		int	true	"Tag ID"
//	@Success		200		{object}	models.PostListByTagAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/tags/{tagID}/posts [get]
func HandleGetPostsByTag(w http.ResponseWriter, r *http.Request) {
	tagRepo := repository.TagRepository{DB: db.DB}
	postTagRepo := repository.PostTagRepository{DB: db.DB}

	tagID, err := strconv.ParseInt(r.PathValue("tagID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid tag ID", nil)
		return
	}

	tag, err := tagRepo.GetTagByID(tagID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if tag == nil {
		utils.JSON(w, http.StatusNotFound, false, "Tag not found", nil)
		return
	}

	posts, err := postTagRepo.GetPostsByTag(tagID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	responses := make([]models.PostResponse, 0, len(posts))
	for _, p := range posts {
		responses = append(responses, models.PostResponse{
			ID:         p.ID,
			AuthorID:   p.AuthorID,
			CategoryID: p.CategoryID,
			Title:      p.Title,
			Slug:       p.Slug,
			Content:    p.Content,
			Published:  p.Published,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		})
	}

	utils.JSON(w, http.StatusOK, true, "Posts retrieved successfully", responses)
}
