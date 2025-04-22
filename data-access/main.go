package main

import (
	"log"
	"net/http"

	_ "data-access/docs" // This is where the generated docs are
	"data-access/internal/config"
	"data-access/internal/database"
	"data-access/internal/handlers"
	"data-access/internal/repository"
	"data-access/internal/routes"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Album API
// @version 1.0
// @description This is a sample album service API
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	albumRepo := repository.NewAlbumRepository(db.DB)

	// Initialize handler
	albumHandler := handlers.NewAlbumHandler(albumRepo)

	// Initialize router
	router := routes.NewRouter(albumHandler)

	// Add Swagger documentation route
	router.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start server
	log.Printf("Server starting on port 8080...")
	log.Printf("Swagger documentation available at http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
