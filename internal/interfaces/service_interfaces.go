package interfaces

import (
	"vk/ecom/internal/domain"
	"vk/ecom/internal/dto"
)

type AuthServiceInterface interface {
	RegisterUser(login, password string) (*domain.User, error)
	LoginUser(login, password string) (string, *domain.User, error)
	ValidateToken(tokenString string) (*domain.User, error)
}

type ListingServiceInterface interface {
	CreateListing(req *dto.ListingRequest, authorID int64) (*dto.ListingDTO, error)
	GetListings(sortBy, sortOrder string, minPrice, maxPrice *int64, currentUserID *int64) ([]*dto.ListingDTO, error)
	GetListingsWithPagination(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int, currentUserID *int64) (*dto.ListingsResponse, error)
}
