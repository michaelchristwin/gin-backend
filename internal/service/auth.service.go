package service

import (
	"context"
	"errors"
	"gin-backend/internal/model"
	"gin-backend/internal/repository"
	"gin-backend/internal/util"
	"time"

	"github.com/alexedwards/argon2id"
)

type AuthService interface {
	Login(ctx context.Context, userReq model.CreateUserRequest) (*model.Session, error)
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository) AuthService {
	return &authService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *authService) Login(ctx context.Context, userReq model.CreateUserRequest) (*model.Session, error) {
	user, err := s.userRepo.GetByEmail(ctx, userReq.Email)
	if err != nil {
		return nil, errors.New("Invalid credentials")
	}
	match, err := argon2id.ComparePasswordAndHash(userReq.Password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("Invalid credentials")
	}
	sessionID, err := util.GenerateSessionID()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	return s.sessionRepo.Create(ctx, &model.Session{ID: sessionID, UserID: user.ID, ExpiresAt: expiresAt})
}
