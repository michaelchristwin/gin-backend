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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordTooWeak    = errors.New("password does not meet requirements")
	// ... other auth-specific sentinel errors
)

type AuthService interface {
	Login(ctx context.Context, userReq model.RegisterRequest) (*model.Session, error)
	Logout(ctx context.Context, sessionId string) error
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository) AuthService {
	return &authService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *authService) Login(ctx context.Context, userReq model.RegisterRequest) (*model.Session, error) {
	user, err := s.userRepo.GetByEmail(ctx, userReq.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	match, err := argon2id.ComparePasswordAndHash(userReq.Password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrInvalidCredentials
	}
	sessionID, err := util.GenerateSessionID()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	return s.sessionRepo.Create(ctx, &model.Session{ID: sessionID, UserID: user.ID, ExpiresAt: expiresAt})
}

func (s *authService) Logout(ctx context.Context, sessionId string) error {
	err := s.sessionRepo.Delete(ctx, sessionId)
	if err != nil {
		return err
	}
	return nil
}
