package service

import (
	"blogapi/internal/db"
	"blogapi/internal/middlewares"
	"blogapi/internal/models"
	"blogapi/internal/repository"
	"blogapi/internal/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// HandleCreateComment godoc
//
//	@Summary		Add a comment to a post
//	@Description	Creates a new comment on the specified post for the authenticated user.
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Post ID"
//	@Param			request	body		models.CreateCommentRequest	true	"Comment content"
//	@Success		201		{object}	models.CreateCommentAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		401		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{id}/comments [post]
func HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	commentRepo := repository.CommentRepository{
		DB: db.DB,
	}

	postRepo := repository.PostRepository{
		DB: db.DB,
	}

	var req models.CreateCommentRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid request body", nil)
		return
	}
	idString := r.PathValue("id")
	postID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid post ID", nil)
		return
	}

	ctx := r.Context()
	userIDValue := ctx.Value(middlewares.CtxUserID)
	authorID, ok := userIDValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid user", nil)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Content is required", nil)
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

	comment, err := commentRepo.SendComment(postID, authorID, req.Content)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Could not post comment", nil)
		return
	}
	response := models.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		AuthorID:  comment.AuthorID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	utils.JSON(w, http.StatusCreated, true, "comment posted successfully", response)
}

// HandleGetComment godoc
//
//	@Summary		List comments for a post
//	@Description	Retrieves a paginated list of comments for the specified post.
//	@Tags			Comments
//	@Produce		json
//	@Param			id		path		int	true	"Post ID"
//	@Param			page	query		int	false	"Page number"				default(1)
//	@Param			limit	query		int	false	"Items per page (max 50)"	default(10)
//	@Success		200		{object}	models.ListCommentsAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//
// @Router /api/posts/{id}/comments [get]
func HandleGetComment(w http.ResponseWriter, r *http.Request) {
	commentRepo := repository.CommentRepository{
		DB: db.DB,
	}
	postRepo := repository.PostRepository{
		DB: db.DB,
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit

	idString := r.PathValue("id")
	postID, err := strconv.ParseInt(idString, 10, 64)
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

	total, err := commentRepo.CountComments(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	comments, err := commentRepo.GetComments(postID, limit, offset)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	commentResponses := make([]models.CommentResponse, 0, len(comments))
	for _, c := range comments {
		commentResponses = append(commentResponses, models.CommentResponse{
			ID:        c.ID,
			PostID:    c.PostID,
			AuthorID:  c.AuthorID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		})
	}
	totalPages := (total + limit - 1) / limit

	response := models.CommentListResponse{
		Comments:      commentResponses,
		Page:          page,
		Limit:         limit,
		Count:         len(comments),
		TotalComments: total,
		TotalPages:    totalPages,
		HasNext:       page < totalPages,
		HasPrevious:   page > 1,
	}

	utils.JSON(w, http.StatusOK, true, "Comments retrieved successfully", response)
}

// HandleUpdateComment godoc
//
//	@Summary		Update a comment
//	@Description	Updates an existing comment. Only the comment's author can update it.
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Comment ID"
//	@Param			request	body		models.UpdateCommentRequest	true	"Updated comment content"
//	@Success		200		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		403		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/comments/{id} [put]
func HandleUpdateComment(w http.ResponseWriter, r *http.Request) {
	commentRepo := repository.CommentRepository{
		DB: db.DB,
	}

	var req models.UpdateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid body request", nil)
		return
	}

	idString := r.PathValue("id")
	commentID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid comment id", nil)
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Content is required", nil)
		return
	}

	ctx := r.Context()

	userIDValue := ctx.Value(middlewares.CtxUserID)
	authorID, ok := userIDValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
		return
	}

	err = commentRepo.UpdateComment(commentID, authorID, req.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusForbidden, false, "You cannot update this comment", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	utils.JSON(w, http.StatusOK, true, "Comment updated successfully", nil)

}

// HandleDeleteComment godoc
//
//	@Summary		Delete a comment
//	@Description	Deletes an existing comment. Allowed for the comment's author or an admin.
//	@Tags			Comments
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Comment ID"
//	@Success		200	{object}	utils.APIResponse
//	@Failure		400	{object}	utils.APIResponse
//	@Failure		403	{object}	utils.APIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/comments/{id} [delete]
func HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	commentRepo := repository.CommentRepository{
		DB: db.DB,
	}
	idString := r.PathValue("id")
	commentID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid comment id", nil)
		return
	}

	ctx := r.Context()
	userIdValue := ctx.Value(middlewares.CtxUserID)

	authorId, ok := userIdValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
		return
	}
	role, _ := ctx.Value(middlewares.CtxUserRole).(string)
	isAdmin := role == "admin"

	err = commentRepo.DeleteComment(commentID, authorId, isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusForbidden, false, "You cannot delete this comment", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	utils.JSON(w, http.StatusOK, true, "comment Deleted successfully", nil)

}
