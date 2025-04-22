package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"data-access/internal/models"
	"data-access/internal/repository"

	"github.com/gorilla/mux"
)

// AlbumHandler handles HTTP requests for albums
type AlbumHandler struct {
	repo *repository.AlbumRepository
}

// NewAlbumHandler creates a new album handler
func NewAlbumHandler(repo *repository.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{repo: repo}
}

// GetByArtist godoc
// @Summary Get albums by artist name
// @Description Retrieve all albums by a specific artist
// @Tags albums
// @Accept json
// @Produce json
// @Param artist query string true "Artist name"
// @Success 200 {array} models.Album
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /albums [get]
func (h *AlbumHandler) GetByArtist(w http.ResponseWriter, r *http.Request) {
	artist := r.URL.Query().Get("artist")
	if artist == "" {
		http.Error(w, "artist parameter is required", http.StatusBadRequest)
		return
	}

	albums, err := h.repo.GetByArtist(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

// GetByID godoc
// @Summary Get album by ID
// @Description Retrieve a specific album by its ID
// @Tags albums
// @Accept json
// @Produce json
// @Param id path int true "Album ID"
// @Success 200 {object} models.Album
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /albums/{id} [get]
func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid album id", http.StatusBadRequest)
		return
	}

	album, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

// Create godoc
// @Summary Create a new album
// @Description Create a new album with the provided details
// @Tags albums
// @Accept json
// @Produce json
// @Param album body models.Album true "Album object"
// @Success 201 {object} map[string]int64
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /albums [post]
func (h *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {
	var album models.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(&album)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}
