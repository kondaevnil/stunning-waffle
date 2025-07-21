package handler

import (
	"vk/ecom/internal/interfaces"
)

type Handler struct {
	authService    interfaces.AuthServiceInterface
	listingService interfaces.ListingServiceInterface
}

func NewHandler(authService interfaces.AuthServiceInterface, listingService interfaces.ListingServiceInterface) *Handler {
	return &Handler{
		authService:    authService,
		listingService: listingService,
	}
}
