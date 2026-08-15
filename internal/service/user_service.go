package service

import (
	"gin-backend/internal/model"
	"gin-backend/internal/repository"

	"context"
)

type UserService interface {
	RegisterUser(ctx context.Context, name, email string) (*model.User, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int64) ([]model.User, error)
	DeleteUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, name, email string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, name, email string) (*model.User, error) {

	return s.repo.Create(ctx, name, email)
}

func (s *userService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetById(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int64) ([]model.User, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.Delete(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, name, email string) (*model.User, error) {
	return s.repo.Update(ctx, name, email)
}
