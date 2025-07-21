package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"vk/ecom/internal/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateListing(c *gin.Context) {
	var req dto.ListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	listing, err := h.listingService.CreateListing(&req, c.GetInt64("user_id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"listing": listing,
	})
}

func (h *Handler) GetListings(c *gin.Context) {
	sortBy := c.DefaultQuery("sort", "date")
	sortOrder := c.DefaultQuery("order", "desc")

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil && val > 0 {
			page = val
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if val, err := strconv.Atoi(pageSizeStr); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	}

	var minPrice, maxPrice *int64
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if val, err := strconv.ParseInt(minPriceStr, 10, 64); err == nil && val >= 0 {
			minPrice = &val
		}
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if val, err := strconv.ParseInt(maxPriceStr, 10, 64); err == nil && val >= 0 {
			maxPrice = &val
		}
	}

	if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price cannot be greater than max_price"})
		return
	}

	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int64); ok {
			currentUserID = &id
		}
	}

	fmt.Println("Current User ID:", currentUserID)

	response, err := h.listingService.GetListingsWithPagination(sortBy, sortOrder, minPrice, maxPrice, page, pageSize, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve listings"})
		return
	}

	c.JSON(http.StatusOK, response)
}
