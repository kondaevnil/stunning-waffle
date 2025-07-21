package domain

import (
	"time"
)

type Listing struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	Price       int64     `json:"price" db:"price"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
