package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
)

type GenrePostgres struct {
	db *sql.DB
}

func NewGenrePostgres(db *sql.DB) *GenrePostgres {
	return &GenrePostgres{db: db}
}

func (r *GenrePostgres) GetGenres() ([]IndTask.Genre, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not starts transaction:%s", err)
		return nil, err
	}
	defer transaction.Commit()

	var listGenres []IndTask.Genre
	query := fmt.Sprint("SELECT id, genre_name FROM genre")
	rows, err := transaction.Query(query)
	if err != nil {
		logger.Errorf("Can not executes a query:%s", err)
		return nil, err
	}

	for rows.Next() {
		var genre IndTask.Genre
		if err := rows.Scan(&genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		listGenres = append(listGenres, genre)
	}

	return listGenres, err
}

func (r *GenrePostgres) CreateGenre(genre *IndTask.Genre) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not starts transaction:%s", err)
		return 0, err
	}
	defer transaction.Rollback()

	var genreId int
	createGenreQuery := fmt.Sprint("INSERT INTO genre (genre_name) VALUES ($1) RETURNING id")
	row := transaction.QueryRow(createGenreQuery, genre.GenreName)
	if err := row.Scan(&genreId); err != nil {
		logger.Errorf("Scan error:%s", err)
		return 0, err
	}
	return genreId, transaction.Commit()

}

func (r *GenrePostgres) ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not starts transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var genre IndTask.Genre
		query := fmt.Sprint("SELECT id, genre_name FROM genre WHERE id = $1")
		row := transaction.QueryRow(query, genreId)
		if err := row.Scan(&genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		return &genre, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprintf("UPDATE genre SET genre_name=$1 WHERE id = $2")
		_, err := transaction.Exec(query, genre.GenreName, genreId)
		if err != nil {
			logger.Errorf("Update genre error:%s", err)
			return nil, err
		}
		return nil, transaction.Commit()
	}

	if method == "DELETE" {
		query := fmt.Sprint("DELETE FROM genre WHERE id=$1")
		_, err := transaction.Exec(query, genreId)
		if err != nil {
			logger.Errorf("Delete genre error:%s", err)
			return nil, err
		}
		return nil, transaction.Commit()
	}

	return nil, transaction.Rollback()
}
