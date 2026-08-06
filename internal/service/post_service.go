package service

import (
	"blogapi/internal/db"
	"blogapi/internal/middlewares"
	"blogapi/internal/models"
	"blogapi/internal/repository"
	"blogapi/internal/utils"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// HandleCreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Creates a new blog post for the authenticated user.
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.CreatePostRequest	true	"Post data"
//	@Success		201		{object}	models.CreatePostAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		401		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts [post]
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	repo := repository.PostRepository{
		DB: db.DB,
	}
	categoryRepo := repository.CategoryRepository{
		DB: db.DB,
	}

	var req models.CreatePostRequest

	published := true

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid request body", nil)
		return
	}

	ctx := r.Context()
	userIDValue := ctx.Value(middlewares.CtxUserID)
	authorID, ok := userIDValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid user", nil)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Title is required", nil)
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Content is required", nil)
		return
	}
	if strings.TrimSpace(req.Slug) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Slug is required", nil)
		return
	}
	if req.CategoryID == 0 {
		utils.JSON(w, http.StatusBadRequest, false, "Category is required", nil)
		return
	}
	category, err := categoryRepo.GetCategoryByID(req.CategoryID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if category == nil {
		utils.JSON(w, http.StatusNotFound, false, "Category not found", nil)
		return
	}

	post, err := repo.CreatePost(authorID, req.CategoryID, req.Title, req.Slug, req.Content, published)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Could not create post", nil)
		return
	}
	response := models.PostResponse{
		ID:         post.ID,
		AuthorID:   post.AuthorID,
		CategoryID: post.CategoryID,
		Title:      post.Title,
		Slug:       post.Slug,
		Content:    post.Content,
		Published:  post.Published,
		ImageURL:   post.ImageURL,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}

	utils.JSON(w, http.StatusCreated, true, "Post created successfully", response)

}

// HandleGetPosts godoc
//
//	@Summary		List posts
//	@Description	Retrieves a paginated list of posts, optionally filtered by search term.
//	@Tags			Posts
//	@Produce		json
//	@Param			page	query		int		false	"Page number"				default(1)
//	@Param			limit	query		int		false	"Items per page (max 50)"	default(10)
//	@Param			search	query		string	false	"Search term"
//	@Success		200		{object}	models.ListPostsAPIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts [get]
func HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	repo := repository.PostRepository{
		DB: db.DB,
	}

	// Parse and validate page
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Parse and validate limit
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))

	offset := (page - 1) * limit

	total, err := repo.CountPosts(search)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	posts, err := repo.GetPosts(limit, offset, search)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	postResponses := make([]models.PostResponse, 0, len(posts))
	for _, p := range posts {
		postResponses = append(postResponses, models.PostResponse{
			ID:         p.ID,
			AuthorID:   p.AuthorID,
			CategoryID: p.CategoryID,
			Title:      p.Title,
			Slug:       p.Slug,
			Content:    p.Content,
			Published:  p.Published,
			ImageURL:   p.ImageURL,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		})
	}
	totalPages := (total + limit - 1) / limit

	response := models.PostListResponse{
		Posts:       postResponses,
		Search:      search,
		Page:        page,
		Limit:       limit,
		Count:       len(posts),
		TotalPosts:  total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	utils.JSON(w, http.StatusOK, true, "Posts retrieved successfully", response)
}

// HandleGetPostById godoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieves a single post by its ID.
//	@Tags			Posts
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	models.GetPostAPIResponse
//	@Failure		400	{object}	utils.APIResponse
//	@Failure		404	{object}	utils.APIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/posts/{id} [get]
func HandleGetPostById(w http.ResponseWriter, r *http.Request) {
	repo := repository.PostRepository{
		DB: db.DB,
	}

	idString := r.PathValue("id")

	postID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid post ID", nil)
		return
	}

	post, err := repo.GetPostByID(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if post == nil {
		utils.JSON(w, http.StatusNotFound, false, "Post not found", nil)
		return
	}
	response := models.PostResponse{
		ID:         post.ID,
		AuthorID:   post.AuthorID,
		CategoryID: post.CategoryID,
		Title:      post.Title,
		Slug:       post.Slug,
		Content:    post.Content,
		Published:  post.Published,
		ImageURL:   post.ImageURL,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}

	utils.JSON(w, http.StatusOK, true, "Post retrieved successfully", response)

}

// HandleUpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Updates an existing post. Only the post's author can update it.
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Post ID"
//	@Param			request	body		models.UpdatePostRequest	true	"Updated post data"
//	@Success		200		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		403		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{id} [put]
func HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	repo := repository.PostRepository{
		DB: db.DB,
	}
	categoryRepo := repository.CategoryRepository{
		DB: db.DB,
	}
	var req models.UpdatePostRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid request body", nil)
		return
	}

	idString := r.PathValue("id")

	postId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid post id", nil)
		return
	}
	post, err := repo.GetPostByID(postId)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if post == nil {
		utils.JSON(w, http.StatusNotFound, false, "Post not found", nil)
		return
	}

	ctx := r.Context()

	if strings.TrimSpace(req.Title) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Title is required", nil)
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Content is required", nil)
		return
	}
	if strings.TrimSpace(req.Slug) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Slug is required", nil)
		return
	}
	if req.CategoryID == 0 {
		utils.JSON(w, http.StatusBadRequest, false, "Category is required", nil)
		return
	}

	userIdValue := ctx.Value(middlewares.CtxUserID)

	authorId, ok := userIdValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
		return
	}
	category, err := categoryRepo.GetCategoryByID(req.CategoryID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if category == nil {
		utils.JSON(w, http.StatusNotFound, false, "Category not found", nil)
		return
	}

	err = repo.UpdatePost(postId, authorId, req.CategoryID, req.Title, req.Slug, req.Content, req.Published)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusForbidden, false, "You cannot update this post", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Post updated successfully", nil)
}

// HandleDeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Deletes an existing post. Allowed for the post's author or an admin.
//	@Tags			Posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	utils.APIResponse
//	@Failure		400	{object}	utils.APIResponse
//	@Failure		403	{object}	utils.APIResponse
//	@Failure		500	{object}	utils.APIResponse
//	@Router			/api/posts/{id} [delete]
func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	repo := repository.PostRepository{
		DB: db.DB,
	}

	idString := r.PathValue("id")

	postId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid post id", nil)
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

	err = repo.DeletePost(postId, authorId, isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusForbidden, false, "You cannot delete this post", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Post Deleted successfully", nil)

}

// HandleUploadPostImage godoc
//
//	@Summary		Upload a post's cover image
//	@Description	Uploads a cover image for the specified post. Allowed for the post's author or an admin. Replaces any existing image.
//	@Tags			Posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int		true	"Post ID"
//	@Param			image	formData	file	true	"Image file (jpg, jpeg, png, webp; max 5MB)"
//	@Success		200		{object}	models.UploadPostImageAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		403		{object}	utils.APIResponse
//	@Failure		404		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/posts/{id}/image [post]
func HandleUploadPostImage(w http.ResponseWriter, r *http.Request) {
	repo := repository.PostRepository{
		DB: db.DB,
	}

	idString := r.PathValue("id")
	postID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid post id", nil)
		return
	}

	post, err := repo.GetPostByID(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if post == nil {
		utils.JSON(w, http.StatusNotFound, false, "Post not found", nil)
		return
	}

	ctx := r.Context()
	userIDValue := ctx.Value(middlewares.CtxUserID)
	authorID, ok := userIDValue.(int64)
	if !ok {
		utils.JSON(w, http.StatusInternalServerError, false, "Invalid user context", nil)
		return
	}

	role, _ := ctx.Value(middlewares.CtxUserRole).(string)
	isAdmin := role == "admin"

	if post.AuthorID != authorID && !isAdmin {
		utils.JSON(w, http.StatusForbidden, false, "You cannot modify this post", nil)
		return
	}

	err = r.ParseMultipartForm(utils.MaxImageSize)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "invalid form data or file too large", nil)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "image file is required", nil)
		return
	}
	defer file.Close()

	ext, err := utils.ValidateImageFile(fileHeader)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	filename := utils.GenerateImageFilename(ext)
	uploadDir := "uploads/posts"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "could not create upload directory", nil)
		return
	}

	destPath := filepath.Join(uploadDir, filename)
	dest, err := os.Create(destPath)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "could not save file", nil)
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "could not save file", nil)
		return
	}

	newImageURL := "/uploads/posts/" + filename

	// Delete old image file if one exists, before committing the new URL.
	if post.ImageURL != nil && *post.ImageURL != "" {
		oldPath := strings.TrimPrefix(*post.ImageURL, "/")
		_ = os.Remove(oldPath) // best-effort; ignore error if file already missing
	}

	if err := repo.UpdatePostImage(postID, newImageURL); err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "could not update post image", nil)
		return
	}

	updatedPost, err := repo.GetPostByID(postID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	response := models.PostResponse{
		ID:         updatedPost.ID,
		AuthorID:   updatedPost.AuthorID,
		CategoryID: updatedPost.CategoryID,
		Title:      updatedPost.Title,
		Slug:       updatedPost.Slug,
		Content:    updatedPost.Content,
		Published:  updatedPost.Published,
		ImageURL:   updatedPost.ImageURL,
		CreatedAt:  updatedPost.CreatedAt,
		UpdatedAt:  updatedPost.UpdatedAt,
	}

	utils.JSON(w, http.StatusOK, true, "Post image uploaded successfully", response)
}
