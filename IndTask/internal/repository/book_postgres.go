package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (r *BookPostgres) CreateBook(book IndTask.Book) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createBookQuery := fmt.Sprintf("INSERT INTO %s (book_name, cost, cover, published, pages, amount) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", booksTable)
	row := transaction.QueryRow(createBookQuery, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, book.Amount)
	if err := row.Scan(&id); err != nil {
		transaction.Rollback()
		return 0, err
	}

	for author := range book.AuthorsId {
		createBookAuthorQuery := fmt.Sprintf("INSERT INTO %s (book_id, author_id) VALUES ($1, $2)", bookAuthorTable)
		_, err = transaction.Exec(createBookAuthorQuery, id, author)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}
	}

	for genre := range book.GenreId {
		createBookAuthorQuery := fmt.Sprintf("INSERT INTO %s (book_id, genre_id) VALUES ($1, $2)", bookGenreTable)
		_, err = transaction.Exec(createBookAuthorQuery, id, genre)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}
	}

	return id, transaction.Commit()
}

func (r *BookPostgres) GetBooks() []IndTask.Book {
	return nil
}

func (r *BookPostgres) ChangeBook() {

}
