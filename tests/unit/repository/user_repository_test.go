package repository_test

import (
	"database/sql"
	"testing"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/repository/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	t.Run("should successfully create user", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		user := &domain.User{
			Login:    "testuser",
			Password: "hashedpassword",
		}

		expectedID := 1
		mock.ExpectQuery(`INSERT INTO users \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs(user.Login, user.Password).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

		err = repo.Create(user)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, user.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle database error during create", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		user := &domain.User{
			Login:    "testuser",
			Password: "hashedpassword",
		}

		mock.ExpectQuery(`INSERT INTO users \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs(user.Login, user.Password).
			WillReturnError(sql.ErrConnDone)

		err = repo.Create(user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Run("should successfully get user by ID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		expectedUser := &domain.User{
			ID:       1,
			Login:    "testuser",
			Password: "hashedpassword",
		}

		mock.ExpectQuery(`SELECT id, login, password FROM users WHERE id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).
				AddRow(expectedUser.ID, expectedUser.Login, expectedUser.Password))

		user, err := repo.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Login, user.Login)
		assert.Equal(t, expectedUser.Password, user.Password)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		mock.ExpectQuery(`SELECT id, login, password FROM users WHERE id = \$1`).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle database error during get", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		mock.ExpectQuery(`SELECT id, login, password FROM users WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		user, err := repo.GetByID(1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByLogin(t *testing.T) {
	t.Run("should successfully get user by login", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		expectedUser := &domain.User{
			ID:       1,
			Login:    "testuser",
			Password: "hashedpassword",
		}

		mock.ExpectQuery(`SELECT id, login, password FROM users WHERE login = \$1`).
			WithArgs("testuser").
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).
				AddRow(expectedUser.ID, expectedUser.Login, expectedUser.Password))

		user, err := repo.GetByLogin("testuser")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Login, user.Login)
		assert.Equal(t, expectedUser.Password, user.Password)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user not found by login", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		mock.ExpectQuery(`SELECT id, login, password FROM users WHERE login = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByLogin("nonexistent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle database error during get by login", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := postgres.NewUserRepository(db)

		mock.ExpectQuery(`SELECT id, login, password FROM users WHERE login = \$1`).
			WithArgs("testuser").
			WillReturnError(sql.ErrConnDone)

		user, err := repo.GetByLogin("testuser")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
