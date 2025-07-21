package postgres

import (
	"database/sql"
	"fmt"
	"vk/ecom/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id`

	err := r.db.QueryRow(query, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	query := `SELECT id, login, password FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByLogin(login string) (*domain.User, error) {
	query := `SELECT id, login, password FROM users WHERE login = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
