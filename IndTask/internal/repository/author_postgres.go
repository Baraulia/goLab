package repository

import "database/sql"

type AuthorPostgres struct {
	db *sql.DB
}

func NewAuthorPostgres(db *sql.DB) *AuthorPostgres {
	return &AuthorPostgres{db: db}
}

func (r *AuthorPostgres) GetAuthors() {

}

func (r *AuthorPostgres) CreateAuthor() {

}

func (r *AuthorPostgres) ChangeAuthor() {

}
