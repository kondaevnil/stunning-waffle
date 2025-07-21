package mocks

import (
	"vk/ecom/internal/dto"
	"vk/ecom/internal/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockListingService struct {
	mock.Mock
}

// Ensure MockListingService implements ListingServiceInterface
var _ interfaces.ListingServiceInterface = (*MockListingService)(nil)

func (m *MockListingService) CreateListing(req *dto.ListingRequest, authorID int64) (*dto.ListingDTO, error) {
	args := m.Called(req, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListingDTO), args.Error(1)
}

func (m *MockListingService) GetListings(sortBy, sortOrder string, minPrice, maxPrice *int64, currentUserID *int64) ([]*dto.ListingDTO, error) {
	args := m.Called(sortBy, sortOrder, minPrice, maxPrice, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.ListingDTO), args.Error(1)
}

func (m *MockListingService) GetListingsWithPagination(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int, currentUserID *int64) (*dto.ListingsResponse, error) {
	args := m.Called(sortBy, sortOrder, minPrice, maxPrice, page, pageSize, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListingsResponse), args.Error(1)
}
