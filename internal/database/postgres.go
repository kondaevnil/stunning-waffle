package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}

func MigrateDB(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_login ON users(login)`,

		`CREATE TABLE IF NOT EXISTS listings (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			image_url VARCHAR(500),
			price BIGINT NOT NULL,
			author_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_listings_author_id ON listings(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_listings_price ON listings(price)`,
		`CREATE INDEX IF NOT EXISTS idx_listings_created_at ON listings(created_at)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}
