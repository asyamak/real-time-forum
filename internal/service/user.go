package service

import (
	"context"
	"real-time-forum/internal/model"
	"real-time-forum/internal/repository"
)

type User interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (string, error)
	GetByID(ctx context.Context, userID int) (model.User, error)
	GetUsersPosts(ctx context.Context, userID int) ([]model.Post, error)
	GetUsersVotedPosts(ctx context.Context, userID int) ([]model.Post, error)
	SetTokens(ctx context.Context, userID int) (string, error)
}

type UserService struct {
	repo repository.User
}

func NewUser(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserSignUpInput struct {
	Username  string
	FirstName string
	LastName  string
	Age       int
	Gender    int
	Email     string
	Password  string
}

type UserSignInInput struct {
	UsernameOrEmail string
	Password        string
}
