package repository

import "database/sql"

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUsers() {

}

func (r *UserPostgres) CreateUser() {

}

func (r *UserPostgres) ChangeUser() {

}
