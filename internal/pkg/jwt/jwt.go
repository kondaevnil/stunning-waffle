package jwt

import (
	"vk/ecom/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTClaim struct {
	UserID int `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("your_secret_key")

func GenerateToken(user *domain.User) (string, error) {
	claims := &JWTClaim{
		UserID: user.ID,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ecom",
			Subject:   "user_token",
			Audience:  []string{"ecom_users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	return &domain.User{
		ID:    claims.UserID,
		Login: claims.Login,
	}, nil
}