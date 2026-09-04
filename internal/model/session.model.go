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
	UserID    int64     `json:"user_id"`
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteSessionReq struct {
	SessionID string `json:"session_id"`
}
