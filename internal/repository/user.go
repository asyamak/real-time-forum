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
	GetUsersVotedPosts(ctx context.Context, userID int) ([]model.Post, error)
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

	_, err = stmt.ExecContext(
		ctx,
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

	return err
}

func (r *UserRepository) GetByCredentials(ctx context.Context, usernameOrEmail string, password string) (model.User, error) {
	var user model.User

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT 
			id, email, username, password, first_name, last_name, age, gender, avatar, creation_time
		FROM 
			user
		WHERE 
			(username = $1 OR email = $1)
		AND
			(password = $2);`)
	if err != nil {
		return model.User{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, usernameOrEmail, password)
	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Gender,
		&user.Avatar,
		&user.CreationTime,
	)

	if isNoRowsError(err) {
		return model.User{}, ErrNoRows
	}

	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, userID int) (model.User, error) {
	var user model.User

	stmt, err := r.db.PrepareContext(ctx, `
	SELECT
		id, email, username, password, first_name, last_name, age, gender, avatar, creation_time
	FROM
		user
	WHERE
		id = $1`)
	if err != nil {
		return model.User{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userID)
	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Gender,
		&user.Avatar,
		&user.CreationTime,
	)

	if isNoRowsError(err) {
		return model.User{}, ErrNoRows
	}

	return model.User{}, err
}

// GetUsersPosts method receives all posts created by userId
func (r *UserRepository) GetUsersPosts(ctx context.Context, userID int) ([]model.Post, error) {
	var isUserExists bool

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("get users posts: transaction: begin: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `SELECT EXISTS (SELECT id FROM user WHERE id = $1);`)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("get users posts: prepare statement 1: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userID)

	row.Scan(&isUserExists)
	if err != nil {
		return nil, fmt.Errorf("get users posts: row scan: %w", err)
	}

	if !isUserExists {
		return nil, fmt.Errorf("get user posts: user is not exists")
	}

	var posts []model.Post

	stmt, err = tx.PrepareContext(ctx, `
	SELECT 
		p.id, 
		p.title, 
		p.content, 
		p.creation_time, 
		p.image, 
		u.id AS author.id,
		u.username AS author.username, 
		u.first_name AS author.first_name, 
		u.last_name AS author.last_name, 
		u.avatar AS author.avatar
	FROM post p LEFT JOIN user u ON p.user_id = u.id
	WHERE post.user_id = $1
	ORDER BY 1 DESC;`)

	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("get users posts: prepare statement 2: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get users posts: exec statemnt: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreationTime,
			&post.ImagePath,
			&post.Author.ID,
			&post.Author.Username,
			&post.Author.FirstName,
			&post.Author.LastName,
			&post.Author.Avatar,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("get users posts: transaction: commit %w", err)
	}

	return posts, nil
}

// func (r *UserRepository) GetUsersVotedPosts(ctx context.Context, userID int) ([]model.Post, error) {
// 	var isUserExists bool

// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return nil, fmt.Errorf("get users voted posts: transaction: begin: %w", err)
// 	}

// 	stmt, err := tx.PrepareContext(ctx, `SELECT EXISTS (SELECT id FROM user WHERE id = $1);`)
// 	if err != nil {
// 		if err = tx.Rollback(); err != nil {
// 			return nil, fmt.Errorf("get users posts: rollback: %w", err)
// 		}
// 		return nil, fmt.Errorf("get users voted posts: prepare statement 1: %w", err)
// 	}
// 	defer stmt.Close()

// 	row := stmt.QueryRowContext(ctx, userID)

// 	row.Scan(&isUserExists)
// 	if err != nil {
// 		return nil, fmt.Errorf("get users voted posts: row scan: %w", err)
// 	}

// 	if !isUserExists {
// 		return nil, fmt.Errorf("user is not exists")
// 	}

// 	var posts []model.Post

// 	stmt, err = tx.PrepareContext(ctx, `
// 	SELECT
// 		p.id,
// 		p.user_id AS author.id,
// 		u.first_name AS author.first_name,
// 		u.last_name AS author.last_name,
// 		p.title,
// 		p.creation_time
// 	FROM
// 		post p
// 		LEFT JOIN user u ON p.user_id = u.id
// 		LEFT JOIN vote_post

// 		`)
// }

func (r *UserRepository) SetSession(ctx context.Context, session model.Session) error {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) DeleteSession(ctx context.Context, userID int, refreshToken string) error {
	panic("not implemented") // TODO: Implement
}
