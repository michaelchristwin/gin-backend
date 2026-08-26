package service

import "gin-backend/internal/repository"

type SessionService interface {
	GetByID()
}

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) SessionService {
	return &sessionService{sessionRepo: sessionRepo}
}

func (s *sessionService) GetByID() {}
