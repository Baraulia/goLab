package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type BookService struct {
	repo repository.AppBook
}

func NewBookService(repo repository.AppBook) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetBooks() []IndTask.Book {
	return s.repo.GetBooks()
}

func (s *BookService) CreateBook() {

}

func (s *BookService) ChangeBook() {

}
