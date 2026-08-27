package service

import (
	"gin-backend/internal/model"
	"gin-backend/internal/repository"

	"context"
)

type UserService interface {
	GetUserById(ctx context.Context, id int64) (*model.UserWithPassword, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserWithPassword, error)
	ListUsers(ctx context.Context, limit, offset int64) ([]model.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, req *model.UpdateUserRequest, id int64) (*model.User, error)
	GetTotalUsers(ctx context.Context) (*model.TotalUsers, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserById(ctx context.Context, id int64) (*model.UserWithPassword, error) {
	return s.repo.GetById(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.UserWithPassword, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int64) ([]model.User, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)

}

func (s *userService) UpdateUser(ctx context.Context, req *model.UpdateUserRequest, id int64) (*model.User, error) {
	return s.repo.UpdateUser(ctx, req, id)
}

func (s *userService) GetTotalUsers(ctx context.Context) (*model.TotalUsers, error) {
	return s.repo.Count(ctx)
}
