package service

import "github.com/Baraulia/goLab/IndTask.git/internal/repository"

type MoveService struct {
	repo repository.AppMove
}

func NewMoveService(repo repository.AppMove) *MoveService {
	return &MoveService{repo: repo}
}

func (s *MoveService) MoveInBook() {

}

func (s *MoveService) MoveOutBook() {

}
