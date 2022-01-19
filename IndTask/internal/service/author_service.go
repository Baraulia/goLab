package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/translate"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type AuthorService struct {
	repo repository.AppAuthor
}

func NewAuthorService(repo repository.AppAuthor) *AuthorService {
	return &AuthorService{repo: repo}
}

func (a *AuthorService) GetAuthors(page int) ([]IndTask.Author, int, error) {
	authors, pages, err := a.repo.GetAuthors(page)
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting authors from database:%w", err)
	}
	return authors, pages, nil
}

func (a *AuthorService) CreateAuthor(author *IndTask.Author) (*IndTask.Author, error) {
	listAuthors, _, err := a.repo.GetAuthors(0)
	if err != nil {
		return nil, fmt.Errorf("error while getting authors from database:%w", err)
	}
	for _, bdAuthor := range listAuthors {
		if bdAuthor.AuthorName == author.AuthorName {
			logger.Errorf("Author with that name:%s already exists", author.AuthorName)
			return nil, fmt.Errorf("author with that name:%s already exists", author.AuthorName)
		}
	}
	newAuthor, err := a.repo.CreateAuthor(author)
	if err != nil {
		return nil, fmt.Errorf("error while creating author in database:%w", err)
	}
	return newAuthor, nil
}

func (a *AuthorService) ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error) {
	listAuthors, _, err := a.repo.GetAuthors(0)
	if err != nil {
		return nil, fmt.Errorf("error while getting authors from database:%w", err)
	}
	var authorExist = false
	for _, bdAuthor := range listAuthors {
		if bdAuthor.Id == authorId {
			authorExist = true
		}
	}
	if authorExist == false {
		logger.Errorf("Such a author:%d does not exist", authorId)
		return nil, fmt.Errorf("such a author:%d does not exist", authorId)
	}
	if method == "GET" {
		oneAuthor, err := a.repo.GetOneAuthor(authorId)
		if err != nil {
			return nil, fmt.Errorf("error while getting one author from database:%w", err)
		}
		return oneAuthor, nil
	}
	if method == "PUT" {
		for _, bdAuthor := range listAuthors {
			if bdAuthor.AuthorName == author.AuthorName {
				logger.Errorf("Author with that name:%s already exists", author.AuthorName)
				return nil, fmt.Errorf("author with that name:%s already exists", author.AuthorName)
			}
		}
		newAuthor, err := a.repo.ChangeAuthor(author, authorId)
		if err != nil {
			return nil, fmt.Errorf("error while changing author in database:%w", err)
		}
		return newAuthor, nil
	}
	if method == "DELETE" {
		err := a.repo.DeleteAuthor(authorId)
		if err != nil {
			return nil, fmt.Errorf("error while deleting one author from database:%w", err)
		}
		return nil, nil
	}
	return nil, nil
}

func InputAuthorFoto(req *http.Request, input *IndTask.Author) error {
	reqFile, fileHeader, err := req.FormFile("file")
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while getting file from multipart form:%s", err)
		return fmt.Errorf("inputAuthorFoto: error while getting file from multipart form:%w", err)
	}
	defer reqFile.Close()
	filePath := fmt.Sprintf("images/authors/%s.%s", translate.Translate(input.AuthorName), (strings.Split(fileHeader.Filename, "."))[1])
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while opening file %s:%s", filePath, err)
		return fmt.Errorf("inputAuthorFoto: error while opening file %s:%w", filePath, err)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(reqFile)
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while reading file from request:%s", err)
		return fmt.Errorf("inputAuthorFoto: error while reading file from request:%w", err)
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while writing file:%s", err)
		return fmt.Errorf("inputAuthorFoto: error while writing file:%w", err)
	}
	input.AuthorFoto = filePath
	return nil
}
