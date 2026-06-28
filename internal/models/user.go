package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Не возвращается в JSON
	Role      string    `json:"role"` // "user" или "admin"
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
}