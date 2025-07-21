package service

import (
	"errors"
	"fmt"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/dto"
	"vk/ecom/internal/interfaces"
	"vk/ecom/internal/repository"
)

type ListingService struct {
	listingRepo repository.ListingRepository
	userRepo    repository.UserRepository
}

// Ensure ListingService implements ListingServiceInterface
var _ interfaces.ListingServiceInterface = (*ListingService)(nil)

func NewListingService(listingRepo repository.ListingRepository, userRepo repository.UserRepository) *ListingService {
	return &ListingService{
		listingRepo: listingRepo,
		userRepo:    userRepo,
	}
}

func (s *ListingService) CreateListing(req *dto.ListingRequest, authorID int64) (*dto.ListingDTO, error) {
	const (
		minTitleLen       = 3
		maxTitleLen       = 100
		minDescLen        = 10
		maxDescLen        = 2000
		minPrice    int64 = 1
		maxPrice    int64 = 1_000_000_000
	)
	allowedImageFormats := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if len(req.Title) < minTitleLen || len(req.Title) > maxTitleLen {
		errString := "title must be between %d and %d characters"
		return nil, fmt.Errorf(errString, minTitleLen, maxTitleLen)
	}
	if len(req.Description) < minDescLen || len(req.Description) > maxDescLen {
		errString := "description must be between %d and %d characters"
		return nil, fmt.Errorf(errString, minDescLen, maxDescLen)
	}
	if req.Price < minPrice || req.Price > maxPrice {
		errString := "price must be between %d and %d"
		return nil, fmt.Errorf(errString, minPrice, maxPrice)
	}
	if req.ImageURL != "" {
		ext := ""
		if dot := len(req.ImageURL) - 4; dot >= 0 {
			ext = req.ImageURL[dot:]
		}
		if !allowedImageFormats[ext] {
			return nil, errors.New("unsupported image format")
		}
	}

	listing := &domain.Listing{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		AuthorID:    authorID,
	}

	err := s.listingRepo.Create(listing)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(int(authorID))
	if err != nil {
		return dto.ToListingDTOWithAuthor(listing, "", &authorID), nil
	}

	return dto.ToListingDTOWithAuthor(listing, author.Login, &authorID), nil
}

func (s *ListingService) GetListings(sortBy, sortOrder string, minPrice, maxPrice *int64, currentUserID *int64) ([]*dto.ListingDTO, error) {
	listings, err := s.listingRepo.GetAll(sortBy, sortOrder, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}

	var result []*dto.ListingDTO
	for _, listing := range listings {
		author, err := s.userRepo.GetByID(int(listing.AuthorID))
		if err != nil {
			result = append(result, dto.ToListingDTOWithAuthor(listing, "", currentUserID))
		} else {
			result = append(result, dto.ToListingDTOWithAuthor(listing, author.Login, currentUserID))
		}
	}

	return result, nil
}

func (s *ListingService) GetListingsWithPagination(sortBy, sortOrder string, minPrice, maxPrice *int64, page, pageSize int, currentUserID *int64) (*dto.ListingsResponse, error) {
	const (
		defaultPageSize = 10
		maxPageSize     = 100
		minPageSize     = 1
	)

	if page < 1 {
		page = 1
	}
	if pageSize < minPageSize {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	if sortBy != "price" && sortBy != "date" {
		sortBy = "date"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	listings, totalCount, err := s.listingRepo.GetAllWithPagination(sortBy, sortOrder, minPrice, maxPrice, page, pageSize)
	if err != nil {
		return nil, err
	}

	var result []*dto.ListingDTO
	for _, listing := range listings {
		author, err := s.userRepo.GetByID(int(listing.AuthorID))
		if err != nil {
			result = append(result, dto.ToListingDTOWithAuthor(listing, "", currentUserID))
		} else {
			result = append(result, dto.ToListingDTOWithAuthor(listing, author.Login, currentUserID))
		}
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.ListingsResponse{
		Listings:   result,
		Count:      len(result),
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
