package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/lib/pq"
	"time"
)

type MovePostgres struct {
	db *sql.DB
}

func NewMovePostgres(db *sql.DB) *MovePostgres {
	return &MovePostgres{db: db}
}

var actLimit = 10

func (r *MovePostgres) GetIssueActs(page int) ([]IndTask.IssueAct, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var listIssueActs []IndTask.IssueAct
	var rows *sql.Rows
	if page == 0 {
		query := fmt.Sprint("SELECT * FROM issue_act")
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	} else {
		query := fmt.Sprint("SELECT * FROM issue_act ORDER BY Id LIMIT $1 OFFSET $2")
		rows, err = transaction.Query(query, actLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	}
	for rows.Next() {
		var issueAct IndTask.IssueAct
		if err := rows.Scan(&issueAct.Id, &issueAct.UserId, &issueAct.ListBookId, &issueAct.RentalTime, &issueAct.ReturnDate, &issueAct.PreCost, &issueAct.Cost, &issueAct.Status); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		listIssueActs = append(listIssueActs, issueAct)
	}
	return listIssueActs, transaction.Commit()

}
func (r *MovePostgres) CreateIssueAct(issueAct *IndTask.IssueAct) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return 0, err
	}
	defer transaction.Rollback()
	var actId int
	createIssueActQuery := fmt.Sprint("INSERT INTO issue_act (user_id, listbook_id, rental_time, return_date, pre_cost, cost, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	row := transaction.QueryRow(createIssueActQuery, issueAct.UserId, issueAct.ListBookId, issueAct.RentalTime, time.Now().Add(issueAct.RentalTime.Duration), issueAct.PreCost, issueAct.Cost, issueAct.Status)
	if err := row.Scan(&actId); err != nil {
		logger.Errorf("Scan error:%s", err)
		return 0, err
	}
	updateListBookQuery := fmt.Sprintf("UPDATE list_books SET issued = true WHERE id = $1")
	_, err = transaction.Exec(updateListBookQuery, issueAct.ListBookId)
	if err != nil {
		logger.Errorf("Error purpose value for listBook.issue where listBook.Id = %d:%s", issueAct.ListBookId, err)
		return 0, fmt.Errorf("error purpose value for listBook.issue where listBook.Id = %d:%s", issueAct.ListBookId, err)
	}
	return actId, transaction.Commit()
}

func (r *MovePostgres) GetIssueActsByUser(userId int, forCost bool, page int) ([]IndTask.IssueAct, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var listIssueActs []IndTask.IssueAct
	var query string
	var rows *sql.Rows
	if page == 0 {
		if forCost {
			query = fmt.Sprintf("SELECT * FROM issue_act WHERE user_id=%d AND status ='open' AND cost = 0", userId)
		} else {
			query = fmt.Sprintf("SELECT * FROM issue_act WHERE user_id=%d", userId)
		}
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	} else {
		if forCost {
			query = fmt.Sprintf("SELECT * FROM issue_act WHERE user_id=%d AND status ='open' AND cost = 0 ORDER BY Id LIMIT $1 OFFSET $2", userId)
		} else {
			query = fmt.Sprintf("SELECT * FROM issue_act WHERE user_id=%d ORDER BY Id LIMIT $1 OFFSET $2", userId)
		}
		rows, err = transaction.Query(query, actLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	}

	for rows.Next() {
		var issueAct IndTask.IssueAct
		if err := rows.Scan(&issueAct.Id, &issueAct.UserId, &issueAct.ListBookId, &issueAct.RentalTime, &issueAct.ReturnDate, &issueAct.PreCost, &issueAct.Cost, &issueAct.Status); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		listIssueActs = append(listIssueActs, issueAct)
	}
	return listIssueActs, transaction.Commit()
}
func (r *MovePostgres) ChangeIssueAct(issueAct *IndTask.IssueAct, actId int, method string) (*IndTask.IssueAct, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var issueAct IndTask.IssueAct
		query := fmt.Sprint("SELECT * FROM issue_act WHERE id = $1")

		row := transaction.QueryRow(query, actId)
		if err := row.Scan(&issueAct.Id, &issueAct.UserId, &issueAct.ListBookId, &issueAct.RentalTime, &issueAct.ReturnDate, &issueAct.PreCost, &issueAct.Cost, &issueAct.Status); err != nil {
			logger.Errorf("Can not scan select from issueAct where id = %d", actId)
			return nil, err
		}
		return &issueAct, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE issue_act SET user_id=$1, listbook_id=$2, rental_time=$3, return_date=$4, pre_cost=$5, cost=$6, status=$7 WHERE id = $8")
		_, err := transaction.Exec(query, issueAct.UserId, issueAct.ListBookId, issueAct.RentalTime, time.Now().Add(issueAct.RentalTime.Duration), issueAct.PreCost, issueAct.Cost, issueAct.Status, actId)
		if err != nil {
			logger.Errorf("Can not update issueAct where id = %d", actId)
			return nil, err
		}
		return nil, transaction.Commit()
	}

	return nil, transaction.Rollback()
}

func (r *MovePostgres) GetReturnActs(page int) ([]IndTask.ReturnAct, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var listReturnActs []IndTask.ReturnAct
	var rows *sql.Rows
	if page == 0 {
		query := fmt.Sprint("SELECT * FROM return_act")
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	} else {
		query := fmt.Sprint("SELECT * FROM return_act ORDER BY Id LIMIT $1 OFFSET $2")
		rows, err = transaction.Query(query, actLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	}
	for rows.Next() {
		var returnAct IndTask.ReturnAct
		if err := rows.Scan(&returnAct.Id, &returnAct.IssueActId, &returnAct.ReturnDate, pq.Array(&returnAct.Foto), &returnAct.Fine, &returnAct.ConditionDecrese, &returnAct.Rating); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		listReturnActs = append(listReturnActs, returnAct)
	}
	return listReturnActs, transaction.Commit()
}

func (r *MovePostgres) CreateReturnAct(returnAct *IndTask.ReturnAct, listBookId int) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return 0, err
	}
	defer transaction.Rollback()

	var actId int
	createReturnActQuery := fmt.Sprint("INSERT INTO return_act (issue_act_id, return_date, foto, fine, condition_decrese, rating) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	row := transaction.QueryRow(createReturnActQuery, returnAct.IssueActId, time.Now(), pq.Array(returnAct.Foto), returnAct.Fine, returnAct.ConditionDecrese, returnAct.Rating)
	if err := row.Scan(&actId); err != nil {
		logger.Errorf("Scan error:%s", err)
		return 0, err
	}
	updateIssueActQuery := fmt.Sprintf("UPDATE issue_act SET status='closed' WHERE id = $1")
	if _, err := transaction.Exec(updateIssueActQuery, returnAct.IssueActId); err != nil {
		logger.Errorf("Error updating issueAct from returnAct.IssueActId = %d:%s", returnAct.IssueActId, err)
		return 0, err
	}

	updateBookListQuery := fmt.Sprintf("UPDATE list_books SET condition=condition-$1, rent_number = rent_number+1 WHERE id=$2")
	if _, err := transaction.Exec(updateBookListQuery, returnAct.ConditionDecrese, listBookId); err != nil {
		logger.Errorf("Error updating issueAct from returnAct.IssueActId = %d:%s", returnAct.IssueActId, err)
		return 0, err
	}
	return actId, transaction.Commit()
}
func (r *MovePostgres) GetReturnActsByUser(userId int, page int) ([]IndTask.ReturnAct, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}
	var listReturnActs []IndTask.ReturnAct
	var rows *sql.Rows
	if page == 0 {
		query := fmt.Sprint("SELECT DISTINCT return_act.id, return_act.issue_act_id, return_act.return_date, return_act.foto, return_act.fine, return_act.condition_decrese, return_act.rating FROM return_act JOIN issue_act ON user_id=$1")
		rows, err = transaction.Query(query, userId)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	} else {
		query := fmt.Sprint("SELECT DISTINCT return_act.id, return_act.issue_act_id, return_act.return_date, return_act.foto, return_act.fine, return_act.condition_decrese, return_act.rating FROM return_act JOIN issue_act ON user_id=$1 ORDER BY return_act.id LIMIT $2 OFFSET $3")
		rows, err = transaction.Query(query, userId, actLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("Can not executes a query:%s", err)
			return nil, err
		}
	}
	for rows.Next() {
		var returnAct IndTask.ReturnAct
		if err := rows.Scan(&returnAct.Id, &returnAct.IssueActId, &returnAct.ReturnDate, pq.Array(&returnAct.Foto), &returnAct.Fine, &returnAct.ConditionDecrese, &returnAct.Rating); err != nil {
			logger.Errorf("Scan error:%s", err)
			return nil, err
		}
		listReturnActs = append(listReturnActs, returnAct)
	}
	return listReturnActs, transaction.Commit()
}
func (r *MovePostgres) ChangeReturnAct(returnAct *IndTask.ReturnAct, actId int, method string) (*IndTask.ReturnAct, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("Can not begin transaction:%s", err)
		return nil, err
	}

	if method == "GET" {
		var returnAct IndTask.ReturnAct
		query := fmt.Sprint("SELECT * FROM return_act WHERE id = $1")
		row := transaction.QueryRow(query, actId)
		if err := row.Scan(&returnAct.Id, &returnAct.IssueActId, &returnAct.ReturnDate, pq.Array(&returnAct.Foto), &returnAct.Fine, &returnAct.ConditionDecrese, &returnAct.Rating); err != nil {
			logger.Errorf("Can not scan select from returnAct where id = %d", actId)
			return nil, err
		}
		return &returnAct, transaction.Commit()
	}

	if method == "PUT" {
		query := fmt.Sprint("UPDATE return_act SET issue_act_id=$1, return_date=$2, foto=$3, fine=$4, condition_decrese=$5, rating=$6 WHERE id = $7")
		_, err := transaction.Exec(query, returnAct.IssueActId, returnAct.ReturnDate, pq.Array(returnAct.Foto), returnAct.Fine, returnAct.ConditionDecrese, returnAct.Rating, actId)
		if err != nil {
			logger.Errorf("Can not update returnAct where id = %d", actId)
			return nil, err
		}
		return nil, transaction.Commit()
	}

	return nil, transaction.Rollback()
}
