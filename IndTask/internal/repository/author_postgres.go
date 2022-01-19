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

func (r *AuthorPostgres) GetAuthors(page int) ([]IndTask.Author, int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetAuthors: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getAuthors: can not starts transaction:%w", err)
	}
	var listAuthors []IndTask.Author
	var rows *sql.Rows
	var pages int
	if page == 0 {
		query := "SELECT id, author_name, author_foto FROM authors"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetAuthors: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getAuthors: repository error:%w", err)
		}
	} else {
		query := "SELECT id, author_name, author_foto FROM authors ORDER BY Id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, authorLimit, (page-1)*authorLimit)
		if err != nil {
			logger.Errorf("GetAuthors: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getAuthors: repository error:%w", err)
		}
	}

	for rows.Next() {
		var author IndTask.Author
		if err := rows.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("Error while scanning for author:%s", err)
			return nil, 0, fmt.Errorf("getAuthors: repository error:%w", err)
		}
		listAuthors = append(listAuthors, author)
	}
	query := "SELECT CEILING(COUNT(id)/$1::float) FROM authors"
	row := transaction.QueryRow(query, authorLimit)
	if err := row.Scan(&pages); err != nil {
		logger.Errorf("Error while scanning for pages:%s", err)
		return nil, 0, fmt.Errorf("getAuthors: error while scanning for pages:%w", err)
	}
	return listAuthors, pages, transaction.Commit()
}

func (r *AuthorPostgres) CreateAuthor(author *IndTask.Author) (*IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateAuthor: can not starts transaction:%s", err)
		return nil, fmt.Errorf("createAuthor: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var newAuthor IndTask.Author
	createAuthorQuery := "INSERT INTO authors (author_name, author_foto) VALUES ($1, $2) RETURNING id, author_name, author_foto"
	row := transaction.QueryRow(createAuthorQuery, author.AuthorName, author.AuthorFoto)
	if err := row.Scan(&newAuthor.Id, &newAuthor.AuthorName, &newAuthor.AuthorFoto); err != nil {
		logger.Errorf("Error while scanning for author:%s", err)
		return nil, fmt.Errorf("createAuthor: error while scanning for author:%w", err)
	}
	return &newAuthor, transaction.Commit()
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

func (r *AuthorPostgres) ChangeAuthor(author *IndTask.Author, authorId int) (*IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeAuthor: can not starts transaction:%s", err)
		return nil, fmt.Errorf("changeAuthor: can not starts transaction:%w", err)
	}
	var upAuthor IndTask.Author
	defer transaction.Rollback()
	query := "UPDATE authors SET author_name=$1, author_foto=$2 WHERE id=$3 RETURNING id, author_name, author_foto"
	row := transaction.QueryRow(query, author.AuthorName, author.AuthorFoto, authorId)
	if err := row.Scan(&upAuthor.Id, &upAuthor.AuthorName, &upAuthor.AuthorFoto); err != nil {
		logger.Errorf("Error while scanning for author:%s", err)
		return nil, fmt.Errorf("createAuthor: error while scanning for author:%w", err)
	}
	return &upAuthor, transaction.Commit()
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
