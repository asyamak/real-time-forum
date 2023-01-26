package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"real-time-forum/internal/model"
)

type Post interface {
	Create(ctx context.Context, post model.Post) (int, error)
	GetByID(ctx context.Context, postID int, userID int) (model.Post, error)
	Delete(ctx context.Context, userID int, postID int) error
	GetPostsByCategoryID(ctx context.Context, categoryID int, limit int, offset int) ([]model.Post, error)
	LikePost(ctx context.Context, like model.PostVotes) (bool, error)
	DislikePost(ctx context.Context, dislike model.PostVotes) (bool, error)
}

type PostRepository struct {
	db *sql.DB
}

func NewPost(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Create(ctx context.Context, post model.Post) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("repo: create post: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO
			post (user_id, title, content, creation_time, image)
		VALUES
			($1, $2, $3, $4, $5) 
		RETURNING id`)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repo: create post: %w", err)
	}

	var id int
	row := stmt.QueryRowContext(ctx,
		&post.Author.ID,
		&post.Title,
		&post.Content,
		&post.CreationTime,
		&post.ImagePath,
	)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repo: create post: %w", err)
	}

	for _, category := range post.Categories {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO
				post_category (post_id, category_id)
			VALUES
				($1, $2);`)
		if err != nil {
			return 0, fmt.Errorf("repo: create post: %w", err)
		}

		_, err = stmt.Exec(&id, &category.ID)
		if err != nil {
			tx.Rollback()
			if isForeignKeyConstraintError(err) {
				return 0, fmt.Errorf("repo: create post: %w", ErrForeignKeyConstraint)
			}
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("repo: create post: %w", err)
	}

	return id, nil
}

func (r *PostRepository) GetByID(ctx context.Context, postID int, userID int) (model.Post, error) {
	var post model.Post
	var isPostExists bool

	row := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT id FROM post WHERE id = $1`, postID)
	if err := row.Scan(&isPostExists); err != nil {
		return model.Post{}, fmt.Errorf("repo: get post: %w", err)
	}

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT
			post.id,
			post.user_id AS author_id
			user.first_name AS author_first_name
			user.last_name AS author_last_name
			post.title
			post.content
			post.creation_time
			post.image
			IFNULL(vote_post.vote, 0) AS user_vote,
			COUNT(DISTINCT pl.id) - COUNT(DISTINCT pd.id) AS vote
		FROM 
			post 
		LEFT JOIN user 
		ON post.user_id = user.id 
		
		LEFT JOIN vote_post pr 
		ON pr.post_id = posts.id 
		AND pr.user_id = $1 
		
		LEFT JOIN vote_post pl 
		ON pl.post_id = posts.id 
		AND pl.type = $2 
		
		LEFT JOIN vote_post pd ON pd.post_id = post.id 
		AND pd.type = $3
	WHERE 
		post.id = $4`)
	if err != nil {
		return model.Post{}, fmt.Errorf("repo: get post: %w", err)
	}

	row = stmt.QueryRowContext(ctx, 1, -1, postID)
	if err := row.Scan(
		&post.ID,
		&post.Author.ID,
		&post.Author.FirstName,
		&post.Author.LastName,
		&post.Title,
		&post.Content,
		&post.CreationTime,
		&post.ImagePath,
		&post.UserRate,
		&post.Rating,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Post{}, ErrNoRows
		}
		return model.Post{}, fmt.Errorf("repo: get post: %w", err)
	}

	post.Categories, err = r.getPostCategories(postID)
	if err != nil {
		return model.Post{}, fmt.Errorf("repo: get post: %w", err)
	}

	return post, nil
}

func (r *PostRepository) getPostCategories(postID int) ([]model.Category, error) {
	var categories []model.Category

	rows, err := r.db.Query(`
		SELECT
			id, name
		FROM
			categories
		WHERE
			id IN (
					SELECT
						category_id
					FROM
						posts_category
					WHERE
						post_id = $1
			)
		`, postID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category model.Category

		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *PostRepository) Delete(ctx context.Context, userID int, postID int) error {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) GetPostsByCategoryID(ctx context.Context, categoryID int, limit int, offset int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) LikePost(ctx context.Context, like model.PostVotes) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) DislikePost(ctx context.Context, dislike model.PostVotes) (bool, error) {
	panic("not implemented") // TODO: Implement
}
