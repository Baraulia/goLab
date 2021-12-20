package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type AppUser interface {
	GetUsers()
	CreateUser()
	ChangeUser()
}

type AppBook interface {
	GetBooks() []IndTask.Book
	CreateBook()
	ChangeBook()
}

type AppMove interface {
	MoveInBook()
	MoveOutBook()
}

type AppAuthor interface {
	GetAuthors()
	CreateAuthor()
	ChangeAuthor()
}

type AppGenre interface {
	GetGenres()
	CreateGenre()
	ChangeGenre()
}

type Service struct {
	AppUser
	AppBook
	AppMove
	AppAuthor
	AppGenre
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		NewUserService(rep.AppUser),
		NewBookService(rep.AppBook),
		NewMoveService(rep.AppMove),
		NewAuthorService(rep.AppAuthor),
		NewGenreService(rep.AppGenre),
	}
}
