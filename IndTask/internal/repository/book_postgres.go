package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"time"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

var bookLimit = 10

func (r *BookPostgres) GetThreeBooks() ([]IndTask.BookDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetThreeBooks: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getThreeBooks: can not starts transaction:%w", err)
	}
	var listBooks []IndTask.BookDTO
	var rows *sql.Rows
	query := "SELECT id, book_name, cost, cover, published, pages, amount FROM books JOIN " +
		"(SELECT book_id, SUM(rent_number) AS sum FROM list_books GROUP BY book_id ORDER BY sum DESC LIMIT 3) AS list  ON books.id = list.book_id"
	rows, err = transaction.Query(query)
	if err != nil {
		logger.Errorf("GetThreeBooks: can not executes a query:%s", err)
		return nil, fmt.Errorf("getThreeBooks: repository error:%w", err)
	}
	for rows.Next() {
		var book IndTask.BookDTO
		if err := rows.Scan(&book.Id, &book.BookName, &book.Cost, &book.Cover, &book.Published, &book.Pages, &book.Amount); err != nil {
			logger.Errorf("GetThreeBooks: error while scanning for book:%s", err)
			return nil, fmt.Errorf("getThreeBooks: repository error:%w", err)
		}
		book.Authors, err = r.ReturnBindAuthors(book.Id)
		if err != nil {
			return nil, fmt.Errorf("error while getting bound authors:%w", err)
		}
		book.Genre, err = r.ReturnBindGenres(book.Id)
		if err != nil {
			return nil, fmt.Errorf("error while getting bound genres:%w", err)
		}
		listBooks = append(listBooks, book)
	}
	return listBooks, transaction.Commit()
}

func (r *BookPostgres) GetBooks(page int) ([]IndTask.BookDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetBooks: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getBooks: can not starts transaction:%w", err)
	}
	var listBooks []IndTask.BookDTO
	var rows *sql.Rows
	if page == 0 {
		query := "SELECT id, book_name, cost, cover, published, pages, amount FROM books"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetBooks: can not executes a query:%s", err)
			return nil, fmt.Errorf("getBooks: repository error:%w", err)
		}
	} else {
		query := "SELECT id, book_name, cost, cover, published, pages, amount FROM books ORDER BY Id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, bookLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("GetBooks: can not executes a query:%s", err)
			return nil, fmt.Errorf("getBooks: repository error:%w", err)
		}
	}
	for rows.Next() {
		var book IndTask.BookDTO
		if err := rows.Scan(&book.Id, &book.BookName, &book.Cost, &book.Cover, &book.Published, &book.Pages, &book.Amount); err != nil {
			logger.Errorf("GetBooks: error while scanning for book:%s", err)
			return nil, fmt.Errorf("getBooks: repository error:%w", err)
		}
		book.Authors, err = r.ReturnBindAuthors(book.Id)
		if err != nil {
			return nil, fmt.Errorf("error while getting bound authors:%w", err)
		}
		book.Genre, err = r.ReturnBindGenres(book.Id)
		if err != nil {
			return nil, fmt.Errorf("error while getting bound genres:%w", err)
		}
		listBooks = append(listBooks, book)
	}
	return listBooks, transaction.Commit()
}

func (r *BookPostgres) CreateBook(book *IndTask.Book, bookExists bool, bookRentCost float64) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateBook: can not starts transaction:%s", err)
		return 0, fmt.Errorf("createBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var bookId int
	if bookExists {
		selectIdBookQuery := "SELECT id FROM books WHERE book_name=$1 AND published=$2"
		row := transaction.QueryRow(selectIdBookQuery, book.BookName, book.Published)
		if err := row.Scan(&bookId); err != nil {
			logger.Errorf("Error while scanning for bookId:%s", err)
			return 0, fmt.Errorf("createBook: error while scanning for bookId:%w", err)
		}
		for i := 0; i < book.Amount; i++ {
			createListBookQuery := "INSERT INTO list_books (book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped) VALUES ($1, $2, $3, $4, $5, $6, $7)"
			_, err := transaction.Exec(createListBookQuery, bookId, "false", 0, bookRentCost, time.Now(), 100, false)
			if err != nil {
				logger.Errorf("Error while execution query for insert into list_books:%s", err)
				return 0, fmt.Errorf("CreateBook: error while execution query for insert into list_books:%w", err)
			}
		}
		query := "UPDATE books SET amount=amount+$1 WHERE id = $2"
		_, err := transaction.Exec(query, book.Amount, bookId)
		if err != nil {
			logger.Errorf("Error while updating books.amount:%s", err)
			return 0, fmt.Errorf("createBook: Error while updating books.amount:%w", err)
		}
		return bookId, transaction.Commit()
	}
	createBookQuery := "INSERT INTO books (book_name, cost, cover, published, pages, amount) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	row := transaction.QueryRow(createBookQuery, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, book.Amount)
	if err := row.Scan(&bookId); err != nil {
		logger.Errorf("Error while scanning for bookId:%s", err)
		return 0, fmt.Errorf("createBook: error while scanning for bookId:%w", err)
	}
	for _, author := range book.AuthorsId {
		createBookAuthorQuery := "INSERT INTO book_author (book_id, author_id) VALUES ($1, $2)"
		_, err = transaction.Exec(createBookAuthorQuery, bookId, author)
		if err != nil {
			logger.Errorf("Error while execution query for insert into book_author:%s", err)
			return 0, fmt.Errorf("CreateBook: error while execution query for insert into book_author:%w", err)
		}
	}
	for _, genre := range book.GenreId {
		createBookAuthorQuery := "INSERT INTO book_genre (book_id, genre_id) VALUES ($1, $2)"
		_, err = transaction.Exec(createBookAuthorQuery, bookId, genre)
		if err != nil {
			logger.Errorf("Error while execution query for insert into book_genre:%s", err)
			return 0, fmt.Errorf("CreateBook: error while execution query for insert into book_genre:%w", err)
		}
	}
	for i := 0; i < book.Amount; i++ {
		createListBookQuery := "INSERT INTO list_books (book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err := transaction.Exec(createListBookQuery, bookId, "false", 0, bookRentCost, time.Now(), 100, false)
		if err != nil {
			logger.Errorf("Error while execution query for insert into list_books:%s", err)
			return 0, fmt.Errorf("CreateBook: error while execution query for insert into list_books:%w", err)
		}
	}
	return bookId, transaction.Commit()
}

func (r *BookPostgres) GetOneBook(bookId int) (*IndTask.BookDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneBook: can not starts transaction:%w", err)
	}
	var book IndTask.BookDTO
	query := "SELECT id, book_name, cost, cover, published, pages, amount FROM books WHERE id = $1"
	row := transaction.QueryRow(query, bookId)
	if err := row.Scan(&book.Id, &book.BookName, &book.Cost, &book.Cover, &book.Published, &book.Pages, &book.Amount); err != nil {
		logger.Errorf("GetOneBook: error while scanning for book:%s", err)
		return nil, fmt.Errorf("getOneBook: repository error:%w", err)
	}
	book.Authors, err = r.ReturnBindAuthors(book.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting bound authors:%w", err)
	}
	book.Genre, err = r.ReturnBindGenres(book.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting bound genres:%w", err)
	}
	return &book, transaction.Commit()
}

func (r *BookPostgres) ChangeBook(book *IndTask.Book, bookId int, bookRentCost float64) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeBook: can not starts transaction:%s", err)
		return fmt.Errorf("changeBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var query string
	if book.Cover == "" {
		query = "UPDATE books SET book_name=$1, cost=$2, published=$3, pages=$4 WHERE id = $5"
		_, err = transaction.Exec(query, book.BookName, book.Cost, book.Published, book.Pages, bookId)
	} else {
		query = "UPDATE books SET book_name=$1, cost=$2, cover=$3, published=$4, pages=$5 WHERE id = $6"
		_, err = transaction.Exec(query, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, bookId)
	}
	if err != nil {
		logger.Errorf("Repository error while updating book:%s", err)
		return fmt.Errorf("changeBook: repository error:%w", err)
	}
	query = "DELETE FROM book_author WHERE book_id = $1"
	_, err = transaction.Exec(query, bookId)
	if err != nil {
		logger.Errorf("Repository error while deleting bound author from book_author:%s", err)
		return fmt.Errorf("changeBook: repository error:%w", err)
	}
	query = "DELETE FROM book_genre WHERE book_id = $1"
	_, err = transaction.Exec(query, bookId)
	if err != nil {
		logger.Errorf("Repository error while deleting bound genre from book_genre:%s", err)
		return fmt.Errorf("changeBook: repository error:%w", err)
	}
	query = "UPDATE list_books SET rent_cost = $1 WHERE book_id=$2"
	_, err = transaction.Exec(query, bookRentCost, bookId)
	if err != nil {
		logger.Errorf("Repository error while updating rent_cost in list_books:%s", err)
		return fmt.Errorf("changeBook: repository error:%w", err)
	}
	for _, author := range book.AuthorsId {
		createBookAuthorQuery := "INSERT INTO book_author (book_id, author_id) VALUES ($1, $2)"
		_, err = transaction.Exec(createBookAuthorQuery, bookId, author)
		if err != nil {
			logger.Errorf("Repository error while insert into book_author:%s", err)
			return fmt.Errorf("changeBook: repository error:%w", err)
		}
	}
	for _, genre := range book.GenreId {
		createBookGenreQuery := "INSERT INTO book_genre (book_id, genre_id) VALUES ($1, $2)"
		_, err = transaction.Exec(createBookGenreQuery, bookId, genre)
		if err != nil {
			logger.Errorf("Repository error while insert into book_genre:%s", err)
			return fmt.Errorf("changeBook: repository error:%w", err)
		}
	}
	return transaction.Commit()
}

func (r *BookPostgres) DeleteBook(bookId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeBook: can not starts transaction:%s", err)
		return fmt.Errorf("changeBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM books WHERE id = $1"
	_, err = transaction.Exec(query, bookId)
	if err != nil {
		logger.Errorf("Repository error while deleting book:%s", err)
		return fmt.Errorf("deleteBook: repository error:%w", err)
	}
	query = "DELETE FROM list_books WHERE book_id = $1"
	_, err = transaction.Exec(query, bookId)
	if err != nil {
		logger.Errorf("Repository error while deleting listBook:%s", err)
		return fmt.Errorf("deleteBook: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *BookPostgres) GetListBooks(page int) ([]IndTask.ListBooksDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetListBooks: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getListBooks: can not starts transaction:%w", err)
	}
	var listBooks []IndTask.ListBooksDTO
	var rows *sql.Rows
	if page == 0 {
		query := "SELECT id, book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped FROM list_books WHERE issued='false' and scrapped='false'"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetListBooks: can not executes a query:%s", err)
			return nil, fmt.Errorf("getListBooks: repository error:%w", err)
		}
	} else {
		query := "SELECT id, book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped FROM list_books WHERE issued='false' and scrapped='false' ORDER BY Id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, bookLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("GetListBooks: can not executes a query:%s", err)
			return nil, fmt.Errorf("getListBooks: repository error:%w", err)
		}
	}
	for rows.Next() {
		var book IndTask.ListBooksDTO
		var bookId int
		if err := rows.Scan(&book.Id, &bookId, &book.Issued, &book.RentNumber, &book.RentCost, &book.RegDate, &book.Condition, &book.Scrapped); err != nil {
			logger.Errorf("GetListBooks: error while scanning for book:%s", err)
			return nil, fmt.Errorf("getListBooks: repository error:%w", err)
		}
		book.Book, err = r.GetOneBook(bookId)
		if err != nil {
			logger.Errorf("GetListBooks: error while getting book:%s", err)
			return nil, fmt.Errorf("getListBooks: repository error:%w", err)
		}
		listBooks = append(listBooks, book)
	}
	return listBooks, transaction.Commit()
}

func (r *BookPostgres) GetOneListBook(listBookId int) (*IndTask.ListBooksDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneListBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneListBook: can not starts transaction:%w", err)
	}
	var listBook IndTask.ListBooksDTO
	var bookId int
	query := "SELECT id, book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped FROM list_books WHERE id = $1"
	row := transaction.QueryRow(query, listBookId)
	if err := row.Scan(&listBook.Id, &bookId, &listBook.Issued, &listBook.RentNumber, &listBook.RentCost, &listBook.RegDate, &listBook.Condition, &listBook.Scrapped); err != nil {
		logger.Errorf("GetOneListBook: error while scanning for listBook:%s", err)
		return nil, fmt.Errorf("getOneListBook: repository error:%w", err)
	}
	listBook.Book, err = r.GetOneBook(bookId)
	if err != nil {
		logger.Errorf("GetOneListBook: error while getting book:%s", err)
		return nil, fmt.Errorf("getOneListBook: repository error:%w", err)
	}
	return &listBook, transaction.Commit()
}

func (r *BookPostgres) ChangeListBook(listBook *IndTask.ListBooks, listBookId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeListBook: can not starts transaction:%s", err)
		return fmt.Errorf("changeListBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "UPDATE list_books SET issued=$1, rent_number=$2, rent_cost=$3, reg_date=$4, condition=$5, scrapped=$6 WHERE id = $7"
	_, err = transaction.Exec(query, listBook.Issued, listBook.RentNumber, listBook.RentCost, listBook.RegDate, listBook.Condition, listBook.Scrapped, listBookId)
	if err != nil {
		logger.Errorf("Repository error while updating listBook:%s", err)
		return fmt.Errorf("changeListBook: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *BookPostgres) DeleteListBook(listBookId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("DeleteListBook: can not starts transaction:%s", err)
		return fmt.Errorf("deleteListBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM list_books WHERE id = $1"
	_, err = transaction.Exec(query, listBookId)
	if err != nil {
		logger.Errorf("Repository error while deleting listBook:%s", err)
		return fmt.Errorf("deleteListBook: repository error:%w", err)
	}
	query = "UPDATE books SET amount=amount-1 WHERE id = (SELECT book_id FROM list_books WHERE id = $1)"
	_, err = transaction.Exec(query, listBookId)
	if err != nil {
		logger.Errorf("Error while updating books.amount:%s", err)
		return fmt.Errorf("deleteListBook: Error while updating books.amount:%w", err)
	}
	return transaction.Commit()
}

func (r *BookPostgres) GetAuthorsByBookId(bookId int) ([]int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetAuthorsByBookId: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getAuthorsByBookId: can not starts transaction:%s", err)
	}
	var authorsId []int
	query := "SELECT author_id FROM book_author WHERE book_id=$1"
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		logger.Errorf("GetAuthorsByBookId: can not executes a query:%s", err)
		return nil, fmt.Errorf("getAuthorsByBookId: repository error:%w", err)
	}
	for rows.Next() {
		var authorId int
		if err := rows.Scan(&authorId); err != nil {
			logger.Errorf("GetAuthorsByBookId: error while scanning for authorId:%s", err)
			return nil, fmt.Errorf("getAuthorsByBookId: repository error:%w", err)
		}
		authorsId = append(authorsId, authorId)
	}
	return authorsId, transaction.Commit()
}

func (r *BookPostgres) ReturnBindAuthors(bookId int) ([]IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ReturnBindAuthors: can not starts transaction:%s", err)
		return nil, fmt.Errorf("returnBindAuthors: can not starts transaction:%s", err)
	}
	var authors []IndTask.Author
	query := "SELECT id, author_name, author_foto FROM authors JOIN book_author ON authors.id = book_author.author_id AND book_author.book_id = $1"
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		logger.Errorf("ReturnBindAuthors: can not executes a query:%s", err)
		return nil, fmt.Errorf("returnBindAuthors: repository error:%w", err)
	}
	for rows.Next() {
		var author IndTask.Author
		if err := rows.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("ReturnBindAuthors: error while scanning for author:%s", err)
			return nil, fmt.Errorf("returnBindAuthors: repository error:%w", err)
		}
		authors = append(authors, author)
	}
	return authors, transaction.Commit()
}

func (r *BookPostgres) ReturnBindGenres(bookId int) ([]IndTask.Genre, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ReturnBindGenres: can not starts transaction:%s", err)
		return nil, fmt.Errorf("returnBindGenres: can not starts transaction:%s", err)
	}
	var genres []IndTask.Genre
	query := "SELECT id, genre_name FROM genre JOIN book_genre ON genre.id = book_genre.genre_id AND book_genre.book_id = $1"
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		logger.Errorf("ReturnBindGenres: can not executes a query:%s", err)
		return nil, fmt.Errorf("returnBindGenres: repository error:%w", err)
	}
	for rows.Next() {
		var genre IndTask.Genre
		if err := rows.Scan(&genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("ReturnBindGenres: error while scanning for genre:%s", err)
			return nil, fmt.Errorf("returnBindGenres: repository error:%w", err)
		}
		genres = append(genres, genre)
	}
	return genres, transaction.Commit()
}
