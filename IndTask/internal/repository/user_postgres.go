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

var userLimit = 10

func (r *UserPostgres) GetUsers(page int) ([]IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetUsers: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getUsers: can not starts transaction:%w", err)
	}
	var listUsers []IndTask.User
	var rows *sql.Rows
	if page == 0 {
		query := "SELECT id, surname, user_name, patronymic, pasp_number, email, adress, birth_date FROM users"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetUsers: can not executes a query:%s", err)
			return nil, fmt.Errorf("getUsers:repository error:%w", err)
		}
	} else {
		query := "SELECT id, surname, user_name, patronymic, pasp_number, email, adress, birth_date FROM users ORDER BY Id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, actLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("GetUsers: can not executes a query:%s", err)
			return nil, fmt.Errorf("getUsers:repository error:%w", err)
		}
	}
	for rows.Next() {
		var user IndTask.User
		if err := rows.Scan(&user.Id, &user.Surname, &user.UserName, &user.Patronymic, &user.PaspNumber, &user.Email, &user.Adress, &user.BirthDate); err != nil {
			logger.Errorf("Error while scanning for user:%s", err)
			return nil, fmt.Errorf("getUsers:repository error:%w", err)
		}
		listUsers = append(listUsers, user)
	}

	return listUsers, transaction.Commit()
}

func (r *UserPostgres) CreateUser(user *IndTask.User) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateUser: can not starts transaction:%s", err)
		return 0, fmt.Errorf("createUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var userId int
	createUserQuery := "INSERT INTO users (surname, user_name, patronymic, pasp_number, email, adress, birth_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	row := transaction.QueryRow(createUserQuery, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Adress, user.BirthDate)
	if err := row.Scan(&userId); err != nil {
		logger.Errorf("Error while scanning for userId:%s", err)
		return 0, fmt.Errorf("createUser: error while scanning for userId:%w", err)
	}
	return userId, transaction.Commit()
}

func (r *UserPostgres) GetOneUser(userId int) (*IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneUser: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneUser: can not starts transaction:%w", err)
	}
	var user IndTask.User
	query := "SELECT id, surname, user_name, patronymic, pasp_number, email, adress, birth_date FROM users WHERE id = $1"
	row := transaction.QueryRow(query, userId)
	if err := row.Scan(&user.Id, &user.Surname, &user.UserName, &user.Patronymic, &user.PaspNumber, &user.Email, &user.Adress, &user.BirthDate); err != nil {
		logger.Errorf("Error while scanning for user:%s", err)
		return nil, fmt.Errorf("getOneUser: repository error:%w", err)
	}
	return &user, transaction.Commit()
}

func (r *UserPostgres) ChangeUser(user *IndTask.User, userId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeUser: can not starts transaction:%s", err)
		return fmt.Errorf("changeUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "UPDATE users SET surname=$1, user_name=$2, patronymic=$3, pasp_number=$4, email=$5, adress=$6, birth_date=$7 WHERE id = $8"
	_, err = transaction.Exec(query, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Adress, user.BirthDate, userId)
	if err != nil {
		logger.Errorf("Repository error while updating user:%s", err)
		return fmt.Errorf("changeUser: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *UserPostgres) DeleteUser(userId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("DeleteUser: can not starts transaction:%s", err)
		return fmt.Errorf("deleteUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM users WHERE id = $1"
	_, err = transaction.Exec(query, userId)
	if err != nil {
		logger.Errorf("Repository error while deleting user:%s", err)
		return fmt.Errorf("deleteUser: repository error:%w", err)
	}
	return transaction.Commit()
}
