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

func (r *ActPostgres) GetActs(page int) ([]IndTask.Act, int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetActs: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getActs: can not starts transaction:%w", err)
	}
	var listActs []IndTask.Act
	var rows *sql.Rows
	var pages int
	if page == 0 {
		query := "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetActs: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getActs:repository error:%w", err)
		}
	} else {
		query := "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act ORDER BY id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, actLimit, (page-1)*actLimit)
		if err != nil {
			logger.Errorf("GetActs: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getActs:repository error:%w", err)
		}
	}
	for rows.Next() {
		var act IndTask.Act
		if err := rows.Scan(&act.Id, &act.UserId, &act.ListBookId, &act.RentalTime, &act.ReturnDate, &act.PreCost, &act.Cost, &act.Status, &act.ActualReturnDate, pq.Array(&act.Foto), &act.Fine, &act.ConditionDecrese, &act.Rating); err != nil {
			logger.Errorf("Error while scanning for act:%s", err)
			return nil, 0, fmt.Errorf("getActs:repository error:%w", err)
		}
		listActs = append(listActs, act)
	}
	query := "SELECT CEILING(COUNT(id)/$1::float) FROM act"
	row := transaction.QueryRow(query, actLimit)
	if err := row.Scan(&pages); err != nil {
		logger.Errorf("Error while scanning for pages:%s", err)
		return nil, 0, fmt.Errorf("getActs: error while scanning for pages:%w", err)
	}
	return listActs, pages, transaction.Commit()

}
func (r *ActPostgres) CreateIssueAct(act *IndTask.Act) (*IndTask.Act, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateIssueAct: can not starts transaction:%s", err)
		return nil, fmt.Errorf("createIssueAct: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var newAct IndTask.Act
	createIssueActQuery := "INSERT INTO act (user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) " +
		"RETURNING id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating"
	row := transaction.QueryRow(createIssueActQuery, act.UserId, act.ListBookId, act.RentalTime, time.Now().Add(act.RentalTime.Duration), act.PreCost, act.Cost, act.Status, act.ActualReturnDate, pq.Array(act.Foto), act.Fine, act.ConditionDecrese, act.Rating)
	if err := row.Scan(&newAct.Id, &newAct.UserId, &newAct.ListBookId, &newAct.RentalTime, &newAct.ReturnDate, &newAct.PreCost, &newAct.Cost, &newAct.Status, &newAct.ActualReturnDate, pq.Array(&newAct.Foto), &newAct.Fine, &newAct.ConditionDecrese, &newAct.Rating); err != nil {
		logger.Errorf("Error while scanning for act:%s", err)
		return nil, fmt.Errorf("createIssueAct: error while scanning for act:%w", err)
	}
	updateListBookQuery := "UPDATE list_books SET issued = true WHERE id = $1"
	_, err = transaction.Exec(updateListBookQuery, act.ListBookId)
	if err != nil {
		logger.Errorf("Error purpose value for listBook.issue where listBook.Id = %d:%s", act.ListBookId, err)
		return nil, fmt.Errorf("error purpose value for listBook.issue where listBook.Id = %d:%s", act.ListBookId, err)
	}
	return &newAct, transaction.Commit()
}

func (r *ActPostgres) GetActsByUser(userId int, forCost bool, page int) ([]IndTask.Act, int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetActsByUser: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getActsByUser: can not starts transaction:%w", err)
	}
	var listActs []IndTask.Act
	var query string
	var rows *sql.Rows
	var pages int
	if page == 0 {
		if forCost {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=$1 AND status ='open'"
		} else {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=%d"
		}
		rows, err = transaction.Query(query, userId)
		if err != nil {
			logger.Errorf("GetActsByUser: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getActsByUser: repository error:%w", err)
		}
	} else {
		if forCost {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=$1 AND status ='open' AND cost = 0 ORDER BY Id LIMIT $2 OFFSET $3"
		} else {
			query = "SELECT id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating FROM act WHERE user_id=$1 ORDER BY Id LIMIT $2 OFFSET $3"
		}
		rows, err = transaction.Query(query, userId, actLimit, (page-1)*actLimit)
		if err != nil {
			logger.Errorf("GetActsByUser: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getActsByUser: repository error:%w", err)
		}
	}
	for rows.Next() {
		var act IndTask.Act
		if err := rows.Scan(&act.Id, &act.UserId, &act.ListBookId, &act.RentalTime, &act.ReturnDate, &act.PreCost, &act.Cost, &act.Status, &act.ActualReturnDate, pq.Array(&act.Foto), &act.Fine, &act.ConditionDecrese, &act.Rating); err != nil {
			logger.Errorf("GetActsByUser: error while scanning for act:%s", err)
			return nil, 0, fmt.Errorf("getActsByUser: repository error:%w", err)
		}
		listActs = append(listActs, act)
	}
	query = "SELECT CEILING(COUNT(id)/$1::float) FROM act WHERE user_id=$2"
	row := transaction.QueryRow(query, actLimit, userId)
	if err := row.Scan(&pages); err != nil {
		logger.Errorf("Error while scanning for pages:%s", err)
		return nil, 0, fmt.Errorf("getActsByUser: error while scanning for pages:%w", err)
	}
	return listActs, pages, transaction.Commit()
}

func (r *ActPostgres) ChangeAct(act *IndTask.Act, actId int) (*IndTask.Act, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeAct: can not starts transaction:%s", err)
		return nil, fmt.Errorf("changeAct: can not starts transaction:%w", err)
	}
	var upAct IndTask.Act
	defer transaction.Rollback()
	query := "UPDATE act SET user_id=$1, listbook_id=$2, rental_time=$3, return_date=$4, pre_cost=$5, cost=$6, status=$7, actual_return_date=$8, foto=$9, fine=$10, condition_decrese=$11, rating=$12 WHERE id = $13 " +
		"RETURNING id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating"
	row := transaction.QueryRow(query, act.UserId, act.ListBookId, act.RentalTime, act.ReturnDate, act.PreCost, act.Cost, act.Status, act.ActualReturnDate, pq.Array(act.Foto), act.Fine, act.ConditionDecrese, act.Rating, actId)
	if err := row.Scan(&upAct.Id, &upAct.UserId, &upAct.ListBookId, &upAct.RentalTime, &upAct.ReturnDate, &upAct.PreCost, &upAct.Cost, &upAct.Status, &upAct.ActualReturnDate, pq.Array(&upAct.Foto), &upAct.Fine, &upAct.ConditionDecrese, &upAct.Rating); err != nil {
		logger.Errorf("Error while scanning for act:%s", err)
		return nil, fmt.Errorf("changeAct: error while scanning for act:%w", err)
	}
	return &upAct, transaction.Commit()
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

func (r *ActPostgres) AddReturnAct(returnAct *IndTask.ReturnAct) (*IndTask.Act, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("AddReturnAct: can not starts transaction:%s", err)
		return nil, fmt.Errorf("addReturnAct: can not starts transaction:%w", err)
	}
	var upAct IndTask.Act
	var bookCondition int
	scrappedBook := false
	defer transaction.Rollback()
	addReturnActQuery := "UPDATE act SET cost=cost+$1, status=$2, actual_return_date=$3, foto=$4, fine=$5, condition_decrese=$6, rating=$7 WHERE id = $8 " +
		"RETURNING id, user_id, listbook_id, rental_time, return_date, pre_cost, cost, status, actual_return_date, foto, fine, condition_decrese, rating"
	row := transaction.QueryRow(addReturnActQuery, returnAct.Fine, "closed", time.Now(), pq.Array(returnAct.Foto), returnAct.Fine, returnAct.ConditionDecrese, returnAct.Rating, returnAct.ActId)
	if err := row.Scan(&upAct.Id, &upAct.UserId, &upAct.ListBookId, &upAct.RentalTime, &upAct.ReturnDate, &upAct.PreCost, &upAct.Cost, &upAct.Status, &upAct.ActualReturnDate, pq.Array(&upAct.Foto), &upAct.Fine, &upAct.ConditionDecrese, &upAct.Rating); err != nil {
		logger.Errorf("Error while scanning for act:%s", err)
		return nil, fmt.Errorf("addReturnAct: error while scanning for act:%w", err)
	}
	getBookQuery := "SELECT condition FROM list_books WHERE id=$1"
	row = transaction.QueryRow(getBookQuery, returnAct.ListBookId)
	if err := row.Scan(&bookCondition); err != nil {
		logger.Errorf("Error while scanning for bookCondition:%s", err)
		return nil, fmt.Errorf("addReturnAct:repository error:%w", err)
	}
	if (bookCondition - returnAct.ConditionDecrese) <= 50 {
		logger.Infof("The instance of with id=%d should be scrapped", returnAct.ListBookId)
		scrappedBook = true
	}
	updateBookListQuery := "UPDATE list_books SET issued=$1, condition=condition-$2, rent_number = rent_number+1, scrapped=$3 WHERE id=$4"
	if _, err := transaction.Exec(updateBookListQuery, false, returnAct.ConditionDecrese, scrappedBook, returnAct.ListBookId); err != nil {
		logger.Errorf("Repository error while updating list_books:%s", err)
		return nil, fmt.Errorf("addReturnAct: repository error:%w", err)
	}
	return &upAct, transaction.Commit()
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
	query := "SELECT EXISTS(SELECT books.id FROM books " +
		"JOIN list_books ON books.id=list_books.book_id AND list_books.id=$1 " +
		"JOIN act ON act.listbook_id=list_books.id AND act.user_id=$2 AND act.status='open')"
	row := transaction.QueryRow(query, act.ListBookId, act.UserId)
	if err := row.Scan(&exist); err != nil {
		logger.Errorf("Error while scanning for issued book:%s", err)
		return fmt.Errorf("checkDuplicateBook: repository error:%w", err)
	}
	if exist {
		return fmt.Errorf("such book is already issued to this user")
	}
	return transaction.Commit()
}
