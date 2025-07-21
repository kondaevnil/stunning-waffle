package mocks

import (
	"vk/ecom/internal/domain"
	"vk/ecom/internal/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

var _ interfaces.AuthServiceInterface = (*MockAuthService)(nil)

func (m *MockAuthService) RegisterUser(login, password string) (*domain.User, error) {
	args := m.Called(login, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) LoginUser(login, password string) (string, *domain.User, error) {
	args := m.Called(login, password)
	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*domain.User), args.Error(2)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*domain.User, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
