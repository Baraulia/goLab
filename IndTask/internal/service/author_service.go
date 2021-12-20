package service

import "github.com/Baraulia/goLab/IndTask.git/internal/repository"

type AuthorService struct {
	repo repository.AppAuthor
}

func NewAuthorService(repo repository.AppAuthor) *AuthorService {
	return &AuthorService{repo: repo}
}

func (u *AuthorService) GetAuthors() {

}

func (u *AuthorService) CreateAuthor() {

}

func (u *AuthorService) ChangeAuthor() {

}
