package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
)

var logger = logging.GetLogger()

type AppUser interface {
	GetUsers(page int) ([]IndTask.User, int, error)
	CreateUser(user *IndTask.User) (*IndTask.User, error)
	ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error)
}

type AppBook interface {
	GetThreeBooks() ([]IndTask.MostPopularBook, error)
	GetBooks(page int) ([]IndTask.BookResponse, int, error)
	CreateBook(*IndTask.Book) (*IndTask.OneBookResponse, error)
	ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.OneBookResponse, error)
	GetListBooks(page int) ([]IndTask.ListBooksResponse, int, error)
	ChangeListBook(books *IndTask.ListBook, bookId int, method string) (*IndTask.ListBooksResponse, error)
	CreateListBook(bookId int) (*IndTask.ListBooksResponse, error)
}

type AppAct interface {
	GetActs(page int) ([]IndTask.Act, int, error)
	CreateIssueAct(act *IndTask.Act) (*IndTask.Act, error)
	GetActsByUser(userId int, page int) ([]IndTask.Act, int, error)
	ChangeAct(act *IndTask.Act, actId int, method string) (*IndTask.Act, error)
	AddReturnAct(act *IndTask.ReturnAct) (*IndTask.Act, error)
}

type AppAuthor interface {
	GetAuthors(page int) ([]IndTask.Author, int, error)
	CreateAuthor(author *IndTask.Author) (*IndTask.Author, error)
	ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error)
}

type AppGenre interface {
	GetGenres() ([]IndTask.Genre, error)
	CreateGenre(genre *IndTask.Genre) (*IndTask.Genre, error)
	ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error)
}

type Validation interface {
	GetGenreById(int) error
	GetAuthorById(int) error
	GetUserById(int) error
	GetListBookById(int) error
	GetActById(int, bool) error
	GetBookById(int) error
}

type Service struct {
	AppUser
	AppBook
	AppAct
	AppAuthor
	AppGenre
	Validation
}

func NewService(rep *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		NewUserService(rep.AppUser),
		NewBookService(*rep, cfg),
		NewActService(*rep),
		NewAuthorService(rep.AppAuthor),
		NewGenreService(rep.AppGenre),
		NewValidationService(rep.Validation),
	}
}
