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

func (s *BookService) GetBooks() ([]IndTask.Book, error) {
	return s.repo.GetBooks()
}

func (s *BookService) CreateBook(book *IndTask.Book) (int, error) {
	return s.repo.CreateBook(book)
}

func (s *BookService) ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.Book, error) {
	return s.repo.ChangeBook(book, bookId, method)
}
