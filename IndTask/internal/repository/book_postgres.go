package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/sirupsen/logrus"
	"time"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (r *BookPostgres) CreateBook(book *IndTask.Book) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return 0, err
	}

	var id int
	createBookQuery := fmt.Sprint("INSERT INTO books (book_name, cost, cover, published, pages, amount, rent_cost) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	row := transaction.QueryRow(createBookQuery, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, book.Amount, book.RentCost)
	if err := row.Scan(&id); err != nil {
		transaction.Rollback()
		return 0, err
	}

	for author := range book.AuthorsId {
		createBookAuthorQuery := fmt.Sprint("INSERT INTO book_author (book_id, author_id) VALUES ($1, $2)")
		_, err = transaction.Exec(createBookAuthorQuery, id, author)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}
	}

	for genre := range book.GenreId {
		createBookAuthorQuery := fmt.Sprint("INSERT INTO book_genre (book_id, genre_id) VALUES ($1, $2)")
		_, err = transaction.Exec(createBookAuthorQuery, id, genre)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}
	}

	for i := 0; i < book.Amount; i++ {
		createListBookQuery := fmt.Sprint("INSERT INTO list_books (book_id, issued, rent_number, rent_cost, reg_date, condition) VALUES ($1, $2, $3, $4, $5, $6)")
		_, err := transaction.Exec(createListBookQuery, id, "false", 0, book.RentCost, time.Now(), 100)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}
	}

	return id, transaction.Commit()
}

func (r *BookPostgres) GetBooks() ([]IndTask.Book, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	var listBooks []IndTask.Book
	query := fmt.Sprint("SELECT * FROM books")

	rows, err := transaction.Query(query)

	for rows.Next() {
		var book IndTask.Book
		if err := rows.Scan(&book); err != nil {
			logrus.Fatal(err)
		}
		listBooks = append(listBooks, book)
	}

	return listBooks, err
}

func (r *BookPostgres) ChangeBook(book *IndTask.Book, bookId int, method string) (*IndTask.Book, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var book *IndTask.Book
		query := fmt.Sprint("SELECT * FROM books WHERE id = $1")

		row := transaction.QueryRow(query, bookId)
		if err := row.Scan(&book); err != nil {
			return nil, err
		}

		return book, nil
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE books SET (book_name=$1, cost=$2, cover=$3, published=$4, pages=$5, amount=$6, rent_cost=$7) WHERE id = $8")
		_, err := transaction.Exec(query, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, book.Amount, book.RentCost, bookId)
		return nil, err
	}

	if method == "DELETE" {
		query := fmt.Sprint("DELETE FROM books WHERE id = $1")
		_, err := transaction.Exec(query, bookId)
		return nil, err
	}

	return nil, transaction.Commit()
}
