package service

import (
	"context"
	"gin-backend/internal/model"
	"gin-backend/internal/repository"
)

type SessionService interface {
	GetByID(ctx context.Context, id string) (*model.Session, error)
}

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) SessionService {
	return &sessionService{sessionRepo: sessionRepo}
}

func (s *sessionService) GetByID(ctx context.Context, id string) (*model.Session, error) {
	return s.sessionRepo.GetById(ctx, id)

}
