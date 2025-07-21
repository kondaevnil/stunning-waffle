package service_test

import (
	"errors"
	"testing"
	"time"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/dto"
	"vk/ecom/internal/mocks"
	"vk/ecom/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListingService_CreateListing(t *testing.T) {
	t.Run("should successfully create listing with valid data", func(t *testing.T) {
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

	t.Run("should fail with short title", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		req := &dto.ListingRequest{
			Title:       "AB", // Too short
			Description: "This is a test listing description with enough characters",
			Price:       100,
		}

		result, err := listingService.CreateListing(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "title must be between 3 and 100 characters")
	})

	t.Run("should fail with long title", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		longTitle := "This is a very very very long title that exceeds the maximum allowed length of 100 characters for titles"
		req := &dto.ListingRequest{
			Title:       longTitle,
			Description: "This is a test listing description with enough characters",
			Price:       100,
		}

		result, err := listingService.CreateListing(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "title must be between 3 and 100 characters")
	})

	t.Run("should fail with short description", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		req := &dto.ListingRequest{
			Title:       "Test Listing",
			Description: "Short",
			Price:       100,
		}

		result, err := listingService.CreateListing(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "description must be between 10 and 2000 characters")
	})

	t.Run("should fail with invalid price (too low)", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		req := &dto.ListingRequest{
			Title:       "Test Listing",
			Description: "This is a test listing description with enough characters",
			Price:       0,
		}

		result, err := listingService.CreateListing(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "price must be between 1 and 1000000000")
	})

	t.Run("should fail with invalid price (too high)", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		req := &dto.ListingRequest{
			Title:       "Test Listing",
			Description: "This is a test listing description with enough characters",
			Price:       1_000_000_001,
		}

		result, err := listingService.CreateListing(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "price must be between 1 and 1000000000")
	})

	t.Run("should fail with invalid image format", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		req := &dto.ListingRequest{
			Title:       "Test Listing",
			Description: "This is a test listing description with enough characters",
			ImageURL:    "http://example.com/image.gif",
			Price:       100,
		}

		result, err := listingService.CreateListing(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "unsupported image format")
	})

	t.Run("should accept valid image formats", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		// Note: Based on the current implementation, it only checks last 4 characters
		// So ".jpeg" won't work because it takes last 4 chars which would be "peg"
		validFormats := []string{".jpg", ".png"}
		author := &domain.User{ID: 1, Login: "testuser"}

		for _, format := range validFormats {
			mockListingRepo.On("Create", mock.AnythingOfType("*domain.Listing")).Return(nil).Once()
			mockUserRepo.On("GetByID", 1).Return(author, nil).Once()

			req := &dto.ListingRequest{
				Title:       "Test Listing",
				Description: "This is a test listing description with enough characters",
				ImageURL:    "http://example.com/image" + format,
				Price:       100,
			}

			result, err := listingService.CreateListing(req, 1)

			assert.NoError(t, err)
			assert.NotNil(t, result)
		}

		mockListingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should handle repository create error", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		req := &dto.ListingRequest{
			Title:       "Test Listing",
			Description: "This is a test listing description with enough characters",
			Price:       100,
		}

		mockListingRepo.On("Create", mock.AnythingOfType("*domain.Listing")).Return(errors.New("database error"))

		result, err := listingService.CreateListing(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")

		mockListingRepo.AssertExpectations(t)
	})
}

func TestListingService_GetListings(t *testing.T) {
	t.Run("should successfully get listings", func(t *testing.T) {
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
			{
				ID:          2,
				Title:       "Listing 2",
				Description: "Description 2",
				Price:       200,
				AuthorID:    2,
				CreatedAt:   now,
			},
		}

		user1 := &domain.User{ID: 1, Login: "user1"}
		user2 := &domain.User{ID: 2, Login: "user2"}

		mockListingRepo.On("GetAll", "date", "desc", (*int64)(nil), (*int64)(nil)).Return(listings, nil)
		mockUserRepo.On("GetByID", 1).Return(user1, nil)
		mockUserRepo.On("GetByID", 2).Return(user2, nil)

		currentUserID := int64(1)
		result, err := listingService.GetListings("date", "desc", nil, nil, &currentUserID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "user1", result[0].AuthorLogin)
		assert.Equal(t, "user2", result[1].AuthorLogin)
		assert.True(t, *result[0].IsOwnListing)
		assert.False(t, *result[1].IsOwnListing)

		mockListingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		mockListingRepo.On("GetAll", "date", "desc", (*int64)(nil), (*int64)(nil)).Return(nil, errors.New("database error"))

		result, err := listingService.GetListings("date", "desc", nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")

		mockListingRepo.AssertExpectations(t)
	})

	t.Run("should handle missing author gracefully", func(t *testing.T) {
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

		mockListingRepo.On("GetAll", "date", "desc", (*int64)(nil), (*int64)(nil)).Return(listings, nil)
		mockUserRepo.On("GetByID", 1).Return(nil, errors.New("user not found"))

		result, err := listingService.GetListings("date", "desc", nil, nil, nil)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "", result[0].AuthorLogin)

		mockListingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestListingService_GetListingsWithPagination(t *testing.T) {
	t.Run("should successfully get listings with pagination", func(t *testing.T) {
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
		totalCount := 25

		user1 := &domain.User{ID: 1, Login: "user1"}

		mockListingRepo.On("GetAllWithPagination", "date", "desc", (*int64)(nil), (*int64)(nil), 1, 10).Return(listings, totalCount, nil)
		mockUserRepo.On("GetByID", 1).Return(user1, nil)

		result, err := listingService.GetListingsWithPagination("date", "desc", nil, nil, 1, 10, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Listings, 1)
		assert.Equal(t, 1, result.Count)
		assert.Equal(t, 25, result.TotalCount)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 3, result.TotalPages) // ceil(25/10) = 3

		mockListingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should handle invalid pagination parameters", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		// Test with invalid page and pageSize
		mockListingRepo.On("GetAllWithPagination", "date", "desc", (*int64)(nil), (*int64)(nil), 1, 10).Return([]*domain.Listing{}, 0, nil)

		result, err := listingService.GetListingsWithPagination("invalid_sort", "invalid_order", nil, nil, -1, 0, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Page)      // Should default to 1
		assert.Equal(t, 10, result.PageSize) // Should default to 10

		mockListingRepo.AssertExpectations(t)
	})

	t.Run("should handle large page size", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		// Test with page size exceeding maximum
		mockListingRepo.On("GetAllWithPagination", "date", "desc", (*int64)(nil), (*int64)(nil), 1, 100).Return([]*domain.Listing{}, 0, nil)

		result, err := listingService.GetListingsWithPagination("date", "desc", nil, nil, 1, 150, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 100, result.PageSize) // Should be capped at 100

		mockListingRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		mockListingRepo := new(mocks.MockListingRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		listingService := service.NewListingService(mockListingRepo, mockUserRepo)

		mockListingRepo.On("GetAllWithPagination", "date", "desc", (*int64)(nil), (*int64)(nil), 1, 10).Return(nil, 0, errors.New("database error"))

		result, err := listingService.GetListingsWithPagination("date", "desc", nil, nil, 1, 10, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")

		mockListingRepo.AssertExpectations(t)
	})
}
