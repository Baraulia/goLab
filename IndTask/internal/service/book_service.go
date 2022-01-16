package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/translate"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
)

type BookService struct {
	repo repository.Repository
	cfg  *config.Config
}

func NewBookService(repo repository.Repository, cfg *config.Config) *BookService {
	return &BookService{repo: repo, cfg: cfg}
}

func (b *BookService) GetThreeBooks() ([]IndTask.BookDTO, error) {
	books, err := b.repo.GetThreeBooks()
	if err != nil {
		return nil, fmt.Errorf("error while getting three books from database:%w", err)
	}
	return books, nil
}

func (b *BookService) GetBooks(page int) ([]IndTask.BookDTO, error) {
	books, err := b.repo.GetBooks(page)
	if err != nil {
		return nil, fmt.Errorf("error while getting books from database:%w", err)
	}
	return books, nil
}

func (b *BookService) GetListBooks(page int) ([]IndTask.ListBooksDTO, error) {
	books, err := b.repo.GetListBooks(page)
	if err != nil {
		return nil, fmt.Errorf("error while getting instances of books from database:%w", err)
	}
	return books, nil
}

func (b *BookService) CreateBook(book *IndTask.Book) (int, error) {
	listBooks, err := b.repo.GetBooks(0)
	if err != nil {
		return 0, fmt.Errorf("error while getting books from database:%w", err)
	}
	bookExists := false
	for _, bdBook := range listBooks {
		if bdBook.BookName == book.BookName {
			if bdBook.Published == book.Published {
				authorsId, err := b.repo.GetAuthorsByBookId(bdBook.Id)
				if err != nil {
					logger.Errorf("Error when getting authorsId from bookId = %d:%s", bdBook.Id, err)
					return 0, fmt.Errorf("error when getting authorsId from bookId = %d:%w", bdBook.Id, err)
				}
				if len(authorsId) == len(book.AuthorsId) {
					amountAuthors := 0
					for _, bdAuthor := range authorsId {
						for _, reqAuthor := range book.AuthorsId {
							if bdAuthor == reqAuthor {
								amountAuthors = amountAuthors + 1
							}
						}
					}
					if amountAuthors == len(authorsId) {
						bookExists = true
					}
				}
			}
		}
	}
	if bookExists {
		logger.Infof("BookExists = true for book_name = %s, book_published = %d", book.BookName, book.Published)
		if book.Cover != "" {
			if err := os.Remove(book.Cover); err != nil {
				logger.Errorf("BookExists = true, error deleting file %s:%s", book.Cover, err)
				return 0, fmt.Errorf("BookExists = true, error deleting file %s:%s", book.Cover, err)
			}
		}
	}
	bookRentCost := b.CalcRentCost(book.Cost)
	bookId, err := b.repo.CreateBook(book, bookExists, bookRentCost)
	if err != nil {
		return 0, fmt.Errorf("error while creating book in database:%w", err)
	}
	return bookId, nil
}

func (b *BookService) ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.BookDTO, error) {
	listBooks, err := b.repo.GetBooks(0)
	if err != nil {
		return nil, fmt.Errorf("error while getting books from database:%w", err)
	}

	var bookExist = false
	for _, bdBook := range listBooks {
		if bdBook.Id == bookId {
			bookExist = true
		}
	}
	if bookExist == false {
		logger.Errorf("Such a book:%d does not exist", bookId)
		return nil, fmt.Errorf("such a book:%d does not exist", bookId)
	}
	if method == "GET" {
		book, err := b.repo.GetOneBook(bookId)
		if err != nil {
			return nil, fmt.Errorf("error while getting one book from database:%w", err)
		}
		return book, nil
	}
	if method == "PUT" {
		bookRentCost := b.CalcRentCost(book.Cost)
		err := b.repo.ChangeBook(book, bookId, bookRentCost)
		if err != nil {
			return nil, fmt.Errorf("error while changing book in database:%w", err)
		}
	}
	if method == "DELETE" {
		err := b.repo.DeleteBook(bookId)
		if err != nil {
			return nil, fmt.Errorf("error while deleting one book from database:%w", err)
		}
		return nil, nil
	}
	return nil, nil
}

func (b *BookService) ChangeListBook(listBook *IndTask.ListBooks, listBookId int, method string) (*IndTask.ListBooksDTO, error) {
	listBooks, err := b.repo.GetListBooks(0)
	if err != nil {
		return nil, fmt.Errorf("error while getting instances of books from database:%w", err)
	}
	var bookListExist = false
	for _, listBook := range listBooks {
		if listBook.Id == listBookId {
			bookListExist = true
		}
	}
	if bookListExist == false {
		logger.Errorf("Such a instance of book:%d does not exist", listBookId)
		return nil, fmt.Errorf("such a instance of book:%d does not exist", listBookId)
	}

	if method == "GET" {
		book, err := b.repo.GetOneListBook(listBookId)
		if err != nil {
			return nil, fmt.Errorf("error while getting one instance of book from database:%w", err)
		}
		return book, nil
	}
	if method == "PUT" {
		err := b.repo.ChangeListBook(listBook, listBookId)
		if err != nil {
			return nil, fmt.Errorf("error while changing instance of book in database:%w", err)
		}
	}
	if method == "DELETE" {
		err := b.repo.DeleteListBook(listBookId)
		if err != nil {
			return nil, fmt.Errorf("error while deleting one instance of book from database:%w", err)
		}
		return nil, nil
	}
	return nil, nil
}
func InputCoverFoto(req *http.Request, input *IndTask.Book) error {
	reqFile, fileHeader, err := req.FormFile("file")
	if err != nil {
		logger.Errorf("InputCoverFoto: error while getting file from multipart form:%s", err)
		return fmt.Errorf("inputCoverFoto: error while getting file from multipart form:%w", err)
	}
	defer reqFile.Close()
	filePath := fmt.Sprintf("images/book_covers/%s_%d.%s", translate.Translate(input.BookName), input.Published, (strings.Split(fileHeader.Filename, "."))[1])
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		logger.Errorf("InputCoverFoto: error while opening file %s:%s", filePath, err)
		return fmt.Errorf("inputCoverFoto: error while opening file %s:%w", filePath, err)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(reqFile)
	if err != nil {
		logger.Errorf("InputCoverFoto: error while reading file from request:%s", err)
		return fmt.Errorf("inputCoverFoto: error while reading file from request:%w", err)
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		logger.Errorf("InputCoverFoto: error while writing file:%s", err)
		return fmt.Errorf("inputCoverFoto: error while writing file:%s", err)
	}
	input.Cover = filePath
	return nil
}

func (b *BookService) CalcRentCost(bookCost float32) float64 {
	profitability := b.cfg.ProfitBook.Profitability
	maxRentNumber := b.cfg.ProfitBook.MaxRentalNumber
	rentCost := float64(bookCost * profitability * 100 / maxRentNumber)
	return math.Round(rentCost) / 100
}
