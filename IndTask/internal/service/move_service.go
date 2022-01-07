package service

import (
	"errors"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

var ReturnActDoesNotExists = errors.New("return act with that id does not exists")
var IssueActDoesNotExists = errors.New("issue act with that id does not exists")

type MoveService struct {
	repo repository.Repository
}

func NewMoveService(repo repository.Repository) *MoveService {
	return &MoveService{repo: repo}
}

func (s *MoveService) GetIssueActs(page int) ([]IndTask.IssueAct, error) {
	return s.repo.GetIssueActs(page)
}

func (s *MoveService) CreateIssueAct(issueAct *IndTask.IssueAct, method string) (int, error) {
	var err error
	if err = CheckUserBookExist(s, issueAct, method); err != nil {
		return 0, err
	}
	if issueAct.PreCost, err = CalcRentPreCost(issueAct, s); err != nil {
		return 0, err
	}
	if err = CheckAmoutActsByUser(issueAct, s); err != nil {
		return 0, err
	}
	issueAct.Status = "open"
	return s.repo.CreateIssueAct(issueAct)
}

func (s *MoveService) GetIssueActsByUser(userId int, page int) ([]IndTask.IssueAct, error) {
	if _, err := s.repo.ChangeUser(nil, userId, "GET"); err != nil {
		logger.Errorf("Such user with id = %d does not exist", userId)
		return nil, fmt.Errorf("such user with id = %d does not exist", userId)
	}
	return s.repo.GetIssueActsByUser(userId, false, page)
}

func (s *MoveService) ChangeIssueAct(issueAct *IndTask.IssueAct, actId int, method string) (*IndTask.IssueAct, error) {
	listActs, err := s.repo.GetIssueActs(0)
	if err != nil {
		logger.Errorf("Error when getting list issue acts:%s", err)
		return nil, err
	}
	var actExist = false
	for _, act := range listActs {
		if act.Id == actId {
			actExist = true
		}
	}
	if actExist == false {
		logger.Errorf("Such an act (id=%d) does not exist", actId)
		return nil, IssueActDoesNotExists
	}
	if method == "PUT" {
		err := CheckUserBookExist(s, issueAct, method)
		if err != nil {
			return nil, err
		}
		if issueAct.PreCost, err = CalcRentPreCost(issueAct, s); err != nil {
			return nil, err
		}
	}

	return s.repo.ChangeIssueAct(issueAct, actId, method)
}

func (s *MoveService) GetReturnActs(page int) ([]IndTask.ReturnAct, error) {
	return s.repo.GetReturnActs(page)
}
func (s *MoveService) CreateReturnAct(returnAct *IndTask.ReturnAct) (int, error) {

	listBookId, err := CheckIssueActExist(returnAct, s)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateReturnAct(returnAct, listBookId)
}
func (s *MoveService) GetReturnActsByUser(userId int, page int) ([]IndTask.ReturnAct, error) {
	if _, err := s.repo.ChangeUser(nil, userId, "GET"); err != nil {
		logger.Errorf("Such user with id = %d does not exist", userId)
		return nil, fmt.Errorf("such user with id = %d does not exist", userId)
	}
	return s.repo.GetReturnActsByUser(userId, page)
}
func (s *MoveService) ChangeReturnAct(returnAct *IndTask.ReturnAct, actId int, method string) (*IndTask.ReturnAct, error) {
	listActs, err := s.repo.GetReturnActs(0)
	if err != nil {
		logger.Errorf("Error when getting list return acts:%s", err)
		return nil, err
	}

	var actExist = false
	for _, act := range listActs {
		if act.Id == actId {
			actExist = true
		}
	}
	if actExist == false {
		logger.Errorf("Such an act (id=%d) does not exist", actId)
		return nil, ReturnActDoesNotExists
	}
	return s.repo.ChangeReturnAct(returnAct, actId, method)
}

func CheckUserBookExist(s *MoveService, issueAct *IndTask.IssueAct, method string) error {
	_, err := s.repo.ChangeUser(nil, issueAct.UserId, "GET")
	if err != nil {
		logger.Errorf("Such a user (id=%d) does not exist", issueAct.UserId)
		return fmt.Errorf("such a user (id=%d) does not exist", issueAct.UserId)
	}
	listBook, err := s.repo.ChangeListBook(nil, issueAct.ListBookId, "GET")
	if err != nil {
		logger.Errorf("Such a listbook (id=%d) does not exist", issueAct.ListBookId)
		return fmt.Errorf("such a listbook (id=%d) does not exist", issueAct.ListBookId)
	}
	if method == "POST" {
		if listBook.Issued {
			logger.Errorf("Such a listbook (id=%d) is already issued", issueAct.ListBookId)
			return fmt.Errorf("such a listbook (id=%d) is already issued", issueAct.ListBookId)
		}
	}

	return nil
}

func CalcRentPreCost(issueAct *IndTask.IssueAct, s *MoveService) (float64, error) {
	var book *IndTask.ListBooks
	var err error
	book, err = s.repo.AppBook.ChangeListBook(nil, issueAct.ListBookId, "GET")
	if err != nil {
		logger.Errorf("Error getting listBook by ListBookId=%d:%s", issueAct.ListBookId, err)
		return 0, err
	}
	rentTime := issueAct.RentalTime.Hours() / 24
	rentPreCost := rentTime * book.RentCost
	return rentPreCost, nil
}

func CheckIssueActExist(returnAct *IndTask.ReturnAct, s *MoveService) (int, error) {
	issueAct, err := s.ChangeIssueAct(nil, returnAct.IssueActId, "GET")
	if err != nil {
		logger.Errorf("IssueAct with id = %d does not exist", returnAct.IssueActId)
		return 0, err
	}
	if issueAct.Status == "closed" {
		logger.Errorf("IssueAct with id = %d already closed", returnAct.IssueActId)
		return 0, fmt.Errorf("issueAct with id = %d already closed", returnAct.IssueActId)
	}

	err = setCost(issueAct, returnAct, s)
	if err != nil {
		logger.Errorf("Error setting cost for returnAct where issueAct.Id = %d:%s", returnAct.IssueActId, err)
		return 0, fmt.Errorf("error setting cost for returnAct where issueAct.Id = %d:%s", returnAct.IssueActId, err)
	}
	return issueAct.ListBookId, nil
}

func CheckAmoutActsByUser(issueAct *IndTask.IssueAct, s *MoveService) error {
	listIssueAct, err := s.repo.GetIssueActsByUser(issueAct.UserId, true, 0)
	if err != nil {
		logger.Errorf("Error getting list IssueActs by user.Id = %d:%s", issueAct.UserId, err)
		return fmt.Errorf("error getting list IssueActs by user.Id = %d:%s", issueAct.UserId, err)
	}
	if len(listIssueAct) > 4 {
		logger.Errorf("This user (id = %d) has too many issue acts (%d)", issueAct.UserId, len(listIssueAct))
		return fmt.Errorf("this user (id = %d) has too many issue acts (%d)", issueAct.UserId, len(listIssueAct))
	}
	return nil
}

func setCost(issueAct *IndTask.IssueAct, returnAct *IndTask.ReturnAct, s *MoveService) error {
	var discount float64
	listIssueAct, err := s.repo.GetIssueActsByUser(issueAct.UserId, true, 0)
	if err != nil {
		logger.Errorf("Error getting list IssueActs by user.Id = %d:%s", issueAct.UserId, err)
		return fmt.Errorf("error getting list IssueActs by user.Id = %d:%s", issueAct.UserId, err)
	}
	if len(listIssueAct) > 2 && len(listIssueAct) < 5 {
		discount = 0.9
	} else if len(listIssueAct) == 5 {
		discount = 0.85
	} else {
		discount = 1
	}

	if returnAct.ReturnDate.After(issueAct.ReturnDate) {
		lateness := (returnAct.ReturnDate.Hour() - issueAct.ReturnDate.Hour()) / 24
		discount = discount + float64(lateness/100)
	}

	for _, act := range listIssueAct {
		act.PreCost = act.PreCost * discount
		act.Cost = act.PreCost + returnAct.Fine
		_, err := s.repo.ChangeIssueAct(&act, act.Id, "PUT")
		if err != nil {
			logger.Errorf("Error updating preCost in issueAct (id=%d):%s", act.Id, err)
			return fmt.Errorf("error updating preCost in issueAct (id=%d):%s", act.Id, err)
		}
	}
	return nil
}
