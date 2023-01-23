package service

import "real-time-forum/internal/repository"

type Service struct {
	User User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUser(repo.User),
	}
}
