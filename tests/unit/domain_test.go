package domain

import (
	"testing"
	"time"
	"vk/ecom/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestUserStruct(t *testing.T) {
	t.Run("should create user with all fields", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Login:    "testuser",
			Password: "hashedpassword",
		}

		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "testuser", user.Login)
		assert.Equal(t, "hashedpassword", user.Password)
	})

	t.Run("should handle empty user", func(t *testing.T) {
		user := &domain.User{}

		assert.Equal(t, 0, user.ID)
		assert.Equal(t, "", user.Login)
		assert.Equal(t, "", user.Password)
	})
}

func TestListingStruct(t *testing.T) {
	t.Run("should create listing with all fields", func(t *testing.T) {
		now := time.Now()
		listing := &domain.Listing{
			ID:          1,
			Title:       "Test Listing",
			Description: "Test Description",
			ImageURL:    "http://example.com/image.jpg",
			Price:       100,
			AuthorID:    1,
			CreatedAt:   now,
		}

		assert.Equal(t, int64(1), listing.ID)
		assert.Equal(t, "Test Listing", listing.Title)
		assert.Equal(t, "Test Description", listing.Description)
		assert.Equal(t, "http://example.com/image.jpg", listing.ImageURL)
		assert.Equal(t, int64(100), listing.Price)
		assert.Equal(t, int64(1), listing.AuthorID)
		assert.Equal(t, now, listing.CreatedAt)
	})

	t.Run("should handle empty listing", func(t *testing.T) {
		listing := &domain.Listing{}

		assert.Equal(t, int64(0), listing.ID)
		assert.Equal(t, "", listing.Title)
		assert.Equal(t, "", listing.Description)
		assert.Equal(t, "", listing.ImageURL)
		assert.Equal(t, int64(0), listing.Price)
		assert.Equal(t, int64(0), listing.AuthorID)
		assert.True(t, listing.CreatedAt.IsZero())
	})
}
