package service

import (
	"context"
	"fmt"
	"real-time-forum/internal/model"
	"real-time-forum/internal/repository"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	Gender    string
	Email     string
	Password  string
}

func (s *UserService) SignUp(ctx context.Context, input UserSignUpInput) error {
	var avatar string

	switch input.Gender {
	case "Male":
		avatar = "./database/images/male_default.jpg"
	case "Female":
		avatar = "./database/images/female_default.jpg"
	default:
		return fmt.Errorf("unknown gender")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username:     input.Username,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Age:          input.Age,
		Gender:       input.Gender,
		Email:        strings.ToLower(input.Email),
		Password:     string(password),
		CreationTime: time.Now(),
		Avatar:       avatar,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
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
