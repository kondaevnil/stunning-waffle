package service_test

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

func TestAuthService_RegisterUser(t *testing.T) {
	t.Run("should successfully register new user", func(t *testing.T) {
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
		assert.NotEqual(t, password, user.Password) // Password should be hashed

		// Verify password is properly hashed
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with short login", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		user, err := authService.RegisterUser("ab", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login must be between 3 and 20 characters")
	})

	t.Run("should fail with long login", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		longLogin := "this_is_a_very_long_login_name_that_exceeds_limit"
		user, err := authService.RegisterUser(longLogin, "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login must be between 3 and 20 characters")
	})

	t.Run("should fail with short password", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		user, err := authService.RegisterUser("testuser", "123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "password must be between 6 and 40 characters")
	})

	t.Run("should fail with long password", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		longPassword := "this_is_a_very_long_password_that_definitely_exceeds_the_forty_character_limit_set_by_validation"
		user, err := authService.RegisterUser("testuser", longPassword)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "password must be between 6 and 40 characters")
	})

	t.Run("should fail when user already exists", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		login := "testuser"
		existingUser := &domain.User{ID: 1, Login: login, Password: "hashedpass"}

		// Mock that user already exists
		mockRepo.On("GetByLogin", login).Return(existingUser, nil)

		user, err := authService.RegisterUser(login, "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user with this login already exists")

		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when repository create fails", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		login := "testuser"
		password := "password123"

		// Mock that user doesn't exist
		mockRepo.On("GetByLogin", login).Return(nil, errors.New("user not found"))

		// Mock failed user creation
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("database error"))

		user, err := authService.RegisterUser(login, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to create user")

		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_LoginUser(t *testing.T) {
	t.Run("should successfully login with valid credentials", func(t *testing.T) {
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

	t.Run("should fail with non-existent user", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		mockRepo.On("GetByLogin", "nonexistent").Return(nil, errors.New("user not found"))

		token, user, err := authService.LoginUser("nonexistent", "password")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid login or password")

		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		login := "testuser"
		correctPassword := "password123"
		wrongPassword := "wrongpassword"

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
		existingUser := &domain.User{
			ID:       1,
			Login:    login,
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByLogin", login).Return(existingUser, nil)

		token, user, err := authService.LoginUser(login, wrongPassword)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid login or password")

		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	t.Run("should successfully validate valid token", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		// First create a user and generate token
		user := &domain.User{ID: 1, Login: "testuser"}
		mockRepo.On("GetByID", 1).Return(user, nil)

		// We need to generate a real token for this test
		// Let's create another auth service to generate token
		mockRepo2 := new(mocks.MockUserRepository)
		authService2 := service.NewAuthService(mockRepo2)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		loginUser := &domain.User{ID: 1, Login: "testuser", Password: string(hashedPassword)}

		mockRepo2.On("GetByLogin", "testuser").Return(loginUser, nil)
		token, _, _ := authService2.LoginUser("testuser", "password")

		validatedUser, err := authService.ValidateToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, validatedUser)
		assert.Equal(t, user.ID, validatedUser.ID)
		assert.Equal(t, user.Login, validatedUser.Login)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		invalidToken := "invalid.token.here"

		user, err := authService.ValidateToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("should fail when user not found in repository", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		authService := service.NewAuthService(mockRepo)

		// Generate a valid token but user doesn't exist in repo
		mockRepo2 := new(mocks.MockUserRepository)
		authService2 := service.NewAuthService(mockRepo2)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		loginUser := &domain.User{ID: 1, Login: "testuser", Password: string(hashedPassword)}

		mockRepo2.On("GetByLogin", "testuser").Return(loginUser, nil)
		token, _, _ := authService2.LoginUser("testuser", "password")

		// Mock that user is not found
		mockRepo.On("GetByID", 1).Return(nil, errors.New("user not found"))

		user, err := authService.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")

		mockRepo.AssertExpectations(t)
	})
}
