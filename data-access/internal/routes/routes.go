package routes

import (
	"net/http"

	"data-access/internal/handlers"

	"github.com/gorilla/mux"
)

// Router represents the application router
type Router struct {
	Router *mux.Router
}

// NewRouter creates a new router instance
func NewRouter(albumHandler *handlers.AlbumHandler) *Router {
	r := &Router{
		Router: mux.NewRouter(),
	}

	// Register routes
	r.registerAlbumRoutes(albumHandler)

	return r
}

// registerAlbumRoutes registers all album-related routes
func (r *Router) registerAlbumRoutes(h *handlers.AlbumHandler) {
	// Album routes
	albums := r.Router.PathPrefix("/albums").Subrouter()
	albums.HandleFunc("", h.GetByArtist).Methods(http.MethodGet)
	albums.HandleFunc("/{id:[0-9]+}", h.GetByID).Methods(http.MethodGet)
	albums.HandleFunc("", h.Create).Methods(http.MethodPost)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}
