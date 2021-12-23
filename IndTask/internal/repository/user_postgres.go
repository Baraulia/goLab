package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/sirupsen/logrus"
)

type UserPostgres struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewUserPostgres(db *sql.DB, logger *logging.Logger) *UserPostgres {
	return &UserPostgres{db: db, logger: logger}
}

func (r *UserPostgres) GetUsers() ([]IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var listUsers []IndTask.User
	query := fmt.Sprint("SELECT * FROM users")

	rows, err := transaction.Query(query)

	for rows.Next() {
		var user IndTask.User
		if err := rows.Scan(&user); err != nil {
			logrus.Fatal(err)
		}
		listUsers = append(listUsers, user)
	}

	return listUsers, err

}

func (r *UserPostgres) CreateUser(user *IndTask.User) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var userId int
	createUserQuery := fmt.Sprint("INSERT INTO users (surname, user_name, patronymic, pasp_number, email, adress, birth_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	row := transaction.QueryRow(createUserQuery, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Adress, user.BirthDate)
	if err := row.Scan(&userId); err != nil {
		transaction.Rollback()
		return 0, err
	}
	return userId, transaction.Commit()

}

func (r *UserPostgres) ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	if method == "GET" {
		var user *IndTask.User
		query := fmt.Sprintf("SELECT * FROM users WHERE id = $1")

		row := transaction.QueryRow(query, userId)
		if err := row.Scan(&user); err != nil {
			transaction.Rollback()
			return nil, err
		}
		return user, nil
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE users SET (surname=$1, user_name=$2, patronymic=$3, pasp_number=$4, email=$5, adress=$6, birth_date=$7) WHERE id = $8")
		_, err := transaction.Exec(query, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Adress, user.BirthDate, userId)
		return nil, err
	}
	if method == "DELETE" {
		query := fmt.Sprint("DELETE FROM users WHERE id = $1")
		_, err := transaction.Exec(query, userId)
		return nil, err
	}
	return nil, transaction.Commit()

}
