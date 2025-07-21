package repository

import "vk/ecom/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id int) (*domain.User, error)
	GetByLogin(login string) (*domain.User, error)
}

type ListingRepository interface {
	Create(listing *domain.Listing) error
	GetByID(id int64) (*domain.Listing, error)
	GetAll(sortBy, sortOrder string, minPrice, maxPrice *int64) ([]*domain.Listing, error)
	GetAllWithPagination(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int) ([]*domain.Listing, int, error)
	GetByAuthorID(authorID int64) ([]*domain.Listing, error)
}
