package auth

import "time"

type Auth struct {
	ID           uint      `json:"id"`
	UserID       string    `json:"user_id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthLog struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Login     string    `json:"login"`
	Action    string    `json:"action"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
