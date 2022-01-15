package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type ValidationService struct {
	repo repository.Validation
}

func NewValidationService(repo repository.Validation) *ValidationService {
	return &ValidationService{repo: repo}
}

func (v *ValidationService) GetGenreById(id int) error {
	return fmt.Errorf("validation error:%w", v.repo.GetGenreById(id))
}
func (v *ValidationService) GetAuthorById(id int) error {
	return fmt.Errorf("validation error:%w", v.repo.GetAuthorById(id))
}
func (v *ValidationService) GetUserById(id int) error {
	return fmt.Errorf("validation error:%w", v.repo.GetUserById(id))
}
func (v *ValidationService) GetListBookById(id int) error {
	return fmt.Errorf("validation error:%w", v.repo.GetListBookById(id))
}
func (v *ValidationService) GetIssueActById(id int) error {
	return fmt.Errorf("validation error:%w", v.repo.GetIssueActById(id))
}
