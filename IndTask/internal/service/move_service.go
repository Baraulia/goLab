package service

import (
	"errors"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/translate"
	"io/ioutil"
	"net/http"
	"os"
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
	if _, err := s.repo.GetOneUser(userId); err != nil {
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
	listBookId, err := CheckIssueAct(returnAct, s)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateReturnAct(returnAct, listBookId)
}
func (s *MoveService) GetReturnActsByUser(userId int, page int) ([]IndTask.ReturnAct, error) {
	if _, err := s.repo.GetOneUser(userId); err != nil {
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

func CalcRentPreCost(issueAct *IndTask.IssueAct, s *MoveService) (float64, error) {
	book, err := s.repo.AppBook.GetOneListBook(issueAct.ListBookId)
	if err != nil {
		logger.Errorf("Error getting listBook by ListBookId=%d:%s", issueAct.ListBookId, err)
		return 0, err
	}
	rentTime := issueAct.RentalTime.Hours() / 24
	rentPreCost := rentTime * book.RentCost
	return rentPreCost, nil
}

func CheckIssueAct(returnAct *IndTask.ReturnAct, s *MoveService) (int, error) {
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

func InputFineFoto(req *http.Request, input *IndTask.ReturnAct) error {
	m := req.MultipartForm
	files := m.File["file"]
	for i, headers := range files {
		reqfile, err := files[i].Open()
		if err != nil {
			logger.Errorf("InputFineFoto: error while getting file from multipart form:%s", err)
			return fmt.Errorf("inputFineFoto: error while getting file from multipart form:%w", err)
		}
		defer reqfile.Close()
		fileBytes, err := ioutil.ReadAll(reqfile)
		if err != nil {
			logger.Errorf("InputFineFoto: error while reading file from request:%s", err)
			return fmt.Errorf("inputFineFoto: error while reading file from request:%w", err)
		}
		directoryPath := fmt.Sprintf("images/fines/issueActId%d", input.IssueActId)
		filePath := fmt.Sprintf("%s/%s", directoryPath, translate.Translate(headers.Filename))
		err = os.MkdirAll(directoryPath, 0777)
		if err != nil {
			logger.Errorf("InputFineFoto: error while creating directory (%s):%s", directoryPath, err)
			return fmt.Errorf("inputFineFoto: error while rcreating directory (%s):%w", directoryPath, err)
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			logger.Errorf("InputFineFoto: error while opening file %s:%s", filePath, err)
			return fmt.Errorf("inputFineFoto: error while opening file %s:%w", filePath, err)
		}
		defer file.Close()
		_, err = file.Write(fileBytes)
		if err != nil {
			logger.Errorf("InputFineFoto: error while writing file:%s", err)
			return fmt.Errorf("inputFineFoto: error while writing file:%w", err)
		}
		input.Foto = append(input.Foto, filePath)
	}
	return nil
}
