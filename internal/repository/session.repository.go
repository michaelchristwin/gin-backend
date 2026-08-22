package repository

import (
	"context"
	sqlc "gin-backend/internal/db"
	"gin-backend/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) (*model.Session, error)
	GetById(ctx context.Context, id string) (*model.Session, error)
	Delete(ctx context.Context, id string) error
	DeleteUserSessions(ctx context.Context, userID int64) error
	DeleteExpiredSessions(ctx context.Context) error
	GetSessionWithUser(ctx context.Context, id string) (*model.SessionWithUser, error)
}

type sessionRepository struct {
	q *sqlc.Queries
}

func NewSessionRepository(q *sqlc.Queries) SessionRepository {
	return &sessionRepository{q: q}
}

func (s *sessionRepository) Create(ctx context.Context, session *model.Session) (*model.Session, error) {
	row, err := s.q.CreateSession(ctx, sqlc.CreateSessionParams{ID: session.ID, UserID: session.UserID, ExpiresAt: session.ExpiresAt})
	if err != nil {
		return nil, err
	}
	return &model.Session{ID: row.ID, UserID: row.UserID, ExpiresAt: row.ExpiresAt}, nil
}

func (s *sessionRepository) GetById(ctx context.Context, id string) (*model.Session, error) {
	row, err := s.q.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.Session{ID: row.ID, UserID: row.UserID, ExpiresAt: row.ExpiresAt}, nil
}

func (s *sessionRepository) Delete(ctx context.Context, id string) error {
	err := s.q.DeleteSession(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionRepository) DeleteUserSessions(ctx context.Context, userID int64) error {
	err := s.q.DeleteUserSessions(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	err := s.q.DeleteExpiredSessions(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionRepository) GetSessionWithUser(ctx context.Context, id string) (*model.SessionWithUser, error) {
	row, err := s.q.GetSessionWithUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.SessionWithUser{
		SessionID:    row.SessionID,
		UserID:       row.UserID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		ExpiresAt:    row.ExpiresAt}, nil
}
