package repository

import (
	"context"
	"database/sql"
	"real-time-forum/internal/model"
)

type User interface {
	Create(ctx context.Context, user model.User) error
	GetByCredentials(ctx context.Context, usernameOrEmail, password string) (model.User, error)
	GetByID(ctx context.Context, userID int) (model.User, error)
	GetUsersPosts(ctx context.Context, userID int) ([]model.Post, error)
	GetUsersRatedPosts(ctx context.Context, userID int) ([]model.Post, error)
	SetSession(ctx context.Context, session model.Session) error
	DeleteSession(ctx context.Context, userID int, refreshToken string) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetByCredentials(ctx context.Context, usernameOrEmail string, password string) (model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetByID(ctx context.Context, userID int) (model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetUsersPosts(ctx context.Context, userID int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetUsersRatedPosts(ctx context.Context, userID int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) SetSession(ctx context.Context, session model.Session) error {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) DeleteSession(ctx context.Context, userID int, refreshToken string) error {
	panic("not implemented") // TODO: Implement
}
