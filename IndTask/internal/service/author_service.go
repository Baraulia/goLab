package service

import (
	"errors"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

var AuthorAlreadyExists = errors.New("author with that name already exists")
var AuthorDoesNotExists = errors.New("author with that id does not exists")

type AuthorService struct {
	repo repository.AppAuthor
}

func NewAuthorService(repo repository.AppAuthor) *AuthorService {
	return &AuthorService{repo: repo}
}

func (a *AuthorService) GetAuthors(page int) ([]IndTask.Author, error) {
	return a.repo.GetAuthors(page)
}

func (a *AuthorService) CreateAuthor(author *IndTask.Author) (int, error) {
	listAuthors, err := a.repo.GetAuthors(0)
	if err != nil {
		logger.Errorf("Error when getting authors:%s", err)
		return 0, err
	}
	for _, bdAuthor := range listAuthors {
		if bdAuthor.AuthorName == author.AuthorName {
			logger.Error("Author with the same name already exists")
			return bdAuthor.Id, AuthorAlreadyExists
		}
	}
	return a.repo.CreateAuthor(author)
}

func (a *AuthorService) ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error) {
	listAuthors, err := a.repo.GetAuthors(0)
	if err != nil {
		logger.Errorf("Error when getting authors:%s", err)
		return nil, err
	}
	var authorExist = false
	for _, bdAuthor := range listAuthors {
		if bdAuthor.Id == authorId {
			authorExist = true
		}
	}
	if authorExist == false {
		logger.Error("Such a author does not exist")
		return nil, AuthorDoesNotExists
	}

	return a.repo.ChangeAuthor(author, authorId, method)
}
