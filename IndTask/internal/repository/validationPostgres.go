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
		logger.Errorf("Can not begin transaction:%s", err)
		return err
	}
	var exist bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM genre WHERE id=$1)")
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Scan error:%s", err)
		return err
	}
	if !exist {
		return fmt.Errorf("does not exist")
	}
	return nil
}
func (v *ValidationPostgres) GetAuthorById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return err
	}
	var exist bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM authors WHERE id=$1)")
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Scan error:%s", err)
		return err
	}
	if !exist {
		return fmt.Errorf("does not exist")
	}
	return nil
}
func (v *ValidationPostgres) GetUserById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return err
	}
	var exist bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)")
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Scan error:%s", err)
		return err
	}
	if !exist {
		return fmt.Errorf("does not exist")
	}
	return nil
}
func (v *ValidationPostgres) GetListBookById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return err
	}
	var exist bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM list_books WHERE id=$1)")
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Scan error:%s", err)
		return err
	}
	if !exist {
		return fmt.Errorf("does not exist")
	}
	return nil
}
func (v *ValidationPostgres) GetIssueActById(id int) error {
	transaction, err := v.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return err
	}
	var exist bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM issue_act WHERE id=$1)")
	row := transaction.QueryRow(query, id)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Scan error:%s", err)
		return err
	}
	if !exist {
		return fmt.Errorf("does not exist")
	}
	return nil
}
