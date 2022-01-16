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
	GetOneUser(userId int) (*IndTask.User, error)
	ChangeUser(user *IndTask.User, userId int) error
	DeleteUser(userId int) error
}

type AppBook interface {
	GetThreeBooks() ([]IndTask.BookDTO, error)
	GetBooks(page int) ([]IndTask.BookDTO, error)
	CreateBook(*IndTask.Book, bool, float64) (int, error)
	GetOneBook(bookId int) (*IndTask.BookDTO, error)
	ChangeBook(book *IndTask.Book, bookId int, bookRentCost float64) error
	DeleteBook(bookId int) error
	GetListBooks(page int) ([]IndTask.ListBooksDTO, error)
	GetAuthorsByBookId(bookId int) ([]int, error)
	GetOneListBook(bookId int) (*IndTask.ListBooksDTO, error)
	ChangeListBook(book *IndTask.ListBooks, bookId int) error
	DeleteListBook(bookId int) error
}

type AppAct interface {
	GetActs(page int) ([]IndTask.Act, error)
	CreateIssueAct(act *IndTask.Act) (int, error)
	GetActsByUser(userId int, forCost bool, page int) ([]IndTask.Act, error)
	ChangeAct(act *IndTask.Act, actId int) error
	GetOneAct(actId int) (*IndTask.Act, error)
	AddReturnAct(returnAct *IndTask.ReturnAct) error
	CheckReturnData() ([]IndTask.Debtor, error)
	CheckDuplicateBook(act *IndTask.Act) error
}

type AppAuthor interface {
	GetAuthors(page int) ([]IndTask.Author, error)
	CreateAuthor(author *IndTask.Author) (int, error)
	GetOneAuthor(authorId int) (*IndTask.Author, error)
	ChangeAuthor(author *IndTask.Author, authorId int) error
	DeleteAuthor(authorId int) error
}

type AppGenre interface {
	GetGenres() ([]IndTask.Genre, error)
	CreateGenre(genre *IndTask.Genre) (int, error)
	GetOneGenre(genreId int) (*IndTask.Genre, error)
	ChangeGenre(genre *IndTask.Genre, genreId int) error
	DeleteGenre(genreId int) error
}

type Validation interface {
	GetGenreById(int) error
	GetAuthorById(int) error
	GetUserById(int) error
	GetListBookById(int) error
	GetActById(int, bool) error
}

type Repository struct {
	AppUser
	AppBook
	AppAct
	AppAuthor
	AppGenre
	Validation
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		NewUserPostgres(db),
		NewBookPostgres(db),
		NewActPostgres(db),
		NewAuthorPostgres(db),
		NewGenrePostgres(db),
		NewValidationPostgres(db),
	}
}
