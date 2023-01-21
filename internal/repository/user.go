package repository

import (
	"context"
	"database/sql"
	"fmt"
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
	stmt, err := r.db.PrepareContext(ctx, `
	INSERT INTO 
		user
			(email, username, password, first_name, last_name, age, gender, avatar, creation_time)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9);`)
	if err != nil {
		return fmt.Errorf("repo: create user: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		user.Email,
		user.Username,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Age,
		user.Gender,
		user.Avatar,
		user.CreationTime,
	)

	if isAlreadyExists(err) {
		return ErrUserExists
	}

	return nil
}

func (r *UserRepository) GetByCredentials(ctx context.Context, usernameOrEmail string, password string) (model.User, error) {
	var user model.User

	err := r.db.QueryRowContext(ctx, `
		SELECT 
			id
		FROM 
			users
		WHERE 
			(username = $1 OR email = $1)
		AND
			(password = $2);
		`, usernameOrEmail, password).Scan(&user.ID)

	if isNoRowsError(err) {
		return model.User{}, ErrNoRows
	}

	return user, nil
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
