package repository

import (
	"context"
	"database/sql"
	"fmt"

	"real-time-forum/internal/model"
)

type User interface {
	Create(ctx context.Context, user model.User) error
	GetByEmailOrUsername(ctx context.Context, emailOrUsername string) (string, error)
	GetByCredentials(ctx context.Context, usernameOrEmail, password string) (model.User, error)
	GetByID(ctx context.Context, userID int) (model.User, error)
	GetUsersPosts(ctx context.Context, userID int) ([]model.Post, error)
	GetUsersVotedPosts(ctx context.Context, userID int) ([]model.Post, error)
	SetSession(ctx context.Context, session model.Session) error
	DeleteSession(ctx context.Context, userID int) error
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

func (r *UserRepository) GetByEmailOrUsername(ctx context.Context, emailOrUsername string) (string, error) {
	var password string

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT 
			password
		FROM
			user
		WHERE
			(username = $1 OR email = $1);`)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, emailOrUsername).Scan(&password); err != nil {
		return "", err
	}

	return password, nil
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
		return nil, fmt.Errorf("repo: get users posts: transaction: begin: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `SELECT EXISTS (SELECT id FROM user WHERE id = $1);`)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("repo: get users posts: prepare statement 1: %w", err)
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, userID).Scan(&isUserExists); err != nil {
		return nil, fmt.Errorf("repo: get users posts: row scan: %w", err)
	}

	if !isUserExists {
		return nil, ErrUserExists
	}

	var posts []model.Post

	stmt, err = tx.PrepareContext(ctx, `
		SELECT 
			post.id, 
			post.title, 
			post.content, 
			post.creation_time, 
			post.image, 
			user.id AS author_id,
			user.username AS author_username, 
			user.first_name AS author_first_name, 
			user.last_name AS author_last_name, 
			user.avatar AS author_avatar
		FROM 
			post LEFT JOIN user ON post.user_id = user.id
		WHERE 
			post.user_id = $1
		ORDER BY 1 DESC;`)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("repo: get users posts: prepare statement 2: %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get users posts: exec statemnt: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err := rows.Scan(
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

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("repo: get users posts: transaction: commit %w", err)
	}

	return posts, nil
}

func (r *UserRepository) GetUsersVotedPosts(ctx context.Context, userID int) ([]model.Post, error) {
	var isUserExists bool

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo: get users voted posts: transaction: begin: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `SELECT EXISTS (SELECT id FROM user WHERE id = $1);`)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("repo: get users voted posts: prepare statement 1: %w", err)
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, userID).Scan(&isUserExists); err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("repo: get users voted posts: row scan: %w", err)
	}

	if !isUserExists {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
		}
		return nil, ErrUserExists
	}

	var posts []model.Post

	stmt, err = tx.PrepareContext(ctx, `
		SELECT
			post.id,
			post.user_id AS author_id,
			user.first_name AS author_first_name,
			user.last_name AS author_last_name,
			post.title,
			post.creation_time,
			vote_post.vote AS user_vote
		FROM
			post
		LEFT JOIN 
			user ON post.user_id = user.id
		LEFT JOIN 
			vote_post ON post.id = vote_post.post_id AND vote_post.user_id = $1
		WHERE
			post.id IN (
				SELECT
					post_id
				FROM
					vote_post
				WHERE
					user_id = $1
			)
			ORDER BY post.id DESC;
		`)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("repo: get users voted posts: prepare %w", err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
		}
		return nil, fmt.Errorf("repo: get users voted posts: query %w", err)
	}

	for rows.Next() {
		var post model.Post

		err := rows.Scan(
			&post.ID,
			&post.Author.ID,
			&post.Author.FirstName,
			&post.Author.LastName,
			&post.Title,
			&post.CreationTime,
			&post.Rating,
		)
		if err != nil {
			if err = tx.Rollback(); err != nil {
				return nil, fmt.Errorf("repo: get users posts: rollback: %w", err)
			}
			return nil, fmt.Errorf("repo: get users voted posts: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("repo: get users voted posts: commit %w", err)
	}

	return posts, nil
}

func (r *UserRepository) SetSession(ctx context.Context, session model.Session) error {
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO
			session_tokens (user_id, token, token_expiration_time)
		VALUES
			($1, $2, $3);`)
	if err != nil {
		return fmt.Errorf("repo: set session: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, session.UserID, session.Token, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("repo: set session: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteSession(ctx context.Context, userID int) error {
	res, err := r.db.Exec(`
		DELETE FROM 
			session_tokens
		WHERE
			user_id = $1`, userID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return err
	}

	return nil
}
