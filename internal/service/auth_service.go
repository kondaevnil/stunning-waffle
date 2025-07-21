package service

import (
	"errors"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/interfaces"
	"vk/ecom/internal/pkg/jwt"
	"vk/ecom/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

var _ interfaces.AuthServiceInterface = (*AuthService)(nil)

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) RegisterUser(login, password string) (*domain.User, error) {
	if len(login) < 3 || len(login) > 20 {
		return nil, errors.New("login must be between 3 and 20 characters")
	}

	if len(password) < 6 || len(password) > 40 {
		return nil, errors.New("password must be between 6 and 40 characters")
	}

	_, err := s.userRepo.GetByLogin(login)
	if err == nil {
		return nil, errors.New("user with this login already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Login:    login,
		Password: string(hashedPassword),
	}

	err = s.userRepo.Create(user)

	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *AuthService) LoginUser(login, password string) (string, *domain.User, error) {
	user, err := s.userRepo.GetByLogin(login)

	if err != nil {
		return "", nil, errors.New("invalid login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", nil, errors.New("invalid login or password")
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}
	return token, user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*domain.User, error) {
	user, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	_, err = s.userRepo.GetByID(user.ID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
