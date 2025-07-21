package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"vk/ecom/internal/domain"
)

type ListingRepository struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) Create(listing *domain.Listing) error {
	query := `
		INSERT INTO listings (title, description, image_url, price, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	listing.CreatedAt = time.Now()

	err := r.db.QueryRow(query, listing.Title, listing.Description, listing.ImageURL, listing.Price, listing.AuthorID, listing.CreatedAt).Scan(&listing.ID)
	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}

	return nil
}

func (r *ListingRepository) GetByID(id int64) (*domain.Listing, error) {
	query := `SELECT id, title, description, image_url, price, author_id, created_at FROM listings WHERE id = $1`

	listing := &domain.Listing{}
	err := r.db.QueryRow(query, id).Scan(
		&listing.ID, &listing.Title, &listing.Description, &listing.ImageURL, &listing.Price, &listing.AuthorID, &listing.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found")
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return listing, nil
}

func (r *ListingRepository) GetAll(sortBy, sortOrder string, minPrice, maxPrice *int64) ([]*domain.Listing, error) {
	query, args := r.buildQuery(sortBy, sortOrder, minPrice, maxPrice, 0, 0)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}
	defer rows.Close()

	var listings []*domain.Listing
	for rows.Next() {
		listing := &domain.Listing{}
		err := rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.ImageURL, &listing.Price, &listing.AuthorID, &listing.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		listings = append(listings, listing)
	}

	return listings, nil
}

func (r *ListingRepository) GetAllWithPagination(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int) ([]*domain.Listing, int, error) {
	countQuery, countArgs := r.buildCountQuery(minPrice, maxPrice)
	var total int
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	query, args := r.buildQuery(sortBy, sortOrder, minPrice, maxPrice, page, pageSize)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get listings: %w", err)
	}
	defer rows.Close()

	var listings []*domain.Listing
	for rows.Next() {
		listing := &domain.Listing{}
		err := rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.ImageURL, &listing.Price, &listing.AuthorID, &listing.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan listing: %w", err)
		}
		listings = append(listings, listing)
	}

	return listings, total, nil
}

func (r *ListingRepository) GetByAuthorID(authorID int64) ([]*domain.Listing, error) {
	query := `SELECT id, title, description, image_url, price, author_id, created_at FROM listings WHERE author_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings by author: %w", err)
	}
	defer rows.Close()

	var listings []*domain.Listing
	for rows.Next() {
		listing := &domain.Listing{}
		err := rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.ImageURL, &listing.Price, &listing.AuthorID, &listing.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		listings = append(listings, listing)
	}

	return listings, nil
}

func (r *ListingRepository) buildQuery(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int) (string, []interface{}) {
	query := `SELECT id, title, description, image_url, price, author_id, created_at FROM listings`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if minPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *minPrice)
		argIndex++
	}

	if maxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *maxPrice)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := "created_at"
	switch sortBy {
	case "price":
		orderBy = "price"
	case "title":
		orderBy = "title"
	}

	order := "DESC"
	if sortOrder == "asc" {
		order = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)

	if pageSize > 0 {
		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, pageSize, offset)
	}

	return query, args
}

func (r *ListingRepository) buildCountQuery(minPrice, maxPrice *int64) (string, []interface{}) {
	query := `SELECT COUNT(*) FROM listings`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if minPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *minPrice)
		argIndex++
	}

	if maxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *maxPrice)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}
