package service

import (
	"errors"
	"testing"
	"vk/ecom/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_RegisterUser_InPackage(t *testing.T) {
	t.Run("should successfully register new user", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := NewAuthService(mockRepo)

		login := "testuser"
		password := "password123"

		// Mock that user doesn't exist
		mockRepo.On("GetByLogin", login).Return(nil, errors.New("user not found"))

		// Mock successful user creation
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

		user, err := authService.RegisterUser(login, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, login, user.Login)
		assert.NotEmpty(t, user.Password)
		assert.NotEqual(t, password, user.Password) // Password should be hashed

		// Verify password is properly hashed
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with short login", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := NewAuthService(mockRepo)

		user, err := authService.RegisterUser("ab", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login must be between 3 and 20 characters")
	})
}
