package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		user, err := h.authService.ValidateToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("user", user)
		c.Set("user_id", int64(user.ID))
		c.Next()
	}
}

func (h *Handler) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Next()
			return
		}

		user, err := h.authService.ValidateToken(authHeader)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("user_id", int64(user.ID))
		c.Next()
	}
}
