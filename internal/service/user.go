package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"real-time-forum/internal/model"
	"real-time-forum/internal/repository"
	"strings"
	"time"

	"github.com/gofrs/uuid"
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

	hash := sha256.New()
	hash.Write([]byte(input.Password))
	input.Password = fmt.Sprintf("%x", hash.Sum([]byte("aboba")))

	user := model.User{
		Username:     input.Username,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Age:          input.Age,
		Gender:       input.Gender,
		Email:        strings.ToLower(input.Email),
		Password:     input.Password,
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
	hash := sha256.New()
	hash.Write([]byte(input.Password))
	input.Password = fmt.Sprintf("%x", hash.Sum([]byte("aboba")))

	user, err := s.repo.GetByCredentials(ctx, input.UsernameOrEmail, input.Password)
	if err != nil {
		return "", fmt.Errorf("get by credentials: %w", err)
	}

	tokenUUID, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	session := model.Session{
		UserID:    user.ID,
		Token:     tokenUUID.String(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}

	if err := s.repo.SetSession(ctx, session); err != nil {
		return "", fmt.Errorf("set session: %w", err)
	}

	return session.Token, nil
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
