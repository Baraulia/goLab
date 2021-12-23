package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/sirupsen/logrus"
)

type GenrePostgres struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewGenrePostgres(db *sql.DB, logger *logging.Logger) *GenrePostgres {
	return &GenrePostgres{db: db, logger: logger}
}

func (r *GenrePostgres) GetGenres() ([]IndTask.Genre, error) {
	fmt.Println("works GetGenres repository")
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var listGenres []IndTask.Genre
	query := fmt.Sprint("SELECT * FROM genre")
	rows, err := transaction.Query(query)

	for rows.Next() {
		var genre IndTask.Genre
		if err := rows.Scan(&genre.Id, &genre.GenreName); err != nil {
			logrus.Fatal(err)
		}
		listGenres = append(listGenres, genre)
	}

	return listGenres, err

}

func (r *GenrePostgres) CreateGenre(genre *IndTask.Genre) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var genreId int
	createGenreQuery := fmt.Sprint("INSERT INTO genre (genre_name) VALUES ($1) RETURNING id")
	row := transaction.QueryRow(createGenreQuery, genre.GenreName)
	if err := row.Scan(&genreId); err != nil {
		transaction.Rollback()
		return 0, err
	}
	return genreId, transaction.Commit()

}

func (r *GenrePostgres) ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	if method == "GET" {
		var genre *IndTask.Genre
		query := fmt.Sprint("SELECT * FROM genre WHERE id = $1")

		row := transaction.QueryRow(query, genreId)
		if err := row.Scan(genre.Id, genre.GenreName); err != nil {
			transaction.Rollback()
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		return nil, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprintf("UPDATE genre SET genre_name = $1 WHERE id = $2")
		_, err := transaction.Exec(query, genre.GenreName, genreId)
		if err != nil {
			return nil, err
		}
		return nil, transaction.Commit()
	}

	if method == "DELETE" {
		query := fmt.Sprint("DELETE FROM genre WHERE id = $1")
		_, err := transaction.Exec(query, genreId)
		if err != nil {
			return nil, err
		}
		return nil, transaction.Commit()
	}
	return nil, nil

}
