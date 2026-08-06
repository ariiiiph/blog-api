package main

import (
	_ "blogapi/docs"
	"blogapi/internal/config"
	"blogapi/internal/db"
	"blogapi/internal/middlewares"
	"blogapi/internal/routes"
	"blogapi/internal/utils"
	"log"
	"net/http"
)

//	@title			Blog API
//	@version		1.0
//	@description	This is a sample Blog.
//	@host			localhost:8080

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	cfg := config.LoadConfig()
	utils.InitJWT(cfg.JWTKey)

	db.InitDB(cfg.DBPath, cfg.DBName)
	defer db.CloseDB()

	mux := routes.RegisterRoutes(db.DB)

	loggerMux := middlewares.LoggingMiddleware(mux)
	corsMux := middlewares.CorsMiddleware(loggerMux)
	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: corsMux,
	}
	log.Printf("Server listening on http://%s", server.Addr)
	log.Printf("Swagger UI available at http://%s/swagger/", server.Addr)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed to start: %v", err)
	}
}
