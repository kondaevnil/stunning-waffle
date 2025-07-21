package handler

import (
	"fmt"
	"net/http"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Printf("User data: %v\n", user)

	if user, err := h.authService.RegisterUser(user.Login, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"user": dto.ToUserDTO(user)})
	}
}
func (h *Handler) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, u, err := h.authService.LoginUser(user.Login, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": u.ID, "login": u.Login})
}
