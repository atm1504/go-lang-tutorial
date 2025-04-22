package models

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title" validate:"required"`
	Artist string  `json:"artist" validate:"required"`
	Price  float32 `json:"price" validate:"required,min=0"`
}
