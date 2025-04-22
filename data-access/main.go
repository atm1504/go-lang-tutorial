package main

import (
	"log"
	"net/http"

	"data-access/internal/config"
	"data-access/internal/database"
	"data-access/internal/handlers"
	"data-access/internal/repository"

	"github.com/gorilla/mux"
)

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
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/albums", albumHandler.GetByArtist).Methods("GET")
	router.HandleFunc("/albums/{id:[0-9]+}", albumHandler.GetByID).Methods("GET")
	router.HandleFunc("/albums", albumHandler.Create).Methods("POST")

	// Start server
	log.Printf("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
