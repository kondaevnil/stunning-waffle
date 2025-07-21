package mocks

import (
	"vk/ecom/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockListingRepository struct {
	mock.Mock
}

func (m *MockListingRepository) Create(listing *domain.Listing) error {
	args := m.Called(listing)
	return args.Error(0)
}

func (m *MockListingRepository) GetByID(id int64) (*domain.Listing, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

func (m *MockListingRepository) GetAll(sortBy, sortOrder string, minPrice, maxPrice *int64) ([]*domain.Listing, error) {
	args := m.Called(sortBy, sortOrder, minPrice, maxPrice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Listing), args.Error(1)
}

func (m *MockListingRepository) GetAllWithPagination(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int) ([]*domain.Listing, int, error) {
	args := m.Called(sortBy, sortOrder, minPrice, maxPrice, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Listing), args.Int(1), args.Error(2)
}

func (m *MockListingRepository) GetByAuthorID(authorID int64) ([]*domain.Listing, error) {
	args := m.Called(authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Listing), args.Error(1)
}
