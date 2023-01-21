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
