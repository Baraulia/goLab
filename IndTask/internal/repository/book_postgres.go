package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"math"
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
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var listBooks []IndTask.BookDTO
	var rows *sql.Rows
	query := fmt.Sprint("SELECT id, book_name, cost, cover, published, pages, amount FROM books JOIN " +
		"(SELECT book_id, SUM(rent_number) AS sum FROM list_books GROUP BY book_id ORDER BY sum DESC LIMIT 3) AS list  ON books.id = list.book_id")
	rows, err = transaction.Query(query)
	if err != nil {
		logger.Errorf("Can not executes a query:%s", err)
		return nil, err
	}

	for rows.Next() {
		var book IndTask.BookDTO
		if err := rows.Scan(&book.Id, &book.BookName, &book.Cost, &book.Cover, &book.Published, &book.Pages, &book.Amount); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		book.Authors, err = r.ReturnBindAuthors(book.Id)
		if err != nil {
			logger.Errorf("Error returning autors binded with book_id=%d:%s", book.Id, err)
			return nil, err
		}
		book.Genre, err = r.ReturnBindGenres(book.Id)
		if err != nil {
			logger.Errorf("Error returning genres binded with book_id=%d:%s", book.Id, err)
			return nil, err
		}
		listBooks = append(listBooks, book)
	}
	return listBooks, err
}

func (r *BookPostgres) GetBooks(page int) ([]IndTask.BookDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var listBooks []IndTask.BookDTO
	var rows *sql.Rows
	if page == 0 {
		query := fmt.Sprint("SELECT id, book_name, cost, cover, published, pages, amount FROM books")
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	} else {
		query := fmt.Sprint("SELECT id, book_name, cost, cover, published, pages, amount FROM books ORDER BY Id LIMIT $1 OFFSET $2")
		rows, err = transaction.Query(query, bookLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	}
	for rows.Next() {
		var book IndTask.BookDTO
		if err := rows.Scan(&book.Id, &book.BookName, &book.Cost, &book.Cover, &book.Published, &book.Pages, &book.Amount); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		book.Authors, err = r.ReturnBindAuthors(book.Id)
		if err != nil {
			logger.Errorf("Error returning autors binded with book_id=%d:%s", book.Id, err)
			return nil, err
		}
		book.Genre, err = r.ReturnBindGenres(book.Id)
		if err != nil {
			logger.Errorf("Error returning genres binded with book_id=%d:%s", book.Id, err)
			return nil, err
		}
		listBooks = append(listBooks, book)
	}
	return listBooks, err
}

func (r *BookPostgres) CreateBook(book *IndTask.Book, bookExists bool) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return 0, err
	}
	defer transaction.Rollback()

	if bookExists {
		var bookId int
		selectIdBookQuery := fmt.Sprintf("SELECT id FROM books WHERE book_name=$1 AND published=$2")
		row := transaction.QueryRow(selectIdBookQuery, book.BookName, book.Published)
		if err := row.Scan(&bookId); err != nil {
			logger.Errorf("Scan error:%s", err)
			return 0, err
		}
		for i := 0; i < book.Amount; i++ {
			createListBookQuery := fmt.Sprint("INSERT INTO list_books (book_id, issued, rent_number, rent_cost, reg_date, condition) VALUES ($1, $2, $3, $4, $5, $6)")
			_, err := transaction.Exec(createListBookQuery, bookId, "false", 0, CalcRentCost(book), time.Now(), 100)
			if err != nil {
				logger.Errorf("Can not insert in the list_books:%s", err)
				return 0, err
			}
		}
		query := fmt.Sprint("UPDATE books SET  amount=amount+$1 WHERE id = $2")
		_, err := transaction.Exec(query, book.Amount, bookId)
		if err != nil {
			return 0, err
		}
		return bookId, transaction.Commit()
	}

	var id int
	createBookQuery := fmt.Sprint("INSERT INTO books (book_name, cost, cover, published, pages, amount) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	row := transaction.QueryRow(createBookQuery, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, book.Amount)
	if err := row.Scan(&id); err != nil {
		logger.Errorf("Scan error:%s", err)
		return 0, err
	}

	for _, author := range book.AuthorsId {
		createBookAuthorQuery := fmt.Sprint("INSERT INTO book_author (book_id, author_id) VALUES ($1, $2)")
		_, err = transaction.Exec(createBookAuthorQuery, id, author)
		if err != nil {
			logger.Errorf("Can not insert in the book_author:%s", err)
			return 0, err
		}
	}

	for _, genre := range book.GenreId {
		createBookAuthorQuery := fmt.Sprint("INSERT INTO book_genre (book_id, genre_id) VALUES ($1, $2)")
		_, err = transaction.Exec(createBookAuthorQuery, id, genre)
		if err != nil {
			logger.Errorf("Can not insert in the book_genre:%s", err)
			return 0, err
		}
	}

	for i := 0; i < book.Amount; i++ {
		createListBookQuery := fmt.Sprint("INSERT INTO list_books (book_id, issued, rent_number, rent_cost, reg_date, condition) VALUES ($1, $2, $3, $4, $5, $6)")
		_, err := transaction.Exec(createListBookQuery, id, "false", 0, CalcRentCost(book), time.Now(), 100)
		if err != nil {
			logger.Errorf("Can not insert in the list_books:%s", err)
			return 0, err
		}
	}
	return id, transaction.Commit()
}

func (r *BookPostgres) ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.BookDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var book IndTask.BookDTO
		query := fmt.Sprint("SELECT id, book_name, cost, cover, published, pages, amount FROM books WHERE id = $1")
		row := transaction.QueryRow(query, bookId)
		if err := row.Scan(&book.Id, &book.BookName, &book.Cost, &book.Cover, &book.Published, &book.Pages, &book.Amount); err != nil {
			logger.Errorf("Can not scan select from books where id = %d", bookId)
			return nil, err
		}
		book.Authors, err = r.ReturnBindAuthors(book.Id)
		if err != nil {
			logger.Errorf("Error returning autors binded with book_id=%d:%s", book.Id, err)
			return nil, err
		}
		book.Genre, err = r.ReturnBindGenres(book.Id)
		if err != nil {
			logger.Errorf("Error returning genres binded with book_id=%d:%s", book.Id, err)
			return nil, err
		}

		return &book, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE books SET book_name=$1, cost=$2, cover=$3, published=$4, pages=$5, amount=$6 WHERE id = $7")
		_, err := transaction.Exec(query, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, book.Amount, bookId)
		if err != nil {
			logger.Errorf("Can not update books where id = %d", bookId)
			return nil, err
		}
		query = fmt.Sprint("DELETE FROM book_author WHERE book_id = $1")
		_, err = transaction.Exec(query, bookId)
		if err != nil {
			logger.Errorf("Can not delete related entry from book_author where book_id = %d", bookId)
			return nil, err
		}
		query = fmt.Sprint("DELETE FROM book_genre WHERE book_id = $1")
		_, err = transaction.Exec(query, bookId)
		if err != nil {
			logger.Errorf("Can not delete related entry from book_genre where book_id = %d", bookId)
			return nil, err
		}

		for _, author := range book.AuthorsId {
			createBookAuthorQuery := fmt.Sprint("INSERT INTO book_author (book_id, author_id) VALUES ($1, $2)")
			_, err = transaction.Exec(createBookAuthorQuery, bookId, author)
			if err != nil {
				logger.Errorf("Can not insert in the book_author for book_id = %d, author_id = %d : %s", bookId, author, err)
				return nil, err
			}
		}

		for _, genre := range book.GenreId {
			createBookGenreQuery := fmt.Sprint("INSERT INTO book_genre (book_id, genre_id) VALUES ($1, $2)")
			_, err = transaction.Exec(createBookGenreQuery, bookId, genre)
			if err != nil {
				logger.Errorf("Can not insert in the book_genre for book_id = %d, genre_id = %d : %s", bookId, genre, err)
				return nil, err
			}
		}
		return nil, transaction.Commit()
	}

	if method == "DELETE" {

		query := fmt.Sprint("DELETE FROM books WHERE id = $1")
		_, err := transaction.Exec(query, bookId)
		if err != nil {
			return nil, err
		}
		return nil, transaction.Commit()
	}

	return nil, transaction.Rollback()
}

func (r *BookPostgres) GetListBooks(page int) ([]IndTask.ListBooksDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	var listBooks []IndTask.ListBooksDTO
	var rows *sql.Rows
	if page == 0 {
		query := fmt.Sprint("SELECT id, book_id, issued, rent_number, rent_cost, reg_date, condition FROM list_books WHERE issued='false'")
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	} else {
		query := fmt.Sprint("SELECT id, book_id, issued, rent_number, rent_cost, reg_date, condition FROM list_books WHERE issued='false' ORDER BY Id LIMIT $1 OFFSET $2")
		rows, err = transaction.Query(query, bookLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	}
	for rows.Next() {
		var book IndTask.ListBooksDTO
		var bookId int
		if err := rows.Scan(&book.Id, &bookId, &book.Issued, &book.RentNumber, &book.RentCost, &book.RegDate, &book.Condition); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		book.Book, err = r.ChangeBook(nil, bookId, "GET")
		if err != nil {
			logger.Errorf("Error get book with id = %s:%s", bookId, err)
			return nil, err
		}
		listBooks = append(listBooks, book)
	}
	return listBooks, err
}

func (r *BookPostgres) ChangeListBook(listBook *IndTask.ListBooks, listBookId int, method string) (*IndTask.ListBooksDTO, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var listBook IndTask.ListBooksDTO
		var bookId int
		query := fmt.Sprint("SELECT id, book_id, issued, rent_number, rent_cost, reg_date, condition FROM list_books WHERE id = $1")
		row := transaction.QueryRow(query, listBookId)
		if err := row.Scan(&listBook.Id, &bookId, &listBook.Issued, &listBook.RentNumber, &listBook.RentCost, &listBook.RegDate, &listBook.Condition); err != nil {
			return nil, err
		}
		listBook.Book, err = r.ChangeBook(nil, bookId, "GET")
		if err != nil {
			logger.Errorf("Error get book with id = %s:%s", bookId, err)
			return nil, err
		}
		return &listBook, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE list_books SET book_id=$1, issued=$2, rent_number=$3, rent_cost=$4, reg_date=$5, condition=$6 WHERE id = $7")
		_, err := transaction.Exec(query, listBook.BookId, listBook.Issued, listBook.RentNumber, listBook.RentCost, listBook.RegDate, listBook.Condition)
		if err != nil {
			return nil, err
		}
		return nil, transaction.Commit()
	}

	if method == "DELETE" {

		query := fmt.Sprint("DELETE FROM list_books WHERE id = $1")
		_, err := transaction.Exec(query, listBookId)
		if err != nil {
			return nil, err
		}
		return nil, transaction.Commit()
	}

	return nil, transaction.Rollback()
}

func CalcRentCost(book *IndTask.Book) float64 {
	rentCost := float64(book.Cost * 1.15 * 10)
	return math.Round(rentCost) / 100
}

func (r *BookPostgres) GetAuthorsByBookId(bookId int) ([]int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	var authorsId []int
	query := fmt.Sprint("SELECT author_id FROM book_author WHERE book_id=$1")
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		logger.Errorf("Can not executes a query:%s", err)
		return nil, err
	}

	for rows.Next() {
		var authorId int
		if err := rows.Scan(&authorId); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		authorsId = append(authorsId, authorId)
	}
	return authorsId, transaction.Commit()
}

func (r *BookPostgres) ReturnBindAuthors(bookId int) ([]IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var authors []IndTask.Author

	query := fmt.Sprint("SELECT id, author_name, author_foto FROM authors JOIN book_author ON authors.id = book_author.author_id AND book_author.book_id = $1")
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var author IndTask.Author
		if err := rows.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (r *BookPostgres) ReturnBindGenres(bookId int) ([]IndTask.Genre, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var genres []IndTask.Genre
	query := fmt.Sprint("SELECT id, genre_name FROM genre JOIN book_genre ON genre.id = book_genre.genre_id AND book_genre.book_id = $1")
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var genre IndTask.Genre
		if err := rows.Scan(&genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}
