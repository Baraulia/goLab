package repository

import "database/sql"

type GenrePostgres struct {
	db *sql.DB
}

func NewGenrePostgres(db *sql.DB) *GenrePostgres {
	return &GenrePostgres{db: db}
}

func (r *GenrePostgres) GetGenres() {

}

func (r *GenrePostgres) CreateGenre() {

}

func (r *GenrePostgres) ChangeGenre() {

}
