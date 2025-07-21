package memory

import (
	"errors"
	"sync"
	"vk/ecom/internal/domain"
)

type InMemoryUserRepository struct {
	users  map[int]*domain.User
	nextID int
	mu     sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++
	return nil
}

func (r *InMemoryUserRepository) GetByID(id int) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByLogin(login string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Login == login {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
