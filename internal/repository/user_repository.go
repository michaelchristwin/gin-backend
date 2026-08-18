package repository

import (
	"context"
	sqlc "gin-backend/internal/db"
	"gin-backend/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, name, email string) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*model.User, error)
	GetAll(ctx context.Context, limit, offset int64) ([]model.User, error)
	Update(ctx context.Context, req *model.UpdateUserRequest, id int64) (*model.User, error)
}

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) Create(ctx context.Context, name, email string) (*model.User, error) {
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{Name: name, Email: email})
	if err != nil {
		return nil, err
	}
	return &model.User{ID: row.ID, Name: row.Name, Email: row.Email}, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	row, err := r.q.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: row.ID, Name: row.Name, Email: row.Email}, nil
}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int64) ([]model.User, error) {
	rows, err := r.q.ListUsers(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	users := make([]model.User, len(rows))
	for i, row := range rows {
		users[i] = model.User{ID: row.ID, Name: row.Name, Email: row.Email}
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, req *model.UpdateUserRequest, id int64) (*model.User, error) {
	row, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{ID: id, Name: *req.Name, Email: *req.Email})
	if err != nil {
		return nil, err
	}
	return &model.User{ID: row.ID, Name: row.Name, Email: row.Email}, nil
}
