package repository

import (
	"database/sql"
	"fmt"
)

type ValidationPostgres struct {
	db *sql.DB
}

func NewValidationPostgres(db *sql.DB) *ValidationPostgres {
	return &ValidationPostgres{db: db}
}

func (v *ValidationPostgres) GetGenreById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("GetGenreById: can not starts transaction:%s", err)
		return fmt.Errorf("getGenreById: can not starts transaction:%w", err)
	}
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM genre WHERE id=$1)"
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Error while scanning for exist genre:%s", err)
		return fmt.Errorf("getGenreById: repository error:%w", err)
	}
	if !exist {
		return fmt.Errorf("such genre %d does not exist", id)
	}
	return transaction.Commit()
}
func (v *ValidationPostgres) GetAuthorById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("GetAuthorById: can not starts transaction:%s", err)
		return fmt.Errorf("getAuthorById: can not starts transaction:%w", err)
	}
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM authors WHERE id=$1)"
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Error while scanning for exist author:%s", err)
		return fmt.Errorf("getAuthorById: repository error:%w", err)
	}
	if !exist {
		return fmt.Errorf("such author %d does not exist", id)
	}
	return transaction.Commit()
}
func (v *ValidationPostgres) GetUserById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("GetUserById: can not starts transaction:%s", err)
		return fmt.Errorf("getUserById: can not starts transaction:%w", err)
	}
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)"
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Error while scanning for exist user:%s", err)
		return fmt.Errorf("getUserById: repository error:%w", err)
	}
	if !exist {
		return fmt.Errorf("such user %d does not exist", id)
	}
	return transaction.Commit()
}
func (v *ValidationPostgres) GetListBookById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("GetListBookById: can not starts transaction:%s", err)
		return fmt.Errorf("getListBookById: can not starts transaction:%w", err)
	}
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM list_books WHERE id=$1)"
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Error while scanning for exist instance of book:%s", err)
		return fmt.Errorf("getListBookById: repository error:%w", err)
	}
	if !exist {
		return fmt.Errorf("such instance of book %d does not exist", id)
	}
	return transaction.Commit()
}
func (v *ValidationPostgres) GetActById(id int, changing bool) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("GetActById: can not starts transaction:%s", err)
		return fmt.Errorf("getActById: can not starts transaction:%w", err)
	}
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM act WHERE id=$1)"
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Error while scanning for exist act:%s", err)
		return fmt.Errorf("getIssueActById: repository error:%w", err)
	}
	if !exist {
		return fmt.Errorf("such act %d does not exist", id)
	}
	if !changing {
		query = "SELECT EXISTS(SELECT 1 FROM act WHERE id=$1 AND status='open')"
		row = transaction.QueryRow(query, id)
		if err := row.Scan(&exist); err != nil {
			logger.Errorf("Error while scanning for exist act:%s", err)
			return fmt.Errorf("getIssueActById: repository error:%w", err)
		}
		if !exist {
			return fmt.Errorf("such act %d already closed", id)
		}
	}
	return transaction.Commit()
}
