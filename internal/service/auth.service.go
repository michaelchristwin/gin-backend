package service

import (
	"context"
	"errors"
	"gin-backend/internal/model"
	"gin-backend/internal/repository"
	"gin-backend/internal/util"
	"gin-backend/internal/validation"
	"time"

	"github.com/alexedwards/argon2id"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordTooWeak    = errors.New("password does not meet requirements")
	// ... other auth-specific sentinel errors
)

type AuthService interface {
	RegisterUser(ctx context.Context, userReq model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, user model.RegisterRequest) (*model.Session, error)
	Logout(ctx context.Context, sessionId string) error
	ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error
}

type authService struct {
	repo     repository.UserRepository
	sessions repository.SessionRepository
}

func NewAuthService(repo repository.UserRepository,
	sessions repository.SessionRepository) AuthService {
	return &authService{repo: repo, sessions: sessions}
}

func (s *authService) RegisterUser(ctx context.Context, user model.RegisterRequest) (*model.User, error) {
	passwordHash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := validation.Validate(ctx, user.Password); err != nil {
		return nil, err // e.g. ErrTooShort, ErrPwned
	}
	return s.repo.Create(ctx, user.Email, passwordHash)
}

func (s *authService) Login(ctx context.Context, userReq model.RegisterRequest) (*model.Session, error) {
	user, err := s.repo.GetByEmail(ctx, userReq.Email)
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
	return s.sessions.Create(ctx, &model.Session{ID: sessionID, UserID: user.ID, ExpiresAt: expiresAt})
}

func (s *authService) Logout(ctx context.Context, sessionId string) error {
	err := s.sessions.Delete(ctx, sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	user, err := s.repo.GetById(ctx, userID)
	if err != nil {
		return err
	}
	match, err := argon2id.ComparePasswordAndHash(currentPassword, user.PasswordHash)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidCredentials
	}
	hash, err := argon2id.CreateHash(newPassword, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePasswordHash(ctx, userID, hash); err != nil {
		return err
	}
	return nil
}
