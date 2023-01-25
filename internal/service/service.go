package service

import (
	"real-time-forum/internal/config"
	"real-time-forum/internal/repository"
	hash "real-time-forum/pkg/hasher"
)

type Service struct {
	User User
}

func NewService(
	repo *repository.Repository,
	h *hash.HasherService,
	cfg *config.Config) *Service {
	userService := NewUser(repo.User, h, cfg)

	return &Service{
		User: userService,
	}
}
