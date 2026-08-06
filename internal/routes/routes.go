package routes

import (
	"blogapi/internal/middlewares"
	"blogapi/internal/service"
	"database/sql"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(conn *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Swagger
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	//Users
	mux.HandleFunc("POST /api/users/register", service.HandleRegister)
	mux.HandleFunc("POST /api/users/login", service.HandleLogin)
	mux.HandleFunc("POST /api/users/refresh", service.HandleRefreshToken)
	mux.HandleFunc("POST /api/users/logout", service.HandleLogout)

	//Posts
	mux.HandleFunc("POST /api/posts", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleCreatePost)))
	mux.HandleFunc("GET /api/posts/", service.HandleGetPosts)
	mux.HandleFunc("GET /api/posts/{id}", service.HandleGetPostById)
	mux.HandleFunc("PUT /api/posts/{id}", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleUpdatePost)))
	mux.HandleFunc("DELETE /api/posts/{id}", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleDeletePost)))

	//comments
	mux.HandleFunc("POST /api/posts/{id}/comments", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleCreateComment)))
	mux.HandleFunc("GET /api/posts/{id}/comments", service.HandleGetComment)
	mux.HandleFunc("PUT /api/comments/{id}", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleUpdateComment)))
	mux.HandleFunc("DELETE /api/comments/{id}", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleDeleteComment)))

	//Category
	mux.HandleFunc("POST /api/categories", middlewares.AdminOnlyHandler(http.HandlerFunc(service.HandleCreateCategory)))
	mux.HandleFunc("GET /api/categories", service.HandleGetCategories)
	mux.HandleFunc("GET /api/categories/{id}", service.HandleGetCategory)
	mux.HandleFunc("PUT /api/categories/{id}", middlewares.AdminOnlyHandler(http.HandlerFunc(service.HandleUpdateCategory)))
	mux.HandleFunc("DELETE /api/categories/{id}", middlewares.AdminOnlyHandler(http.HandlerFunc(service.HandleDeleteCategory)))

	//Tags
	mux.HandleFunc("POST /api/tags", middlewares.AdminOnlyHandler(http.HandlerFunc(service.HandleCreateTag)))
	mux.HandleFunc("GET /api/tags", service.HandleGetTags)
	mux.HandleFunc("GET /api/tags/{id}", service.HandleGetTag)
	mux.HandleFunc("PUT /api/tags/{id}", middlewares.AdminOnlyHandler(http.HandlerFunc(service.HandleUpdateTag)))
	mux.HandleFunc("DELETE /api/tags/{id}", middlewares.AdminOnlyHandler(http.HandlerFunc(service.HandleDeleteTag)))

	// Post Tags
	mux.HandleFunc("POST /api/posts/{postID}/tags/{tagID}", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleAddTagToPost)))
	mux.HandleFunc("DELETE /api/posts/{postID}/tags/{tagID}", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleRemoveTagFromPost)))
	mux.HandleFunc("GET /api/posts/{postID}/tags", service.HandleGetTagsByPost)
	mux.HandleFunc("GET /api/tags/{tagID}/posts", service.HandleGetPostsByTag)

	//likes
	mux.HandleFunc("POST /api/posts/{postID}/like", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleLikePost)))
	mux.HandleFunc("DELETE /api/posts/{postID}/like", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleUnlikePost)))
	mux.HandleFunc("GET /api/posts/{postID}/likes", service.HandleGetPostLikes)

	//Posts
	mux.HandleFunc("POST /api/posts/{id}/image", middlewares.AuthenticateHandler(http.HandlerFunc(service.HandleUploadPostImage)))
	// Static file serving for uploaded images
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	return mux
}
