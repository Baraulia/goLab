package repository

import "database/sql"

type MovePostgres struct {
	db *sql.DB
}

func NewMovePostgres(db *sql.DB) *MovePostgres {
	return &MovePostgres{db: db}
}

func (r *MovePostgres) MoveInBook() {

}

func (r *MovePostgres) MoveOutBook() {

}
