package models

import "time"

type Link struct {
	ID             int       `json:"id"`
	OriginalURL    string    `json:"original_url"`
	ShortID        string    `json:"short_id"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	ClickCount     int       `json:"click_count"`
	Status         string    `json:"status"` // "active", "deleted", "blocked"
	ExpiresAt      time.Time `json:"expires_at"`
}