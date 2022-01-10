package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
)

type AuthorPostgres struct {
	db *sql.DB
}

var authorLimit = 10

func NewAuthorPostgres(db *sql.DB) *AuthorPostgres {
	return &AuthorPostgres{db: db}
}

func (r *AuthorPostgres) GetAuthors(page int) ([]IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	var listAuthors []IndTask.Author
	var rows *sql.Rows
	if page == 0 {
		query := fmt.Sprint("SELECT id, author_name, author_foto FROM authors")
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	} else {
		query := fmt.Sprint("SELECT id, author_name, author_foto FROM authors ORDER BY Id LIMIT $1 OFFSET $2")
		rows, err = transaction.Query(query, authorLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	}

	for rows.Next() {
		var author IndTask.Author
		if err := rows.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		listAuthors = append(listAuthors, author)
	}
	return listAuthors, err
}

func (r *AuthorPostgres) CreateAuthor(author *IndTask.Author) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return 0, err
	}
	defer transaction.Rollback()

	var authorId int
	createAuthorQuery := fmt.Sprint("INSERT INTO authors (author_name, author_foto) VALUES ($1, $2) RETURNING id")
	row := transaction.QueryRow(createAuthorQuery, author.AuthorName, author.AuthorFoto)
	if err := row.Scan(&authorId); err != nil {
		logger.Errorf("Scan error:%s", err)
		return 0, err
	}
	return authorId, transaction.Commit()

}

func (r *AuthorPostgres) ChangeAuthor(author *IndTask.Author, authorId int, method string) (*IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var author IndTask.Author
		query := fmt.Sprint("SELECT id, author_name, author_foto FROM authors WHERE id =$1")
		row := transaction.QueryRow(query, authorId)
		if err := row.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		return &author, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE authors SET author_name=$1, author_foto=$2 WHERE id=$3")
		_, err := transaction.Exec(query, author.AuthorName, author.AuthorFoto, authorId)
		if err != nil {
			logger.Errorf("Update author error:%s", err)
			return nil, err
		}
		return nil, transaction.Commit()
	}

	if method == "DELETE" {
		query := fmt.Sprint("DELETE FROM authors WHERE id = $1")
		_, err := transaction.Exec(query, authorId)
		if err != nil {
			logger.Errorf("Delete author error:%s", err)
			return nil, err
		}
		return nil, transaction.Commit()
	}

	return nil, transaction.Rollback()
}
