package repository

import (
	"database/sql"
	"real-time-forum/internal/model"
)

type Post interface {
	Create(post model.Post) (int, error)
	GetByID(postID int, userID int) (model.Post, error)
	Delete(userID int, postID int) error
	GetPostsByCategoryID(categoryID int, limit int, offset int) ([]model.Post, error)
	LikePost(like model.PostVotes) (bool, error)
	DislikePost(dislike model.PostVotes) (bool, error)
}

type PostRepository struct {
	db *sql.DB
}

func NewPost(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Create(post model.Post) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) GetByID(postID int, userID int) (model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) Delete(userID int, postID int) error {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) GetPostsByCategoryID(categoryID int, limit int, offset int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) LikePost(like model.PostVotes) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepository) DislikePost(dislike model.PostVotes) (bool, error) {
	panic("not implemented") // TODO: Implement
}
