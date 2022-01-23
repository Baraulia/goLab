package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/myErrors"
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

func (b *BookService) GetThreeBooks() ([]IndTask.MostPopularBook, error) {
	books, err := b.repo.GetThreeBooks()
	if err != nil {
		return nil, fmt.Errorf("error while getting three books from database:%w", err)
	}
	return books, nil
}

func (b *BookService) GetBooks(page int, sorting string) ([]*IndTask.BookResponse, int, error) {
	books, pages, err := b.repo.GetBooks(page, sorting)
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting books from database:%w", err)
	}
	return books, pages, nil
}

func (b *BookService) GetListBooks(page int) ([]IndTask.ListBooksResponse, int, error) {
	books, pages, err := b.repo.GetListBooks(page)
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting instances of books from database:%w", err)
	}
	return books, pages, nil
}

func (b *BookService) CreateBook(book *IndTask.Book) (*IndTask.OneBookResponse, error) {
	listBooks, _, err := b.repo.GetBooks(0, bookNameDesc)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting books from database:%w", err), Code: 500}
	}
	bookExists := false
	for _, bdBook := range listBooks {
		if bdBook.BookName == book.BookName {
			if bdBook.Published == book.Published {
				authorsId, err := b.repo.GetAuthorsByBookId(bdBook.Id)
				if err != nil {
					logger.Errorf("Error when getting authorsId from bookId = %d:%s", bdBook.Id, err)
					return nil, &myErrors.MyError{Err: fmt.Errorf("error when getting authorsId from bookId = %d:%w", bdBook.Id, err), Code: 500}
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
				return nil, &myErrors.MyError{Err: fmt.Errorf("BookExists = true, error deleting file %s:%s", book.Cover, err), Code: 500}
			}
		}
		return nil, &myErrors.MyError{Err: fmt.Errorf("BookExists = true for book_name = %s, book_published = %d", book.BookName, book.Published), Code: 400}
	}
	bookRentCost := b.CalcRentCost(book.Cost)
	newBook, err := b.repo.CreateBook(book, bookRentCost)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while creating book in database:%w", err), Code: 500}
	}
	return newBook, nil
}

func (b *BookService) ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.OneBookResponse, error) {
	err := b.repo.Validation.GetBookById(bookId)
	if err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			logger.Errorf("Such a book:%d does not exist", bookId)
			return nil, &myErrors.MyError{Err: fmt.Errorf("such a book:%d does not exist", bookId), Code: e.Status()}
		default:
			logger.Errorf("changeBook:%s", err)
			return nil, fmt.Errorf("changeBook:%w", err)
		}
	}
	if method == "GET" {
		oneBook, err := b.repo.GetOneBook(bookId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting one book from database:%w", err), Code: 500}
		}
		return oneBook, nil
	}
	if method == "PUT" {
		bookRentCost := b.CalcRentCost(book.Cost)
		upBook, err := b.repo.ChangeBook(book, bookId, bookRentCost)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while changing book in database:%w", err), Code: 500}
		}
		return upBook, nil
	}
	if method == "DELETE" {
		err := b.repo.DeleteBook(bookId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while deleting one book from database:%w", err), Code: 500}
		}
		return nil, nil
	}
	return nil, nil
}

func (b *BookService) CreateListBook(bookId int) (*IndTask.ListBooksResponse, error) {
	err := b.repo.Validation.GetBookById(bookId)
	if err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			logger.Errorf("Such book:%d does not exist", bookId)
			return nil, &myErrors.MyError{Err: fmt.Errorf("createListBook:%w", err), Code: e.Status()}
		default:
			logger.Errorf("createListBook:%s", err)
			return nil, fmt.Errorf("createListBook:%w", err)
		}
	}
	book, err := b.repo.AppBook.GetOneBook(bookId)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting book from database:%w", err), Code: 500}
	}
	bookRentCost := b.CalcRentCost(book.Cost)
	listBook, err := b.repo.AppBook.CreateListBook(bookId, bookRentCost)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while creating listBook in database:%w", err), Code: 500}
	}
	return listBook, nil
}

func (b *BookService) ChangeListBook(listBook *IndTask.ListBook, listBookId int, method string) (*IndTask.ListBooksResponse, error) {
	err := b.repo.Validation.GetListBookById(listBookId)
	if err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			logger.Errorf("Such a instance of book:%d does not exist", listBookId)
			return nil, &myErrors.MyError{Err: fmt.Errorf("such a instance of book:%d does not exist", listBookId), Code: e.Status()}
		default:
			logger.Errorf("changeListBook:%s", err)
			return nil, fmt.Errorf("changeListBook:%w", err)
		}
	}
	if method == "GET" {
		book, err := b.repo.GetOneListBook(listBookId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting one instance of book from database:%w", err), Code: 500}
		}
		return book, nil
	}
	if method == "PUT" {
		upBook, err := b.repo.ChangeListBook(listBook, listBookId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while changing instance of book in database:%w", err), Code: 500}
		}
		return upBook, nil
	}
	if method == "DELETE" {
		err := b.repo.DeleteListBook(listBookId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while deleting one instance of book from database:%w", err), Code: 500}
		}
		return nil, nil
	}
	return nil, nil
}

func (b *BookService) InputCoverFoto(req *http.Request, input *IndTask.Book) error {
	reqFile, fileHeader, err := req.FormFile("file")
	if err != nil {
		logger.Errorf("InputCoverFoto: error while getting file from multipart form:%s", err)
		return &myErrors.MyError{Err: fmt.Errorf("inputCoverFoto: error while getting file from multipart form:%w", err), Code: 400}
	}
	defer reqFile.Close()
	filePath := fmt.Sprintf("%simages/book_covers/%s_%d.%s", b.cfg.FilePath, translate.Translate(input.BookName), input.Published, (strings.Split(fileHeader.Filename, "."))[1])
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		logger.Errorf("InputCoverFoto: error while opening file %s:%s", filePath, err)
		return &myErrors.MyError{Err: fmt.Errorf("inputCoverFoto: error while opening file %s:%w", filePath, err), Code: 500}
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(reqFile)
	if err != nil {
		logger.Errorf("InputCoverFoto: error while reading file from request:%s", err)
		return &myErrors.MyError{Err: fmt.Errorf("inputCoverFoto: error while reading file from request:%w", err), Code: 500}
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		logger.Errorf("InputCoverFoto: error while writing file:%s", err)
		return &myErrors.MyError{Err: fmt.Errorf("inputCoverFoto: error while writing file:%s", err), Code: 500}
	}
	filePath = strings.Replace(filePath, b.cfg.FilePath, "", 1)
	input.Cover = filePath
	return nil
}

func (b *BookService) CalcRentCost(bookCost float32) float64 {
	profitability := b.cfg.ProfitBook.Profitability
	maxRentNumber := b.cfg.ProfitBook.MaxRentalNumber
	rentCost := float64(bookCost * profitability * 100 / maxRentNumber)
	return math.Round(rentCost) / 100
}

func (b *BookService) SortTypeBook(sorting string) string {
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
