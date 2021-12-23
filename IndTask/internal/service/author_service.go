package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type AuthorService struct {
	repo repository.AppAuthor
}

func NewAuthorService(repo repository.AppAuthor) *AuthorService {
	return &AuthorService{repo: repo}
}

func (a *AuthorService) GetAuthors() ([]IndTask.Author, error) {
	return a.repo.GetAuthors()
}

func (a *AuthorService) CreateAuthor(author *IndTask.Author) (int, error) {
	return a.repo.CreateAuthor(author)
}

func (a *AuthorService) ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error) {
	return a.repo.ChangeAuthor(author, authorId, method)
}
