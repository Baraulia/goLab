package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"net/http"
)

var logger = logging.GetLogger()

type AppUser interface {
	GetUsers(page int, sorting string) ([]IndTask.UserResponse, int, error)
	CreateUser(user *IndTask.User) (*IndTask.User, error)
	ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error)
	FoundUser(userSurname string) (*IndTask.User, error)
}

type AppBook interface {
	GetThreeBooks() ([]IndTask.MostPopularBook, error)
	GetBooks(page int, sorting string) ([]*IndTask.BookResponse, int, error)
	CreateBook(*IndTask.Book) (*IndTask.OneBookResponse, error)
	ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.OneBookResponse, error)
	GetListBooks(page int) ([]IndTask.ListBooksResponse, int, error)
	ChangeListBook(books *IndTask.ListBook, bookId int, method string) (*IndTask.ListBooksResponse, error)
	CreateListBook(bookId int) (*IndTask.ListBooksResponse, error)
	InputCoverFoto(req *http.Request, book *IndTask.Book) error
}

type AppAct interface {
	GetActs(page int) ([]IndTask.Act, int, error)
	CreateIssueAct(act *IndTask.Act) (*IndTask.Act, error)
	GetActsByUser(userId int, page int) ([]IndTask.Act, int, error)
	ChangeAct(act *IndTask.Act, actId int, method string) (*IndTask.Act, error)
	AddReturnAct(act *IndTask.ReturnAct) (*IndTask.Act, error)
	InputFineFoto(req *http.Request, actId int) ([]string, error)
}

type AppAuthor interface {
	GetAuthors(page int) ([]IndTask.Author, int, error)
	CreateAuthor(author *IndTask.Author) (*IndTask.Author, error)
	ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error)
	InputAuthorFoto(req *http.Request, author *IndTask.Author) error
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

type AppFile interface {
	GetFile(path string) ([]byte, error)
}

type Service struct {
	AppUser
	AppBook
	AppAct
	AppAuthor
	AppGenre
	Validation
	AppFile
}

func NewService(rep *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		NewUserService(rep.AppUser),
		NewBookService(*rep, cfg),
		NewActService(*rep, cfg),
		NewAuthorService(rep.AppAuthor, cfg),
		NewGenreService(rep.AppGenre),
		NewValidationService(rep.Validation),
		NewFileService(cfg, logger),
	}
}

const (
	//sorting for books
	bookNameDesc  = "books.book_name DESC"
	bookNameAsc   = "books.book_name ASC"
	publishedDesc = "books.published DESC"
	publishedAsc  = "books.published ASC"
	amountDesc    = "books.amount DESC"
	amountAsc     = "books.amount ASC"
	avAmountDesc  = "av_books DESC"
	avAmountAsc   = "av_books ASC"
	//sorting for users
	userSurnameDesc = "surname DESC"
	userSurnameAsc  = "surname ASC"
	userNameDesc    = "user_name DESC"
	userNameAsc     = "user_name ASC"
	emailDesc       = "email DESC"
	emailAsc        = "email ASC"
	addressDesc     = "address DESC"
	addressAsc      = "address ASC"
	birthDateDesc   = "birth_date DESC"
	birthDateAsc    = "birth_date ASC"
)

func SortTypeBook(sorting string) string {
	switch sorting {
	case "bookNameDesc":
		return bookNameDesc
	case "bookNameAsc":
		return bookNameAsc
	case "publishedDesc":
		return publishedDesc
	case "publishedAsc":
		return publishedAsc
	case "amountDesc":
		return amountDesc
	case "amountAsc":
		return amountAsc
	case "avAmountDesc":
		return avAmountDesc
	case "avAmountAsc":
		return avAmountAsc
	}
	return "av_books DESC, books.book_name"
}

func SortTypeUser(sorting string) string {
	switch sorting {
	case "userSurnameDesc":
		return userSurnameDesc
	case "userSurnameAsc":
		return userSurnameAsc
	case "userNameDesc":
		return userNameDesc
	case "userNameAsc":
		return userNameAsc
	case "emailDesc":
		return emailDesc
	case "emailAsc":
		return emailAsc
	case "addressDesc":
		return addressDesc
	case "addressAsc":
		return addressAsc
	case "birthDateDesc":
		return birthDateDesc
	case "birthDateAsc":
		return birthDateAsc
	}
	return "surname"
}
