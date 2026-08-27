package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserWithPassword struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpdateUserRequest struct {
	Email *string `json:"email" binding:"required,email"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=12,max=64"`
}

type TotalUsers struct {
	Total int64 `json:"total"`
}
