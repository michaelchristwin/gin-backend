package model

import "time"

type Session struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionReq struct {
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionWithUser struct {
	SessionID    string    `json:"session_id"`
	UserID       int64     `json:"user_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
