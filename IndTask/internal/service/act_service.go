package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/translate"
	"io/ioutil"
	"net/http"
	"os"
)

type ActService struct {
	repo repository.Repository
}

func NewActService(repo repository.Repository) *ActService {
	return &ActService{repo: repo}
}

func (s *ActService) GetActs(page int) ([]IndTask.Act, error) {
	acts, err := s.repo.AppAct.GetActs(page)
	if err != nil {
		return nil, fmt.Errorf("error while getting acts from database:%w", err)
	}
	return acts, nil
}

func (s *ActService) CreateIssueAct(act *IndTask.Act) (int, error) {
	var err error
	if err := s.repo.AppAct.CheckDuplicateBook(act); err != nil {
		return 0, fmt.Errorf("checkDuplicateBook:%w", err)
	}
	if err = s.CheckAct(act.ListBookId); err != nil {
		return 0, fmt.Errorf("checkAct:%w", err)
	}
	if act.PreCost, err = s.CalcRentPreCost(act.ListBookId, act.RentalTime); err != nil {
		return 0, fmt.Errorf("calcRentPreCost:%w", err)
	}
	if err = s.CheckAmountActsByUser(act.UserId); err != nil {
		return 0, fmt.Errorf("checkAmoutActsByUser:%w", err)
	}
	act.Status = "open"
	actId, err := s.repo.AppAct.CreateIssueAct(act)
	if err != nil {
		return 0, fmt.Errorf("error while creating act in database:%w", err)
	}
	return actId, nil
}

func (s *ActService) GetActsByUser(userId int, page int) ([]IndTask.Act, error) {
	if err := s.repo.GetUserById(userId); err != nil {
		logger.Errorf("Such user with id = %d does not exist", userId)
		return nil, fmt.Errorf("getActsByUser: %w", err)
	}
	acts, err := s.repo.AppAct.GetActsByUser(userId, false, page)
	if err != nil {
		return nil, fmt.Errorf("error while getting act by usrId=%d in database:%w", userId, err)
	}
	return acts, nil
}

func (s *ActService) ChangeAct(act *IndTask.Act, actId int, method string) (*IndTask.Act, error) {
	if err := s.repo.Validation.GetActById(actId, true); err != nil {
		logger.Errorf("ChangeAct: such act with id=%d does not exists:%s", actId, err)
		return nil, fmt.Errorf("changeAct:%w", err)
	}
	if method == "PUT" {
		if err := s.CheckAct(act.ListBookId); err != nil {
			return nil, fmt.Errorf("changeAct:%w", err)
		}
		if act.PreCost == 0 {
			preCost, err := s.CalcRentPreCost(act.ListBookId, act.RentalTime)
			if err != nil {
				return nil, fmt.Errorf("calcRentPreCost:%w", err)
			}
			act.PreCost = preCost
		}
		if err := s.CheckAmountActsByUser(act.UserId); err != nil {
			return nil, fmt.Errorf("checkAmoutActsByUser:%w", err)
		}
		err := s.repo.AppAct.ChangeAct(act, actId)
		if err != nil {
			return nil, fmt.Errorf("error while updating act in database:%w", err)
		}
	} else {
		oneAct, err := s.repo.AppAct.GetOneAct(actId)
		if err != nil {
			return nil, fmt.Errorf("error while getting act by id=%d in database:%w", actId, err)
		}
		return oneAct, nil
	}
	return nil, nil
}

func (s *ActService) AddReturnAct(returnAct *IndTask.ReturnAct) error {
	if err := s.CheckReturnAct(returnAct); err != nil {
		return fmt.Errorf("addReturnAct:%w", err)
	}
	err := s.repo.AppAct.AddReturnAct(returnAct)
	if err != nil {
		return fmt.Errorf("error while adding return act in database:%w", err)
	}
	return nil
}

func (s *ActService) CheckAct(listBookId int) error {
	book, err := s.repo.AppBook.GetOneListBook(listBookId)
	if err != nil {
		logger.Errorf("Error while getting instance of book whith id=%d", listBookId)
		return fmt.Errorf("error while getting instance of book whith id=%d", listBookId)
	}
	if book.Scrapped == true || book.Issued == true {
		logger.Errorf("That instance of book with id=%d is already issued or scrapped", listBookId)
		return fmt.Errorf("that instance of book with id=%d is already issued or scrapped", listBookId)
	}
	return nil
}

func (s *ActService) CalcRentPreCost(listBookId int, rentalTime IndTask.MyDuration) (float64, error) {
	book, err := s.repo.AppBook.GetOneListBook(listBookId)
	if err != nil {
		logger.Errorf("Error while getting instance of book by ListBookId=%d:%s", listBookId, err)
		return 0, fmt.Errorf("error while getting instance of book by ListBookId=%d:%w", listBookId, err)
	}
	rentTime := rentalTime.Hours() / 24
	rentPreCost := rentTime * book.RentCost
	return rentPreCost, nil
}

func (s *ActService) CheckAmountActsByUser(userId int) error {
	listAct, err := s.repo.GetActsByUser(userId, true, 0)
	if err != nil {
		logger.Errorf("Error getting instance of book by userId = %d:%s", userId, err)
		return fmt.Errorf("error getting instance of book by userId = %d:%w", userId, err)
	}
	if len(listAct) > 4 {
		logger.Errorf("This user (id = %d) has too many acts (%d)", userId, len(listAct))
		return fmt.Errorf("this user (id = %d) has too many acts (%d)", userId, len(listAct))
	}
	return nil
}

func (s *ActService) CheckReturnAct(returnAct *IndTask.ReturnAct) error {
	act, err := s.repo.AppAct.GetOneAct(returnAct.ActId)
	if err != nil {
		logger.Errorf("Error getting act by id = %d:%s", returnAct.ActId, err)
		return fmt.Errorf("error getting act by id = %d:%s", returnAct.ActId, err)
	}
	if act.UserId != returnAct.UserId {
		logger.Errorf("Incorrectly specified userId in return act (expected:%d, received:%d)", act.UserId, returnAct.UserId)
		return fmt.Errorf("incorrectly specified userId in return act (expected:%d, received:%d)", act.UserId, returnAct.UserId)
	}
	if act.ListBookId != returnAct.ListBookId {
		logger.Errorf("Incorrectly specified ListBookId in return act (expected:%d, received:%d)", act.ListBookId, returnAct.ListBookId)
		return fmt.Errorf("incorrectly specified ListBookId in return act (expected:%d, received:%d)", act.ListBookId, returnAct.ListBookId)
	}
	err = s.setCost(returnAct)
	if err != nil {
		return fmt.Errorf("setCost:%w", err)
	}
	return nil
}

func (s *ActService) setCost(returnAct *IndTask.ReturnAct) error {
	var discount float64
	acts, err := s.repo.AppAct.GetActsByUser(returnAct.UserId, true, 0)
	if err != nil {
		return fmt.Errorf("error while getting acts from database:%w", err)
	}
	if len(acts) == 0 {
		logger.Errorf("User with id=%d has already returned all books", returnAct.UserId)
		return fmt.Errorf("user with id=%d has already returned all books", returnAct.UserId)
	}
	var rawActs []IndTask.Act
	for _, act := range acts {
		if act.Cost == 0 {
			rawActs = append(rawActs, act)
		}
	}
	act, err := s.repo.AppAct.GetOneAct(returnAct.ActId)
	if err != nil {
		return fmt.Errorf("error while getting one act from database:%w", err)
	}

	if len(rawActs) > 2 && len(rawActs) < 5 {
		discount = 0.9
	} else if len(rawActs) == 5 {
		discount = 0.85
	} else {
		discount = 1
	}

	if returnAct.ReturnDate.After(act.ReturnDate) {
		lateness := (returnAct.ReturnDate.Hour() - act.ReturnDate.Hour()) / 24
		discount = discount + float64(lateness/100)
	}

	for _, act := range acts {
		act.PreCost = act.PreCost * discount
		act.Cost = act.PreCost
		err := s.repo.AppAct.ChangeAct(&act, act.Id)
		if err != nil {
			logger.Errorf("Error while updating preCost and cost in act (id=%d):%s", act.Id, err)
			return fmt.Errorf("error updating preCost and cost in act (id=%d):%s", act.Id, err)
		}
	}
	return nil
}

func InputFineFoto(req *http.Request, actId int) ([]string, error) {
	var filePathes []string
	m := req.MultipartForm
	files := m.File["file"]
	for i, headers := range files {
		reqfile, err := files[i].Open()
		if err != nil {
			logger.Errorf("InputFineFoto: error while getting file from multipart form:%s", err)
			return nil, fmt.Errorf("inputFineFoto: error while getting file from multipart form:%w", err)
		}
		defer reqfile.Close()
		fileBytes, err := ioutil.ReadAll(reqfile)
		if err != nil {
			logger.Errorf("InputFineFoto: error while reading file from request:%s", err)
			return nil, fmt.Errorf("inputFineFoto: error while reading file from request:%w", err)
		}
		directoryPath := fmt.Sprintf("images/fines/issueActId%d", actId)
		filePath := fmt.Sprintf("%s/%s", directoryPath, translate.Translate(headers.Filename))
		err = os.MkdirAll(directoryPath, 0777)
		if err != nil {
			logger.Errorf("InputFineFoto: error while creating directory (%s):%s", directoryPath, err)
			return nil, fmt.Errorf("inputFineFoto: error while creating directory (%s):%w", directoryPath, err)
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			logger.Errorf("InputFineFoto: error while opening file %s:%s", filePath, err)
			return nil, fmt.Errorf("inputFineFoto: error while opening file %s:%w", filePath, err)
		}
		defer file.Close()
		_, err = file.Write(fileBytes)
		if err != nil {
			logger.Errorf("InputFineFoto: error while writing file:%s", err)
			return nil, fmt.Errorf("inputFineFoto: error while writing file:%w", err)
		}
		filePathes = append(filePathes, filePath)
	}
	return filePathes, nil
}
