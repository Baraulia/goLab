package repository

import (
	"database/sql"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
)

var logger = logging.GetLogger()

type AppUser interface {
	GetUsers() ([]IndTask.User, error)
	CreateUser(user *IndTask.User) (int, error)
	ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error)
}

type AppBook interface {
	GetBooks() ([]IndTask.Book, error)
	CreateBook(*IndTask.Book) (int, error)
	ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.Book, error)
}

type AppMove interface {
	MoveInBook(issueAct *IndTask.IssueAct) (issueActId int, err error)
	GetMoveInBooks(userId int) ([]IndTask.IssueAct, error)
	MoveOutBook(returnAct *IndTask.ReturnAct) (returnActId int, err error)
}

type AppAuthor interface {
	GetAuthors() ([]IndTask.Author, error)
	CreateAuthor(author *IndTask.Author) (int, error)
	ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error)
}

type AppGenre interface {
	GetGenres() ([]IndTask.Genre, error)
	CreateGenre(genre *IndTask.Genre) (int, error)
	ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error)
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
