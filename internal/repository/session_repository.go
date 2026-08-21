package repository

import (
	"context"
	sqlc "gin-backend/internal/db"
)

type SessionRepository interface {
	Create(ctx context.Context)
	Get()
	Delete()
	DeleteUserSessions()
	DeleteExpiredSessions()
	GetSessionWithUser()
}

type sessionRepository struct {
	q *sqlc.Queries
}

// func NewSessionRepository(q *sqlc.Queries) SessionRepository {
// 	return &sessionRepository{q: q}
// }
