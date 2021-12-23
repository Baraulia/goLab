package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/sirupsen/logrus"
	"time"
)

type MovePostgres struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewMovePostgres(db *sql.DB, logger *logging.Logger) *MovePostgres {
	return &MovePostgres{db: db, logger: logger}
}

func (r *MovePostgres) MoveInBook(issueAct *IndTask.IssueAct) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var issueActId int
	createIssueActQuery := fmt.Sprint("INSERT INTO issue_act (user_id, book_id, rental_time, return_date, pre_cost, status) VALUES ($1, $2, $3, $4, $5, $6,) RETURNING id")
	row := transaction.QueryRow(createIssueActQuery, issueAct.UserId, issueAct.BookId, issueAct.RentalTime, time.Now().Add(issueAct.RentalTime), issueAct.PreCost, true)
	if err := row.Scan(&issueActId); err != nil {
		transaction.Rollback()
		return 0, err
	}
	correctingRentNumberQuery := fmt.Sprint("UPDATE listBooks SET (rent_number=rent_number+1) WHERE book_id=$1")
	_, err = transaction.Exec(correctingRentNumberQuery, issueAct.BookId)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	return issueActId, transaction.Commit()
}

func (r *MovePostgres) GetMoveInBooks(userId int) ([]IndTask.IssueAct, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var listIssueActs []IndTask.IssueAct
	query := fmt.Sprint("SELECT * FROM issue_act WHERE user_id=$1")

	rows, err := transaction.Query(query, userId)

	for rows.Next() {
		var issueAct IndTask.IssueAct
		if err := rows.Scan(&issueAct); err != nil {
			logrus.Fatal(err)
		}
		listIssueActs = append(listIssueActs, issueAct)
	}
	return listIssueActs, err
}

func (r *MovePostgres) MoveOutBook(returnAct *IndTask.ReturnAct) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var returnActId int
	createReturnActQuery := fmt.Sprint("INSERT INTO return_act (user_id, book_id, cost, return_date, foto, fine, condition_decrese, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	row := transaction.QueryRow(createReturnActQuery, returnAct.UserId, returnAct.BookId, returnAct.Cost, time.Now(), returnAct.Foto, returnAct.Fine, returnAct.ConditionDecrese, returnAct.Rating)
	if err := row.Scan(&returnActId); err != nil {
		transaction.Rollback()
		return 0, err
	}
	if returnAct.ConditionDecrese != 0 {
		correctingConditionQuery := fmt.Sprint("UPDATE listBooks SET (conditions=conditions-$1) WHERE book_id=$2")
		_, err := transaction.Exec(correctingConditionQuery, returnAct.ConditionDecrese, returnAct.BookId)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}
	}
	return returnActId, transaction.Commit()

}
