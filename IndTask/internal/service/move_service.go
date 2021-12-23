package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type MoveService struct {
	repo repository.AppMove
}

func NewMoveService(repo repository.AppMove) *MoveService {
	return &MoveService{repo: repo}
}

func (s *MoveService) MoveInBook(issueAct *IndTask.IssueAct) (issueActId int, err error) {
	return s.repo.MoveInBook(issueAct)
}

func (s *MoveService) GetMoveInBooks(userId int) ([]IndTask.IssueAct, error) {
	return s.repo.GetMoveInBooks(userId)
}

func (s *MoveService) MoveOutBook(returnAct *IndTask.ReturnAct) (returnActId int, err error) {
	return s.repo.MoveOutBook(returnAct)
}
