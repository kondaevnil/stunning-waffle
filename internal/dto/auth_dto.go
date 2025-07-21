package dto

import (
	"vk/ecom/internal/domain"
)

type UserDTO struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

func ToUserDTO(user *domain.User) *UserDTO {
	if user == nil {
		return nil
	}
	return &UserDTO{
		ID:    user.ID,
		Login: user.Login,
	}
}
