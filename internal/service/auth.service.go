package service

import (
	"context"
	"errors"
	"fmt"
	"gin-backend/internal/domain"
	"gin-backend/internal/model"
	"gin-backend/internal/repository"

	"gin-backend/internal/util"
	"gin-backend/internal/validation"
	"log/slog"
	"time"

	"github.com/alexedwards/argon2id"
)

type AuthService interface {
	RegisterUser(ctx context.Context, userReq model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, user model.RegisterRequest) (*model.Session, error)
	Logout(ctx context.Context, sessionId string) error
	ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error
	checkEmailAvailable(ctx context.Context, email string) error
}

type authService struct {
	repo     repository.UserRepository
	sessions repository.SessionRepository
	logger   *slog.Logger
}

func NewAuthService(repo repository.UserRepository,
	sessions repository.SessionRepository, logger *slog.Logger) AuthService {
	return &authService{repo: repo, sessions: sessions, logger: logger}
}

func (s *authService) RegisterUser(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	if err := s.checkEmailAvailable(ctx, req.Email); err != nil {
		return nil, err // could be ErrEmailTaken or a wrapped infra error
	}
	passwordHash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {

		return nil, domain.ErrInvalidCredentials
	}
	if err := validation.Validate(ctx, req.Password, s.logger); err != nil {
		return nil, err // e.g. ErrTooShort, ErrPwned
	}
	user, err := s.repo.Create(ctx, req.Email, passwordHash)
	if err != nil {
		if validation.IsUniqueViolation(err) { // check pg error code 23505, or your driver's equivalent
			return nil, domain.ErrEmailTaken
		}
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, userReq model.RegisterRequest) (*model.Session, error) {
	user, err := s.repo.GetByEmail(ctx, userReq.Email)
	if err != nil {
		s.logger.Warn("authentication failed",
			"reason", "invalid_credentials",
		)
		return nil, domain.ErrInvalidCredentials
	}
	match, err := argon2id.ComparePasswordAndHash(userReq.Password, user.PasswordHash)
	if err != nil {
		s.logger.Warn("authentication failed",
			"user_id", user.ID,
			"reason", "invalid_credentials",
		)
		return nil, domain.ErrInvalidCredentials
	}
	if !match {
		s.logger.Warn("authentication failed",
			"user_id", user.ID,
			"reason", "invalid_credentials",
		)
		return nil, domain.ErrInvalidCredentials
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
		s.logger.Warn("authentication failed",
			"reason", "invalid_credentials",
		)
		return err
	}
	match, err := argon2id.ComparePasswordAndHash(currentPassword, user.PasswordHash)
	if err != nil {
		s.logger.Warn("authentication failed",
			"user_id", user.ID,
			"reason", "invalid_credentials",
		)
		return err
	}
	if !match {
		s.logger.Warn("authentication failed",
			"user_id", user.ID,
			"reason", "invalid_credentials",
		)
		return domain.ErrInvalidCredentials
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

func (s *authService) checkEmailAvailable(ctx context.Context, email string) error {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return domain.ErrEmailTaken // found a user, email is taken
	}
	if errors.Is(err, domain.ErrUserNotFound) { // or sql.ErrNoRows, whatever your repo returns
		return nil // expected — email is available
	}
	return fmt.Errorf("checking email availability: %w", err) // real error — propagate it
}
