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
	err := v.repo.GetGenreById(id)
	if err != nil {
		logger.Errorf("validation error:%s", err)
		return fmt.Errorf("validation error:%s", err)
	}
	return nil
}
func (v *ValidationService) GetAuthorById(id int) error {
	err := v.repo.GetAuthorById(id)
	if err != nil {
		logger.Errorf("validation error:%s", err)
		return fmt.Errorf("validation error:%s", err)
	}
	return nil
}
func (v *ValidationService) GetUserById(id int) error {
	err := v.repo.GetUserById(id)
	if err != nil {
		logger.Errorf("validation error:%s", err)
		return fmt.Errorf("validation error:%s", err)
	}
	return nil
}
func (v *ValidationService) GetListBookById(id int) error {
	err := v.repo.GetListBookById(id)
	if err != nil {
		logger.Errorf("validation error:%s", err)
		return fmt.Errorf("validation error:%s", err)
	}
	return nil
}
func (v *ValidationService) GetActById(id int, changing bool) error {
	err := v.repo.GetActById(id, changing)
	if err != nil {
		logger.Errorf("validation error:%s", err)
		return fmt.Errorf("validation error:%s", err)
	}
	return nil
}
