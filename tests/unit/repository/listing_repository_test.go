package repository_test

import (
	"database/sql"
	"testing"
	"time"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/repository/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestListingRepository_Create(t *testing.T) {
	t.Run("should successfully create listing", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewListingRepository(db)

		listing := &domain.Listing{
			Title:       "Test Listing",
			Description: "Test Description",
			ImageURL:    "http://example.com/image.jpg",
			Price:       100,
			AuthorID:    1,
		}

		expectedID := int64(1)
		mock.ExpectQuery(`INSERT INTO listings \(title, description, image_url, price, author_id, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`).
			WithArgs(listing.Title, listing.Description, listing.ImageURL, listing.Price, listing.AuthorID, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

		err = repo.Create(listing)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, listing.ID)
		assert.False(t, listing.CreatedAt.IsZero())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle database error during create", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewListingRepository(db)

		listing := &domain.Listing{
			Title:       "Test Listing",
			Description: "Test Description",
			Price:       100,
			AuthorID:    1,
		}

		mock.ExpectQuery(`INSERT INTO listings \(title, description, image_url, price, author_id, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`).
			WithArgs(listing.Title, listing.Description, listing.ImageURL, listing.Price, listing.AuthorID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		err = repo.Create(listing)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create listing")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestListingRepository_GetByID(t *testing.T) {
	t.Run("should successfully get listing by ID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewListingRepository(db)

		now := time.Now()
		expectedListing := &domain.Listing{
			ID:          1,
			Title:       "Test Listing",
			Description: "Test Description",
			ImageURL:    "http://example.com/image.jpg",
			Price:       100,
			AuthorID:    1,
			CreatedAt:   now,
		}

		mock.ExpectQuery(`SELECT id, title, description, image_url, price, author_id, created_at FROM listings WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "image_url", "price", "author_id", "created_at"}).
				AddRow(expectedListing.ID, expectedListing.Title, expectedListing.Description, expectedListing.ImageURL,
					expectedListing.Price, expectedListing.AuthorID, expectedListing.CreatedAt))

		listing, err := repo.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, listing)
		assert.Equal(t, expectedListing.ID, listing.ID)
		assert.Equal(t, expectedListing.Title, listing.Title)
		assert.Equal(t, expectedListing.Description, listing.Description)
		assert.Equal(t, expectedListing.ImageURL, listing.ImageURL)
		assert.Equal(t, expectedListing.Price, listing.Price)
		assert.Equal(t, expectedListing.AuthorID, listing.AuthorID)
		assert.Equal(t, expectedListing.CreatedAt.Unix(), listing.CreatedAt.Unix())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when listing not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewListingRepository(db)

		mock.ExpectQuery(`SELECT id, title, description, image_url, price, author_id, created_at FROM listings WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnError(sql.ErrNoRows)

		listing, err := repo.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, listing)
		assert.Contains(t, err.Error(), "listing not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle database error during get", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewListingRepository(db)

		mock.ExpectQuery(`SELECT id, title, description, image_url, price, author_id, created_at FROM listings WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnError(sql.ErrConnDone)

		listing, err := repo.GetByID(1)

		assert.Error(t, err)
		assert.Nil(t, listing)
		assert.Contains(t, err.Error(), "failed to get listing")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
