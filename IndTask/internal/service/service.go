package service

import "github.com/Baraulia/goLab/IndTask.git/internal/repository"

type AppUser interface {
}

type AppBook interface {
}

type AppMove interface {
}

type Service struct {
	AppUser
	AppBook
	AppMove
}

func NewService(rep *repository.Repository) *Service {
	return &Service{}
}
