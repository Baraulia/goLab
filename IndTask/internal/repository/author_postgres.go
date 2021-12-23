package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/sirupsen/logrus"
)

type AuthorPostgres struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewAuthorPostgres(db *sql.DB, logger *logging.Logger) *AuthorPostgres {
	return &AuthorPostgres{db: db, logger: logger}
}

func (r *AuthorPostgres) GetAuthors() ([]IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var listAuthors []IndTask.Author
	query := fmt.Sprint("SELECT * FROM authors")

	rows, err := transaction.Query(query)

	for rows.Next() {
		var author IndTask.Author
		if err := rows.Scan(&author); err != nil {
			logrus.Fatal(err)
		}
		listAuthors = append(listAuthors, author)
	}

	return listAuthors, err
}

func (r *AuthorPostgres) CreateAuthor(author *IndTask.Author) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var authorId int
	createAuthorQuery := fmt.Sprint("INSERT INTO authors (author_name, author_foto) VALUES ($1, $2) RETURNING id")
	row := transaction.QueryRow(createAuthorQuery, author.AuthorName, author.AuthorFoto)
	if err := row.Scan(&authorId); err != nil {
		transaction.Rollback()
		return 0, err
	}
	return authorId, transaction.Commit()

}

func (r *AuthorPostgres) ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	if method == "GET" {
		var author *IndTask.Author
		query := fmt.Sprint("SELECT * FROM authors WHERE id =$1")

		row := transaction.QueryRow(query, authorId)
		if err := row.Scan(&author); err != nil {
			transaction.Rollback()
			return nil, err
		}

		return author, nil
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE authors SET (author_name=$1, author_foto=$2) WHERE id = $3")
		_, err := transaction.Exec(query, author.AuthorName, author.AuthorFoto, authorId)
		return nil, err
	}
	if method == "DELETE" {
		query := fmt.Sprint("DELETE FROM authors WHERE id = $1")
		_, err := transaction.Exec(query, authorId)
		return nil, err
	}
	return nil, transaction.Commit()

}
