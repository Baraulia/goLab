package repository

import (
	"database/sql"
	"github.com/Baraulia/goLab/IndTask.git"
)

type AppUser interface {
	GetUsers()
	CreateUser()
	ChangeUser()
}

type AppBook interface {
	GetBooks() []IndTask.Book
	CreateBook(IndTask.Book) (int, error)
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

type Repository struct {
	AppUser
	AppBook
	AppMove
	AppAuthor
	AppGenre
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		NewUserPostgres(db),
		NewBookPostgres(db),
		NewMovePostgres(db),
		NewAuthorPostgres(db),
		NewGenrePostgres(db),
	}
}
