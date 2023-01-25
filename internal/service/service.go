package service

import (
	"real-time-forum/internal/repository"
	hash "real-time-forum/pkg/hasher"
)

type Service struct {
	User User
}

func NewService(repo *repository.Repository, h *hash.HasherService) *Service {
	userService := NewUser(repo.User, h)

	return &Service{
		User: userService,
	}
}
