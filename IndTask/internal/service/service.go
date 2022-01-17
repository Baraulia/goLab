package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

var logger = logging.GetLogger()

type AppUser interface {
	GetUsers(page int) ([]IndTask.User, error)
	CreateUser(user *IndTask.User) (int, error)
	ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error)
}

type AppBook interface {
	GetThreeBooks() ([]IndTask.MostPopularBook, error)
	GetBooks(page int) ([]IndTask.BookResponse, error)
	CreateBook(*IndTask.Book) (int, error)
	ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.BookResponse, error)
	GetListBooks(page int) ([]IndTask.ListBooksResponse, error)
	ChangeListBook(books *IndTask.ListBooks, bookId int, method string) (*IndTask.ListBooksResponse, error)
}

type AppAct interface {
	GetActs(page int) ([]IndTask.Act, error)
	CreateIssueAct(act *IndTask.Act) (int, error)
	GetActsByUser(userId int, page int) ([]IndTask.Act, error)
	ChangeAct(act *IndTask.Act, actId int, method string) (*IndTask.Act, error)
	AddReturnAct(act *IndTask.ReturnAct) error
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
	GetActById(int, bool) error
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
