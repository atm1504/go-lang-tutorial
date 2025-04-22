package repository

import (
	"database/sql"
	"fmt"

	"data-access/internal/models"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) GetByArtist(name string) ([]models.Album, error) {
	var albums []models.Album

	rows, err := r.db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("error querying albums by artist: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb models.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("error scanning album: %v", err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating albums: %v", err)
	}

	return albums, nil
}

func (r *AlbumRepository) GetByID(id int64) (*models.Album, error) {
	album := &models.Album{}
	err := r.db.QueryRow("SELECT * FROM album WHERE id = ?", id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("album not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying album: %v", err)
	}
	return album, nil
}

func (r *AlbumRepository) Create(album *models.Album) (int64, error) {
	result, err := r.db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)",
		album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("error creating album: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %v", err)
	}

	return id, nil
}
