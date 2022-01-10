package repository

import (
	"database/sql"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
)

var logger = logging.GetLogger()

type AppUser interface {
	GetUsers(page int) ([]IndTask.User, error)
	CreateUser(user *IndTask.User) (int, error)
	ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error)
}

type AppBook interface {
	GetThreeBooks() ([]IndTask.BookDTO, error)
	GetBooks(page int) ([]IndTask.BookDTO, error)
	CreateBook(*IndTask.Book, bool) (int, error)
	ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.BookDTO, error)
	GetListBooks(page int) ([]IndTask.ListBooksDTO, error)
	GetAuthorsByBookId(bookId int) ([]int, error)
	ChangeListBook(listBook *IndTask.ListBooks, listBookId int, method string) (*IndTask.ListBooksDTO, error)
}

type AppMove interface {
	GetIssueActs(page int) ([]IndTask.IssueAct, error)
	CreateIssueAct(issueAct *IndTask.IssueAct) (int, error)
	GetIssueActsByUser(userId int, forCost bool, page int) ([]IndTask.IssueAct, error)
	ChangeIssueAct(issueAct *IndTask.IssueAct, actId int, method string) (*IndTask.IssueAct, error)
	GetReturnActs(page int) ([]IndTask.ReturnAct, error)
	CreateReturnAct(returnAct *IndTask.ReturnAct, listBookId int) (int, error)
	GetReturnActsByUser(userId int, page int) ([]IndTask.ReturnAct, error)
	ChangeReturnAct(returnAct *IndTask.ReturnAct, actId int, method string) (*IndTask.ReturnAct, error)
	CheckReturnData() ([]IndTask.Debtor, error)
}

type AppAuthor interface {
	GetAuthors(page int) ([]IndTask.Author, error)
	CreateAuthor(author *IndTask.Author) (int, error)
	ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error)
}

type AppGenre interface {
	GetGenres() ([]IndTask.Genre, error)
	CreateGenre(genre *IndTask.Genre) (int, error)
	ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error)
}

type Validation interface {
	GetGenreById(int) error
	GetAuthorById(int) error
	GetUserById(int) error
	GetListBookById(int) error
	GetIssueActById(int) error
}

type Repository struct {
	AppUser
	AppBook
	AppMove
	AppAuthor
	AppGenre
	Validation
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		NewUserPostgres(db),
		NewBookPostgres(db),
		NewMovePostgres(db),
		NewAuthorPostgres(db),
		NewGenrePostgres(db),
		NewValidationPostgres(db),
	}
}
