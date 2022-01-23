package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/myErrors"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/translate"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type AuthorService struct {
	repo repository.Repository
	cfg  *config.Config
}

func NewAuthorService(repo repository.Repository, cfg *config.Config) *AuthorService {
	return &AuthorService{repo: repo, cfg: cfg}
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
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting authors from database:%w", err), Code: 500}
	}
	for _, bdAuthor := range listAuthors {
		if bdAuthor.AuthorName == author.AuthorName {
			logger.Errorf("Author with that name:%s already exists", author.AuthorName)
			return nil, &myErrors.MyError{Err: fmt.Errorf("author with that name:%s already exists", author.AuthorName), Code: 400}
		}
	}
	newAuthor, err := a.repo.CreateAuthor(author)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while creating author in database:%w", err), Code: 500}
	}
	return newAuthor, nil
}

func (a *AuthorService) ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error) {
	err := a.repo.Validation.GetAuthorById(authorId)
	if err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			logger.Errorf("Such author:%d does not exist", authorId)
			return nil, &myErrors.MyError{Err: fmt.Errorf("such author:%d does not exist", authorId), Code: e.Status()}
		default:
			logger.Errorf("changeAuthor:%s", err)
			return nil, fmt.Errorf("changeAuthor:%w", err)
		}
	}
	listAuthors, _, err := a.repo.GetAuthors(0)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting authors from database:%w", err), Code: 500}
	}
	var authorExist = false
	for _, bdAuthor := range listAuthors {
		if bdAuthor.Id == authorId {
			authorExist = true
		}
	}
	if authorExist == false {
		logger.Errorf("Such a author:%d does not exist", authorId)
		return nil, &myErrors.MyError{Err: fmt.Errorf("such a author:%d does not exist", authorId), Code: 400}
	}
	if method == "GET" {
		oneAuthor, err := a.repo.GetOneAuthor(authorId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting one author from database:%w", err), Code: 500}
		}
		return oneAuthor, nil
	}
	if method == "PUT" {
		for _, bdAuthor := range listAuthors {
			if bdAuthor.AuthorName == author.AuthorName {
				logger.Errorf("Author with that name:%s already exists", author.AuthorName)
				return nil, &myErrors.MyError{Err: fmt.Errorf("author with that name:%s already exists", author.AuthorName), Code: 400}
			}
		}
		newAuthor, err := a.repo.ChangeAuthor(author, authorId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while changing author in database:%w", err), Code: 500}
		}
		return newAuthor, nil
	}
	if method == "DELETE" {
		err := a.repo.DeleteAuthor(authorId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while deleting one author from database:%w", err), Code: 500}
		}
		return nil, nil
	}
	return nil, nil
}

func (a *AuthorService) InputAuthorFoto(req *http.Request, input *IndTask.Author) error {
	reqFile, fileHeader, err := req.FormFile("file")
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while getting file from multipart form:%s", err)
		return &myErrors.MyError{Err: fmt.Errorf("inputAuthorFoto: error while getting file from multipart form:%w", err), Code: 400}
	}
	defer reqFile.Close()
	filePath := fmt.Sprintf("%simages/authors/%s.%s", a.cfg.FilePath, translate.Translate(input.AuthorName), (strings.Split(fileHeader.Filename, "."))[1])
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while opening file %s:%s", filePath, err)
		return &myErrors.MyError{Err: fmt.Errorf("inputAuthorFoto: error while opening file %s:%w", filePath, err), Code: 500}
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(reqFile)
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while reading file from request:%s", err)
		return &myErrors.MyError{Err: fmt.Errorf("inputAuthorFoto: error while reading file from request:%w", err), Code: 500}
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		logger.Errorf("InputAuthorFoto: error while writing file:%s", err)
		return &myErrors.MyError{Err: fmt.Errorf("inputAuthorFoto: error while writing file:%w", err), Code: 500}
	}
	filePath = strings.Replace(filePath, a.cfg.FilePath, "", 1)
	input.AuthorFoto = filePath
	return nil
}
