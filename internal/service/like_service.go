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
	"strings"
)

// HandleLikePost godoc
//
//	@Summary		Like a post
//	@Description	Adds a like to the specified post from the authenticated user.
//	@Tags			Likes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postID	path		int	true	"Post ID"
//	@Success		201		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		409		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{postID}/like [post]
func HandleLikePost(w http.ResponseWriter, r *http.Request) {
	likeRepo := repository.LikeRepository{DB: db.DB}
	postRepo := repository.PostRepository{DB: db.DB}

	postID, err := strconv.ParseInt(r.PathValue("postID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid post ID", nil)
		return
	}

	userIDValue := r.Context().Value(middlewares.CtxUserID)
	userID, ok := userIDValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
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
	err = likeRepo.LikePost(userID, postID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			utils.JSON(w, http.StatusConflict, false, "Post is already liked", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusCreated, true, "Post liked successfully", nil)
}

// HandleUnlikePost godoc
//
//	@Summary		Unlike a post
//	@Description	Removes the authenticated user's like from the specified post.
//	@Tags			Likes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{postID}/like [delete]
func HandleUnlikePost(w http.ResponseWriter, r *http.Request) {
	likeRepo := repository.LikeRepository{DB: db.DB}
	postRepo := repository.PostRepository{DB: db.DB}

	postID, err := strconv.ParseInt(r.PathValue("postID"), 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid post ID", nil)
		return
	}

	userIDValue := r.Context().Value(middlewares.CtxUserID)
	userID, ok := userIDValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
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

	err = likeRepo.UnlikePost(userID, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusNotFound, false, "Post is not liked", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Post unliked successfully", nil)
}

// HandleGetPostLikes godoc
//
//	@Summary		Get like count for a post
//	@Description	Retrieves the total number of likes for the specified post.
//	@Tags			Likes
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	models.LikeCountAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{postID}/likes [get]
func HandleGetPostLikes(w http.ResponseWriter, r *http.Request) {
	likeRepo := repository.LikeRepository{DB: db.DB}
	postRepo := repository.PostRepository{DB: db.DB}

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

	count, err := likeRepo.CountLikes(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	response := models.LikeCountResponse{
		PostID: postID,
		Likes:  count,
	}

	utils.JSON(w, http.StatusOK, true, "Like count retrieved successfully", response)
}
