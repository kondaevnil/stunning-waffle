package memory

import (
	"errors"
	"sort"
	"sync"
	"time"
	"vk/ecom/internal/domain"
)

type InMemoryListingRepository struct {
	listings map[int64]*domain.Listing
	nextID   int64
	mu       sync.RWMutex
}

func NewInMemoryListingRepository() *InMemoryListingRepository {
	return &InMemoryListingRepository{
		listings: make(map[int64]*domain.Listing),
		nextID:   1,
	}
}

func (r *InMemoryListingRepository) Create(listing *domain.Listing) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	listing.ID = r.nextID
	listing.CreatedAt = time.Now()
	r.listings[r.nextID] = listing
	r.nextID++
	return nil
}

func (r *InMemoryListingRepository) GetByID(id int64) (*domain.Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if listing, exists := r.listings[id]; exists {
		return listing, nil
	}
	return nil, errors.New("listing not found")
}

func (r *InMemoryListingRepository) GetAll(sortBy, sortOrder string, minPrice, maxPrice *int64) ([]*domain.Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Listing
	for _, listing := range r.listings {
		if minPrice != nil && listing.Price < *minPrice {
			continue
		}
		if maxPrice != nil && listing.Price > *maxPrice {
			continue
		}
		result = append(result, listing)
	}

	switch sortBy {
	case "price":
		sort.Slice(result, func(i, j int) bool {
			if sortOrder == "desc" {
				return result[i].Price > result[j].Price
			}
			return result[i].Price < result[j].Price
		})
	case "date":
		sort.Slice(result, func(i, j int) bool {
			if sortOrder == "desc" {
				return result[i].CreatedAt.After(result[j].CreatedAt)
			}
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		})
	}

	return result, nil
}

func (r *InMemoryListingRepository) GetAllWithPagination(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int) ([]*domain.Listing, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filteredListings []*domain.Listing
	for _, listing := range r.listings {
		if minPrice != nil && listing.Price < *minPrice {
			continue
		}
		if maxPrice != nil && listing.Price > *maxPrice {
			continue
		}
		filteredListings = append(filteredListings, listing)
	}

	switch sortBy {
	case "price":
		sort.Slice(filteredListings, func(i, j int) bool {
			if sortOrder == "desc" {
				return filteredListings[i].Price > filteredListings[j].Price
			}
			return filteredListings[i].Price < filteredListings[j].Price
		})
	case "date":
		sort.Slice(filteredListings, func(i, j int) bool {
			if sortOrder == "desc" {
				return filteredListings[i].CreatedAt.After(filteredListings[j].CreatedAt)
			}
			return filteredListings[i].CreatedAt.Before(filteredListings[j].CreatedAt)
		})
	default:
		sort.Slice(filteredListings, func(i, j int) bool {
			return filteredListings[i].CreatedAt.After(filteredListings[j].CreatedAt)
		})
	}

	totalCount := len(filteredListings)

	start := (page - 1) * pageSize
	if start < 0 {
		start = 0
	}
	if start >= totalCount {
		return []*domain.Listing{}, totalCount, nil
	}

	end := start + pageSize
	if end > totalCount {
		end = totalCount
	}

	return filteredListings[start:end], totalCount, nil
}

func (r *InMemoryListingRepository) GetByAuthorID(authorID int64) ([]*domain.Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var authorListings []*domain.Listing
	for _, listing := range r.listings {
		if listing.AuthorID == authorID {
			authorListings = append(authorListings, listing)
		}
	}
	if len(authorListings) == 0 {
		return nil, errors.New("no listings found for this author")
	}
	return authorListings, nil
}
