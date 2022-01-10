package service

import (
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type ValidationService struct {
	repo repository.Validation
}

func NewValidationService(repo repository.Validation) *ValidationService {
	return &ValidationService{repo: repo}
}

func (v *ValidationService) GetGenreById(id int) error {
	return v.repo.GetGenreById(id)
}
func (v *ValidationService) GetAuthorById(id int) error {
	return v.repo.GetAuthorById(id)
}
func (v *ValidationService) GetUserById(id int) error {
	return v.repo.GetUserById(id)
}
func (v *ValidationService) GetListBookById(id int) error {
	return v.repo.GetListBookById(id)
}
func (v *ValidationService) GetIssueActById(id int) error {
	return v.repo.GetIssueActById(id)
}
