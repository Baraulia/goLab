package service

import (
	"errors"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"os"
)

var BookDoesNotExists = errors.New("book with that id does not exists")
var ListBookDoesNotExists = errors.New("listBook with that id does not exists")

type BookService struct {
	repo repository.Repository
}

func NewBookService(repo repository.Repository) *BookService {
	return &BookService{repo: repo}
}

func (b *BookService) GetThreeBooks() ([]IndTask.BookDTO, error) {
	return b.repo.GetThreeBooks()
}

func (b *BookService) GetBooks(page int) ([]IndTask.BookDTO, error) {
	return b.repo.GetBooks(page)
}

func (b *BookService) GetListBooks(page int) ([]IndTask.ListBooksDTO, error) {
	return b.repo.GetListBooks(page)
}

func (b *BookService) CreateBook(book *IndTask.Book) (int, error) {
	listBooks, err := b.repo.GetBooks(0)
	if err != nil {
		logger.Errorf("Error when getting books:%s", err.Error())
		return 0, err
	}

	bookExists := false
	for _, bdBook := range listBooks {
		if bdBook.BookName == book.BookName {
			if bdBook.Published == book.Published {
				authorsId, err := b.repo.GetAuthorsByBookId(bdBook.Id)
				if err != nil {
					logger.Errorf("Error when getting authorsId from bookId = %d :%s", bdBook.Id, err.Error())
					return 0, err
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
			break
		}
	}
	if bookExists {
		logger.Errorf("BookExists = true for book_name = %s, book_published = %d", book.BookName, book.Published)
		if err := os.Remove(book.Cover); err != nil {
			return 0, fmt.Errorf("BookExists = true, error deleting file %s:%s", book.Cover, err)
		}

		return 0, fmt.Errorf("BookExists = true for book_name = %s, book_published = %d", book.BookName, book.Published)
	}
	return b.repo.CreateBook(book, bookExists)
}

func (b *BookService) ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.BookDTO, error) {
	listBooks, err := b.repo.GetBooks(0)
	if err != nil {
		logger.Errorf("Error when getting books:%s", err)
		return nil, err
	}

	var bookExist = false
	for _, bdBook := range listBooks {
		if bdBook.Id == bookId {
			bookExist = true
		}
	}
	if bookExist == false {
		logger.Error("Such a book does not exist")
		return nil, BookDoesNotExists
	}
	return b.repo.ChangeBook(book, bookId, method)
}

func (b *BookService) ChangeListBook(listBook *IndTask.ListBooks, listBookId int, method string) (*IndTask.ListBooksDTO, error) {
	listBooks, err := b.repo.GetListBooks(0)
	if err != nil {
		logger.Errorf("Error when getting listBooks:%s", err)
		return nil, err
	}

	var bookListExist = false
	for _, listBook := range listBooks {
		if listBook.Id == listBookId {
			bookListExist = true
		}
	}
	if bookListExist == false {
		logger.Error("Such a listBook does not exist")
		return nil, ListBookDoesNotExists
	}

	return b.repo.ChangeListBook(listBook, listBookId, method)
}
