package database

import (
	"database/sql"
	"fmt"

	"mini_url/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) DEFAULT 'user',
			created_at TIMESTAMP DEFAULT NOW(),
			last_login TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS links (
			id SERIAL PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_id VARCHAR(20) UNIQUE NOT NULL,
			user_id INTEGER REFERENCES users(id),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP,
			last_accessed_at TIMESTAMP,
			click_count INTEGER DEFAULT 0,
			status VARCHAR(20) DEFAULT 'active',
			expires_at TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_links_short_id ON links(short_id)`,
		`CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}