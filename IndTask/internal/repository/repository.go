package repository

import (
	"database/sql"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
)

var logger = logging.GetLogger()

type AppUser interface {
	GetUsers(page int) ([]IndTask.User, int, error)
	CreateUser(user *IndTask.User) (*IndTask.User, error)
	GetOneUser(userId int) (*IndTask.User, error)
	ChangeUser(user *IndTask.User, userId int) (*IndTask.User, error)
	DeleteUser(userId int) error
}

type AppBook interface {
	GetThreeBooks() ([]IndTask.MostPopularBook, error)
	GetBooks(page int) ([]IndTask.BookResponse, int, error)
	CreateBook(book *IndTask.Book, bookRentCost float64) (*IndTask.OneBookResponse, error)
	GetOneBook(bookId int) (*IndTask.OneBookResponse, error)
	ChangeBook(book *IndTask.Book, bookId int, bookRentCost float64) (*IndTask.OneBookResponse, error)
	DeleteBook(bookId int) error
	GetListBooks(page int) ([]IndTask.ListBooksResponse, int, error)
	GetAuthorsByBookId(bookId int) ([]int, error)
	GetOneListBook(bookId int) (*IndTask.ListBooksResponse, error)
	ChangeListBook(book *IndTask.ListBook, bookId int) (*IndTask.ListBooksResponse, error)
	CreateListBook(bookId int, bookRentCost float64) (*IndTask.ListBooksResponse, error)
	DeleteListBook(bookId int) error
	GetListBooksId() ([]int, error)
	GetBooksId() ([]int, error)
}

type AppAct interface {
	GetActs(page int) ([]IndTask.Act, int, error)
	CreateIssueAct(act *IndTask.Act) (*IndTask.Act, error)
	GetActsByUser(userId int, forCost bool, page int) ([]IndTask.Act, int, error)
	ChangeAct(act *IndTask.Act, actId int) (*IndTask.Act, error)
	GetOneAct(actId int) (*IndTask.Act, error)
	AddReturnAct(returnAct *IndTask.ReturnAct) (*IndTask.Act, error)
	CheckReturnData() ([]IndTask.Debtor, error)
	CheckDuplicateBook(act *IndTask.Act) error
}

type AppAuthor interface {
	GetAuthors(page int) ([]IndTask.Author, int, error)
	CreateAuthor(author *IndTask.Author) (*IndTask.Author, error)
	GetOneAuthor(authorId int) (*IndTask.Author, error)
	ChangeAuthor(author *IndTask.Author, authorId int) (*IndTask.Author, error)
	DeleteAuthor(authorId int) error
}

type AppGenre interface {
	GetGenres() ([]IndTask.Genre, error)
	CreateGenre(genre *IndTask.Genre) (*IndTask.Genre, error)
	GetOneGenre(genreId int) (*IndTask.Genre, error)
	ChangeGenre(genre *IndTask.Genre, genreId int) (*IndTask.Genre, error)
	DeleteGenre(genreId int) error
}

type Validation interface {
	GetGenreById(int) error
	GetAuthorById(int) error
	GetUserById(int) error
	GetListBookById(int) error
	GetActById(int, bool) error
	GetBookById(int) error
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
