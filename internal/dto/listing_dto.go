package dto

import (
	"time"
	"vk/ecom/internal/domain"
)

type ListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Price       int64  `json:"price"`
}

type ListingDTO struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	Price        int64     `json:"price"`
	AuthorID     int64     `json:"author_id"`
	AuthorLogin  string    `json:"author_login"`
	CreatedAt    time.Time `json:"created_at"`
	IsOwnListing *bool     `json:"is_own_listing,omitempty"`
}

type ListingsResponse struct {
	Listings   []*ListingDTO `json:"listings"`
	Count      int           `json:"count"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

func ToListingDTO(listing *domain.Listing) *ListingDTO {
	if listing == nil {
		return nil
	}
	return &ListingDTO{
		ID:          listing.ID,
		Title:       listing.Title,
		Description: listing.Description,
		ImageURL:    listing.ImageURL,
		Price:       listing.Price,
		AuthorID:    listing.AuthorID,
		CreatedAt:   listing.CreatedAt,
	}
}

func ToListingDTOWithAuthor(listing *domain.Listing, authorLogin string, currentUserID *int64) *ListingDTO {
	if listing == nil {
		return nil
	}

	dto := &ListingDTO{
		ID:          listing.ID,
		Title:       listing.Title,
		Description: listing.Description,
		ImageURL:    listing.ImageURL,
		Price:       listing.Price,
		AuthorID:    listing.AuthorID,
		AuthorLogin: authorLogin,
		CreatedAt:   listing.CreatedAt,
	}

	if currentUserID != nil {
		isOwn := listing.AuthorID == *currentUserID
		dto.IsOwnListing = &isOwn
	}

	return dto
}
