package coverage_test

import (
	"errors"
	"testing"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/mocks"
	"vk/ecom/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// These tests are for coverage of the main packages
// They test the same functionality but from external package perspective

func TestAuthService_Coverage(t *testing.T) {
	t.Run("should register user successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

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
		assert.NotEqual(t, password, user.Password)

		// Verify password is properly hashed
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should login user successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		login := "testuser"
		password := "password123"

		// Create hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		existingUser := &domain.User{
			ID:       1,
			Login:    login,
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByLogin", login).Return(existingUser, nil)

		token, user, err := authService.LoginUser(login, password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, user)
		assert.Equal(t, existingUser.ID, user.ID)
		assert.Equal(t, existingUser.Login, user.Login)

		mockRepo.AssertExpectations(t)
	})
}
