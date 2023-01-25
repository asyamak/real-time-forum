package service

import (
	"context"
	"errors"
	"fmt"
	"real-time-forum/internal/model"
	"real-time-forum/internal/repository"
	"real-time-forum/pkg/hasher"
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
	DeleteToken(ctx context.Context, userID int) error
}

type UserService struct {
	repo   repository.User
	hasher *hasher.HasherService
}

func NewUser(repo repository.User, hasher *hasher.HasherService) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
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

	input.Password = s.hasher.HashPassword(input.Password)

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
	input.Password = s.hasher.HashPassword(input.Password)

	user, err := s.repo.GetByCredentials(ctx, input.UsernameOrEmail, input.Password)
	if err != nil {
		return "", fmt.Errorf("get by credentials: %w", err)
	}

	return s.SetToken(ctx, user.ID)
}

var ErrUserDoesNotExists = errors.New("user doesn't exists")

func (s *UserService) GetByID(ctx context.Context, userID int) (model.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return model.User{}, ErrUserDoesNotExists
		}
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) GetUsersPosts(ctx context.Context, userID int) ([]model.Post, error) {
	posts, err := s.repo.GetUsersPosts(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, ErrUserDoesNotExists
		}
		return nil, err
	}
	return posts, nil
}

func (s *UserService) GetUsersVotedPosts(ctx context.Context, userID int) ([]model.Post, error) {
	likedPosts, err := s.repo.GetUsersVotedPosts(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, ErrUserDoesNotExists
		}
		return nil, err
	}
	return likedPosts, nil
}

func (s *UserService) SetToken(ctx context.Context, userID int) (string, error) {
	tokenUUID, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	session := model.Session{
		UserID:    userID,
		Token:     tokenUUID.String(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}

	return session.Token, s.repo.SetSession(ctx, session)
}

func (s *UserService) DeleteToken(ctx context.Context, userID int) error {
	if err := s.repo.DeleteSession(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return ErrUserDoesNotExists
		}
		return err
	}

	return nil
}
