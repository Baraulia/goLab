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
		logger.Errorf("GetGenres: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getGenres: can not starts transaction:%w", err)
	}
	var listGenres []IndTask.Genre
	query := "SELECT id, genre_name FROM genre"
	rows, err := transaction.Query(query)
	if err != nil {
		logger.Errorf("GetGenres: can not executes a query:%s", err)
		return nil, fmt.Errorf("getGenres:repository error:%w", err)
	}
	for rows.Next() {
		var genre IndTask.Genre
		if err := rows.Scan(&genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("Error while scanning for genre:%s", err)
			return nil, fmt.Errorf("getGenres:repository error:%w", err)
		}
		listGenres = append(listGenres, genre)
	}
	return listGenres, transaction.Commit()
}

func (r *GenrePostgres) CreateGenre(genre *IndTask.Genre) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateGenre: can not starts transaction:%s", err)
		return 0, fmt.Errorf("createGenre: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var genreId int
	createGenreQuery := "INSERT INTO genre (genre_name) VALUES ($1) RETURNING id"
	row := transaction.QueryRow(createGenreQuery, genre.GenreName)
	if err := row.Scan(&genreId); err != nil {
		logger.Errorf("Error while scanning for genreId:%s", err)
		return 0, fmt.Errorf("createGenre: error while scanning for genreId:%w", err)
	}
	return genreId, transaction.Commit()
}

func (r *GenrePostgres) GetOneGenre(genreId int) (*IndTask.Genre, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneGenre: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneGenre: can not starts transaction:%w", err)
	}
	var genre IndTask.Genre
	query := "SELECT id, genre_name FROM genre WHERE id = $1"
	row := transaction.QueryRow(query, genreId)
	if err := row.Scan(&genre.Id, &genre.GenreName); err != nil {
		logger.Errorf("Error while scanning for genre:%s", err)
		return nil, fmt.Errorf("getOneGenre: repository error:%w", err)
	}
	return &genre, transaction.Commit()
}

func (r *GenrePostgres) ChangeGenre(genre *IndTask.Genre, genreId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeGenre: can not starts transaction:%s", err)
		return fmt.Errorf("changeGenre: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "UPDATE genre SET genre_name=$1 WHERE id = $2"
	_, err = transaction.Exec(query, genre.GenreName, genreId)
	if err != nil {
		logger.Errorf("Repository error while updating genre:%s", err)
		return fmt.Errorf("changeGenre: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *GenrePostgres) DeleteGenre(genreId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("DeleteGenre: can not starts transaction:%s", err)
		return fmt.Errorf("deleteGenre: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM genre WHERE id=$1"
	_, err = transaction.Exec(query, genreId)
	if err != nil {
		logger.Errorf("Repository error while deleting genre:%s", err)
		return fmt.Errorf("deleteGenre: repository error:%w", err)
	}
	return transaction.Commit()
}
