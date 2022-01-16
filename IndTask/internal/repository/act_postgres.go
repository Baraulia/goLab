package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/lib/pq"
	"time"
)

type ActPostgres struct {
	db *sql.DB
}

func NewActPostgres(db *sql.DB) *ActPostgres {
	return &ActPostgres{db: db}
}

var actLimit = 10

func (r *ActPostgres) GetActs(page int) ([]IndTask.Act, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetActs: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getActs: can not starts transaction:%w", err)
	}
	var listActs []IndTask.Act
	var rows *sql.Rows
	if page == 0 {
		query := "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetActs: can not executes a query:%s", err)
			return nil, fmt.Errorf("getActs:repository error:%w", err)
		}
	} else {
		query := "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act ORDER BY Id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, actLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("GetActs: can not executes a query:%s", err)
			return nil, fmt.Errorf("getActs:repository error:%w", err)
		}
	}
	for rows.Next() {
		var act IndTask.Act
		if err := rows.Scan(&act.Id, &act.UserId, &act.ListBookId, &act.RentalTime, &act.ReturnDate, &act.PreCost, &act.Cost, &act.Status, &act.ActualReturnDate, pq.Array(&act.Foto), &act.Fine, &act.ConditionDecrese, &act.Rating); err != nil {
			logger.Errorf("Error while scanning for act:%s", err)
			return nil, fmt.Errorf("getActs:repository error:%w", err)
		}
		listActs = append(listActs, act)
	}
	return listActs, transaction.Commit()

}
func (r *ActPostgres) CreateIssueAct(act *IndTask.Act) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateIssueAct: can not starts transaction:%s", err)
		return 0, fmt.Errorf("createIssueAct: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var actId int
	createIssueActQuery := "INSERT INTO act (user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id"
	row := transaction.QueryRow(createIssueActQuery, act.UserId, act.ListBookId, act.RentalTime, time.Now().Add(act.RentalTime.Duration), act.PreCost, act.Cost, act.Status, act.ActualReturnDate, pq.Array(act.Foto), act.Fine, act.ConditionDecrese, act.Rating)
	if err := row.Scan(&actId); err != nil {
		logger.Errorf("Error while scanning for actId:%s", err)
		return 0, fmt.Errorf("createIssueAct: error while scanning for actId:%w", err)
	}
	updateListBookQuery := "UPDATE list_books SET issued = true WHERE id = $1"
	_, err = transaction.Exec(updateListBookQuery, act.ListBookId)
	if err != nil {
		logger.Errorf("Error purpose value for listBook.issue where listBook.Id = %d:%s", act.ListBookId, err)
		return 0, fmt.Errorf("error purpose value for listBook.issue where listBook.Id = %d:%s", act.ListBookId, err)
	}
	return actId, transaction.Commit()
}

func (r *ActPostgres) GetActsByUser(userId int, forCost bool, page int) ([]IndTask.Act, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetActsByUser: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getActsByUser: can not starts transaction:%w", err)
	}
	var listActs []IndTask.Act
	var query string
	var rows *sql.Rows
	if page == 0 {
		if forCost {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=$1 AND status ='open'"
		} else {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=%d"
		}
		rows, err = transaction.Query(query, userId)
		if err != nil {
			logger.Errorf("GetActsByUser: can not executes a query:%s", err)
			return nil, fmt.Errorf("getActsByUser: repository error:%w", err)
		}
	} else {
		if forCost {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=$1 AND status ='open' AND cost = 0 ORDER BY Id LIMIT $2 OFFSET $3"
		} else {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=$1 ORDER BY Id LIMIT $2 OFFSET $3"
		}
		rows, err = transaction.Query(query, userId, actLimit, (page-1)*10)
		if err != nil {
			logger.Errorf("GetActsByUser: can not executes a query:%s", err)
			return nil, fmt.Errorf("getActsByUser: repository error:%w", err)
		}
	}
	for rows.Next() {
		var act IndTask.Act
		if err := rows.Scan(&act.Id, &act.UserId, &act.ListBookId, &act.RentalTime, &act.ReturnDate, &act.PreCost, &act.Cost, &act.Status, &act.ActualReturnDate, pq.Array(&act.Foto), &act.Fine, &act.ConditionDecrese, &act.Rating); err != nil {
			logger.Errorf("GetActsByUser: error while scanning for act:%s", err)
			return nil, fmt.Errorf("getActsByUser: repository error:%w", err)
		}
		listActs = append(listActs, act)
	}
	return listActs, transaction.Commit()
}

func (r *ActPostgres) ChangeAct(act *IndTask.Act, actId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeAct: can not starts transaction:%s", err)
		return fmt.Errorf("changeAct: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "UPDATE act SET user_id=$1, listbook_id=$2, rental_time=$3, return_date=$4, pre_cost=$5, cost=$6, status=$7, actual_return_date=$8, foto=$9, fine=$10, condition_decrese=$11, rating=$12 WHERE id = $13"
	_, err = transaction.Exec(query, act.UserId, act.ListBookId, act.RentalTime, act.ReturnDate, act.PreCost, act.Cost, act.Status, act.ActualReturnDate, pq.Array(act.Foto), act.Fine, act.ConditionDecrese, act.Rating, actId)
	if err != nil {
		logger.Errorf("Repository error while updating act:%s", err)
		return fmt.Errorf("changeAct: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *ActPostgres) GetOneAct(actId int) (*IndTask.Act, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneAct: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneAct: can not starts transaction:%w", err)
	}
	var act IndTask.Act
	query := "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE id = $1"
	row := transaction.QueryRow(query, actId)
	if err := row.Scan(&act.Id, &act.UserId, &act.ListBookId, &act.RentalTime, &act.ReturnDate, &act.PreCost, &act.Cost, &act.Status, &act.ActualReturnDate, pq.Array(&act.Foto), &act.Fine, &act.ConditionDecrese, &act.Rating); err != nil {
		logger.Errorf("Error while scanning for act:%s", err)
		return nil, fmt.Errorf("getOneAct:repository error:%w", err)
	}
	return &act, transaction.Commit()
}

func (r *ActPostgres) AddReturnAct(returnAct *IndTask.ReturnAct) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("AddReturnAct: can not starts transaction:%s", err)
		return fmt.Errorf("addReturnAct: can not starts transaction:%w", err)
	}
	var bookCondition int
	scrappedBook := false
	defer transaction.Rollback()
	addReturnActQuery := "UPDATE act SET cost=cost+$1, status=$2, actual_return_date=$3, foto=$4, fine=$5, condition_decrese=$6, rating=$7 WHERE id = $8"
	_, err = transaction.Exec(addReturnActQuery, returnAct.Fine, "closed", time.Now(), pq.Array(returnAct.Foto), returnAct.Fine, returnAct.ConditionDecrese, returnAct.Rating, returnAct.ActId)
	if err != nil {
		logger.Errorf("Repository error while adding return act:%s", err)
		return fmt.Errorf("addReturnAct: repository error:%w", err)
	}
	getBookQuery := "SELECT condition FROM list_books WHERE id=$1"
	row := transaction.QueryRow(getBookQuery, returnAct.ListBookId)
	if err := row.Scan(&bookCondition); err != nil {
		logger.Errorf("Error while scanning for bookCondition:%s", err)
		return fmt.Errorf("addReturnAct:repository error:%w", err)
	}
	if (bookCondition - returnAct.ConditionDecrese) <= 50 {
		logger.Infof("The instance of with id=%d should be scrapped", returnAct.ListBookId)
		scrappedBook = true
	}
	updateBookListQuery := "UPDATE list_books SET issued=$1, condition=condition-$2, rent_number = rent_number+1, scrapped=$3 WHERE id=$4"
	if _, err := transaction.Exec(updateBookListQuery, false, returnAct.ConditionDecrese, scrappedBook, returnAct.ListBookId); err != nil {
		logger.Errorf("Repository error while updating list_books:%s", err)
		return fmt.Errorf("addReturnAct: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *ActPostgres) CheckReturnData() ([]IndTask.Debtor, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CheckReturnData: can not starts transaction:%s", err)
		return nil, fmt.Errorf("checkReturnData: can not starts transaction:%w", err)
	}
	query := fmt.Sprint("SELECT users.email, users.user_name, books.book_name FROM users JOIN " +
		"act ON act.status = 'open' AND act.return_date < $1 AND users.id = act.user_id JOIN " +
		"list_books ON act.listbook_id = list_books.id JOIN " +
		"books ON list_books.book_id = books.id ")
	rows, err := transaction.Query(query, time.Now().Add(24*time.Hour))
	if err != nil {
		logger.Errorf("CheckReturnData: can not executes a query:%s", err)
		return nil, fmt.Errorf("checkReturnData:repository error:%w", err)
	}
	var listDeptors []IndTask.Debtor
	for rows.Next() {
		var debtor IndTask.Debtor
		if err := rows.Scan(&debtor.Email, &debtor.Name, &debtor.Book); err != nil {
			logger.Errorf("Error while scanning for debtor:%s", err)
			return nil, fmt.Errorf("checkReturnData:repository error:%w", err)
		}
		listDeptors = append(listDeptors, debtor)
	}

	return listDeptors, transaction.Commit()
}

func (r *ActPostgres) CheckDuplicateBook(act *IndTask.Act) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CheckDuplicateBook: can not starts transaction:%s", err)
		return fmt.Errorf("checkDuplicateBook: can not starts transaction:%w", err)
	}
	var exist bool
	fmt.Println("11")
	query := "SELECT EXISTS(SELECT 1 FROM books " +
		"JOIN list_books ON books.id=list_books.book_id AND list_books.id=$1 " +
		"JOIN act ON act.listbook_id=list_books.id AND act.user_id=$2)"
	row := transaction.QueryRow(query, act.ListBookId, act.UserId)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Error while scanning for issued book:%s", err)
		return fmt.Errorf("checkDuplicateBook: repository error:%w", err)
	}
	if !exist {
		return fmt.Errorf("such book is already issued to this user")
	}
	return transaction.Commit()
}
