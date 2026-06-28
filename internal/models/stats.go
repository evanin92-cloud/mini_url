package models

import "time"

type Stats struct {
	ID         int       `json:"id"`
	LinkID     int       `json:"link_id"`
	AccessedAt time.Time `json:"accessed_at"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Referer    string    `json:"referer"`
}