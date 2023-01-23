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
	SetToken(ctx context.Context, userID int) (string, error)
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

func (s *UserService) SignUp(ctx context.Context, input UserSignUpInput) error {
	panic("not implemented") // TODO: Implement
}

type UserSignInInput struct {
	UsernameOrEmail string
	Password        string
}

func (s *UserService) SignIn(ctx context.Context, input UserSignInInput) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (s *UserService) GetByID(ctx context.Context, userID int) (model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (s *UserService) GetUsersPosts(ctx context.Context, userID int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (s *UserService) GetUsersVotedPosts(ctx context.Context, userID int) ([]model.Post, error) {
	panic("not implemented") // TODO: Implement
}

func (s *UserService) SetToken(ctx context.Context, userID int) (string, error) {
	panic("not implemented") // TODO: Implement
}
