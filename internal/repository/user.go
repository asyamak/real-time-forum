package repository

import (
	"database/sql"
	"real-time-forum/internal/model"
)

type User interface {
	Create(user model.User) error
	GetByCredentials(usernameOrEmail, password string) (model.User, error)
	GetByID(userID int) (model.User, error)
	GetUsersPosts(userID int) ([]model.Post, error)
	GetUsersRatedPosts(userID int) ([]model.Post, error)
	SetSession(session model.Session) error
	DeleteSession(userID int, refreshToken string) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user model.User) error {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetByCredentials(usernameOrEmail string, password string) (model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetByID(userID int) (model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetUsersPosts(userID int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) GetUsersRatedPosts(userID int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) SetSession(session model.Session) error {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) DeleteSession(userID int, refreshToken string) error {
	panic("not implemented") // TODO: Implement
}
