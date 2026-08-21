package repository

import (
	"context"
	"database/sql"
	sqlc "gin-backend/internal/db"
	"gin-backend/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, email, password_hash string) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.UserWithPassword, error)
	GetAll(ctx context.Context, limit, offset int64) ([]model.User, error)
	Update(ctx context.Context, req *model.UpdateUserRequest, id int64) (*model.User, error)
	Count(ctx context.Context) (*model.TotalUsers, error)
}

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) Create(ctx context.Context, email, password_hash string) (*model.User, error) {
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{Email: email, PasswordHash: password_hash})
	if err != nil {
		return nil, err
	}
	return &model.User{ID: row.ID, Email: row.Email, CreatedAt: row.CreatedAt}, nil
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
	return &model.User{ID: row.ID, Email: row.Email, CreatedAt: row.CreatedAt}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.UserWithPassword, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &model.UserWithPassword{ID: row.ID, Email: row.Email, PasswordHash: row.PasswordHash, CreatedAt: row.CreatedAt}, nil

}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int64) ([]model.User, error) {
	rows, err := r.q.ListUsers(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	users := make([]model.User, len(rows))
	for i, row := range rows {
		users[i] = model.User{ID: row.ID, Email: row.Email}
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, req *model.UpdateUserRequest, id int64) (*model.User, error) {
	row, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID: id,
		Email: sql.NullString{
			String: deref(req.Email),
			Valid:  req.Email != nil,
		},
		PasswordHash: sql.NullString{
			String: deref(req.PasswordHash),
			Valid:  req.PasswordHash != nil,
		},
	})
	if err != nil {
		return nil, err
	}
	return &model.User{ID: row.ID, Email: row.Email, CreatedAt: row.CreatedAt}, nil
}

func (r *userRepository) Count(ctx context.Context) (*model.TotalUsers, error) {
	count, err := r.q.GetTotalUsers(ctx)
	if err != nil {
		return nil, err
	}
	return &model.TotalUsers{Total: count}, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
