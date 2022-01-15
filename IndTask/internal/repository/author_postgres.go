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
		logger.Errorf("GetAuthors: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getAuthors: can not starts transaction:%w", err)
	}
	var listAuthors []IndTask.Author
	var rows *sql.Rows
	if page == 0 {
		query := "SELECT id, author_name, author_foto FROM authors"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetAuthors: can not executes a query:%s", err)
			return nil, fmt.Errorf("getAuthors: repository error:%w", err)
		}
	} else {
		query := "SELECT id, author_name, author_foto FROM authors ORDER BY Id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, authorLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("GetAuthors: can not executes a query:%s", err)
			return nil, fmt.Errorf("getAuthors: repository error:%w", err)
		}
	}

	for rows.Next() {
		var author IndTask.Author
		if err := rows.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("Error while scanning for author:%s", err)
			return nil, fmt.Errorf("getAuthors: repository error:%w", err)
		}
		listAuthors = append(listAuthors, author)
	}
	return listAuthors, transaction.Commit()
}

func (r *AuthorPostgres) CreateAuthor(author *IndTask.Author) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateAuthor: can not starts transaction:%s", err)
		return 0, fmt.Errorf("createAuthor: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var authorId int
	createAuthorQuery := "INSERT INTO authors (author_name, author_foto) VALUES ($1, $2) RETURNING id"
	row := transaction.QueryRow(createAuthorQuery, author.AuthorName, author.AuthorFoto)
	if err := row.Scan(&authorId); err != nil {
		logger.Errorf("Error while scanning for author:%s", err)
		return 0, fmt.Errorf("createAuthor: error while scanning for author:%w", err)
	}
	return authorId, transaction.Commit()
}

func (r *AuthorPostgres) GetOneAuthor(authorId int) (*IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneAuthor: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneAuthor: can not starts transaction:%w", err)
	}
	var author IndTask.Author
	query := "SELECT id, author_name, author_foto FROM authors WHERE id =$1"
	row := transaction.QueryRow(query, authorId)
	if err := row.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
		logger.Errorf("Error while scanning for author:%s", err)
		return nil, fmt.Errorf("getOneAuthor: repository error:%w", err)
	}
	return &author, transaction.Commit()
}

func (r *AuthorPostgres) ChangeAuthor(author *IndTask.Author, authorId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeAuthor: can not starts transaction:%s", err)
		return fmt.Errorf("changeAuthor: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "UPDATE authors SET author_name=$1, author_foto=$2 WHERE id=$3"
	_, err = transaction.Exec(query, author.AuthorName, author.AuthorFoto, authorId)
	if err != nil {
		logger.Errorf("Update author error:%s", err)
		return fmt.Errorf("changeAuthor: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *AuthorPostgres) DeleteAuthor(authorId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("DeleteAuthor: can not starts transaction:%s", err)
		return fmt.Errorf("deleteAuthor: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM authors WHERE id = $1"
	_, err = transaction.Exec(query, authorId)
	if err != nil {
		logger.Errorf("Repository error while deleting author:%s", err)
		return fmt.Errorf("deleteAuthor: repository error:%w", err)
	}
	return transaction.Commit()
}
