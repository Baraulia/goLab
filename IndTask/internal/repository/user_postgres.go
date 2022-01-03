package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUsers() ([]IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	var listUsers []IndTask.User
	query := fmt.Sprint("SELECT * FROM users")
	rows, err := transaction.Query(query)
	if err != nil {
		logger.Errorf("Can not executes a query:%s", err)
		return nil, err
	}

	for rows.Next() {
		var user IndTask.User
		if err := rows.Scan(&user.Id, &user.Surname, &user.UserName, &user.Patronymic, &user.PaspNumber, &user.Email, &user.Adress, &user.BirthDate); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		listUsers = append(listUsers, user)
	}

	return listUsers, transaction.Commit()
}

func (r *UserPostgres) CreateUser(user *IndTask.User) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return 0, err
	}
	defer transaction.Rollback()

	var userId int
	createUserQuery := fmt.Sprint("INSERT INTO users (surname, user_name, patronymic, pasp_number, email, adress, birth_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	row := transaction.QueryRow(createUserQuery, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Adress, user.BirthDate)
	if err := row.Scan(&userId); err != nil {
		logger.Errorf("Scan error:%s", err)
		return 0, err
	}
	return userId, transaction.Commit()

}

func (r *UserPostgres) ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var user IndTask.User
		query := fmt.Sprintf("SELECT * FROM users WHERE id = $1")
		row := transaction.QueryRow(query, userId)
		if err := row.Scan(&user.Id, &user.Surname, &user.UserName, &user.Patronymic, &user.PaspNumber, &user.Email, &user.Adress, &user.BirthDate); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		return &user, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE users SET surname=$1, user_name=$2, patronymic=$3, pasp_number=$4, email=$5, adress=$6, birth_date=$7 WHERE id = $8")
		_, err := transaction.Exec(query, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Adress, user.BirthDate, userId)
		if err != nil {
			logger.Errorf("Update user error:%s", err)
			return nil, err
		}
		return nil, transaction.Commit()
	}

	if method == "DELETE" {
		query := fmt.Sprint("DELETE FROM users WHERE id = $1")
		_, err := transaction.Exec(query, userId)
		if err != nil {
			logger.Errorf("Delete user error:%s", err)
			return nil, err
		}
		return nil, transaction.Commit()

	}
	return nil, transaction.Rollback()

}
