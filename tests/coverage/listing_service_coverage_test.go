package coverage_test

import (
	"testing"
	"time"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/dto"
	"vk/ecom/internal/mocks"
	"vk/ecom/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListingService_Coverage(t *testing.T) {
	t.Run("should create listing successfully", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		req := &dto.ListingRequest{
			Title:       "Test Listing",
			Description: "This is a test listing description with enough characters",
			ImageURL:    "http://example.com/image.jpg",
			Price:       100,
		}
		authorID := int64(1)

		author := &domain.User{ID: 1, Login: "testuser"}

		mockListingRepo.On("Create", mock.AnythingOfType("*domain.Listing")).Return(nil)
		mockUserRepo.On("GetByID", 1).Return(author, nil)

		result, err := listingService.CreateListing(req, authorID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Title, result.Title)
		assert.Equal(t, req.Description, result.Description)
		assert.Equal(t, req.ImageURL, result.ImageURL)
		assert.Equal(t, req.Price, result.Price)
		assert.Equal(t, authorID, result.AuthorID)
		assert.Equal(t, author.Login, result.AuthorLogin)

		mockListingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should get listings with pagination", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		now := time.Now()
		listings := []*domain.Listing{
			{
				ID:          1,
				Title:       "Listing 1",
				Description: "Description 1",
				Price:       100,
				AuthorID:    1,
				CreatedAt:   now,
			},
		}
		totalCount := 1

		user1 := &domain.User{ID: 1, Login: "user1"}

		mockListingRepo.On("GetAllWithPagination", "date", "desc", (*int64)(nil), (*int64)(nil), 1, 10).Return(listings, totalCount, nil)
		mockUserRepo.On("GetByID", 1).Return(user1, nil)

		result, err := listingService.GetListingsWithPagination("date", "desc", nil, nil, 1, 10, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Listings, 1)
		assert.Equal(t, 1, result.Count)
		assert.Equal(t, 1, result.TotalCount)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 1, result.TotalPages)

		mockListingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
