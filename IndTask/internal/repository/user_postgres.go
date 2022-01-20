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

func (r *UserPostgres) GetUsers(page int, sorting string) ([]IndTask.UserResponse, int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetUsers: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getUsers: can not starts transaction:%w", err)
	}
	var listUsers []IndTask.UserResponse
	var rows *sql.Rows
	var pages int
	if page == 0 {
		query := fmt.Sprintf("SELECT id, surname, user_name, email, address, birth_date FROM users ORDER BY %s ", sorting)
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetUsers: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getUsers:repository error:%w", err)
		}
	} else {
		query := fmt.Sprintf("SELECT id, surname, user_name, email, address, birth_date FROM users ORDER BY %s LIMIT $1 OFFSET $2", sorting)
		rows, err = transaction.Query(query, userLimit, (page-1)*userLimit)
		if err != nil {
			logger.Errorf("GetUsers: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getUsers:repository error:%w", err)
		}
	}
	for rows.Next() {
		var user IndTask.UserResponse
		if err := rows.Scan(&user.Id, &user.Surname, &user.UserName, &user.Email, &user.Address, &user.BirthDate); err != nil {
			logger.Errorf("Error while scanning for user:%s", err)
			return nil, 0, fmt.Errorf("getUsers:repository error:%w", err)
		}
		listUsers = append(listUsers, user)
	}
	query := "SELECT CEILING(COUNT(id)/$1::float) FROM users"
	row := transaction.QueryRow(query, userLimit)
	if err := row.Scan(&pages); err != nil {
		logger.Errorf("Error while scanning for pages:%s", err)
		return nil, 0, fmt.Errorf("getUsers: error while scanning for pages:%w", err)
	}
	return listUsers, pages, transaction.Commit()
}

func (r *UserPostgres) CreateUser(user *IndTask.User) (*IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateUser: can not starts transaction:%s", err)
		return nil, fmt.Errorf("createUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var newUser IndTask.User
	createUserQuery := "INSERT INTO users (surname, user_name, patronymic, pasp_number, email, address, birth_date) VALUES ($1, $2, $3, $4, $5, $6, $7) " +
		"RETURNING id, surname, user_name, patronymic, pasp_number, email, address, birth_date"
	row := transaction.QueryRow(createUserQuery, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Address, user.BirthDate)
	if err := row.Scan(&newUser.Id, &newUser.Surname, &newUser.UserName, &newUser.Patronymic, &newUser.PaspNumber, &newUser.Email, &newUser.Address, &newUser.BirthDate); err != nil {
		logger.Errorf("Error while scanning for user:%s", err)
		return nil, fmt.Errorf("createUser: error while scanning for user:%w", err)
	}
	return &newUser, transaction.Commit()
}

func (r *UserPostgres) FoundUser(userSurname string) (*IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("FoundUser: can not starts transaction:%s", err)
		return nil, fmt.Errorf("foundUser: can not starts transaction:%w", err)
	}
	var user IndTask.User
	query := "SELECT id, surname, user_name, patronymic, pasp_number, email, address, birth_date FROM users WHERE surname = $1"
	row := transaction.QueryRow(query, userSurname)
	if err := row.Scan(&user.Id, &user.Surname, &user.UserName, &user.Patronymic, &user.PaspNumber, &user.Email, &user.Address, &user.BirthDate); err != nil {
		logger.Errorf("Error while scanning for user:%s", err)
		return nil, fmt.Errorf("foundUser: repository error:%w", err)
	}
	return &user, transaction.Commit()
}

func (r *UserPostgres) GetOneUser(userId int) (*IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneUser: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneUser: can not starts transaction:%w", err)
	}
	var user IndTask.User
	query := "SELECT id, surname, user_name, patronymic, pasp_number, email, address, birth_date FROM users WHERE id = $1"
	row := transaction.QueryRow(query, userId)
	if err := row.Scan(&user.Id, &user.Surname, &user.UserName, &user.Patronymic, &user.PaspNumber, &user.Email, &user.Address, &user.BirthDate); err != nil {
		logger.Errorf("Error while scanning for user:%s", err)
		return nil, fmt.Errorf("getOneUser: repository error:%w", err)
	}
	return &user, transaction.Commit()
}

func (r *UserPostgres) ChangeUser(user *IndTask.User, userId int) (*IndTask.User, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeUser: can not starts transaction:%s", err)
		return nil, fmt.Errorf("changeUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var upUser IndTask.User
	query := "UPDATE users SET surname=$1, user_name=$2, patronymic=$3, pasp_number=$4, email=$5, address=$6, birth_date=$7 WHERE id = $8 " +
		"RETURNING id, surname, user_name, patronymic, pasp_number, email, address, birth_date"
	row := transaction.QueryRow(query, user.Surname, user.UserName, user.Patronymic, user.PaspNumber, user.Email, user.Address, user.BirthDate, userId)
	if err := row.Scan(&upUser.Id, &upUser.Surname, &upUser.UserName, &upUser.Patronymic, &upUser.PaspNumber, &upUser.Email, &upUser.Address, &upUser.BirthDate); err != nil {
		logger.Errorf("Error while scanning for user:%s", err)
		return nil, fmt.Errorf("createUser: error while scanning for user:%w", err)
	}
	return &upUser, transaction.Commit()
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
