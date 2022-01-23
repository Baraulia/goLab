package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/myErrors"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/translate"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ActService struct {
	repo repository.Repository
	cfg  *config.Config
}

func NewActService(repo repository.Repository, cfg *config.Config) *ActService {
	return &ActService{repo: repo, cfg: cfg}
}

func (s *ActService) GetActs(page int) ([]IndTask.Act, int, error) {
	acts, pages, err := s.repo.AppAct.GetActs(page)
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting acts from database:%w", err)
	}
	return acts, pages, nil
}

func (s *ActService) CreateIssueAct(act *IndTask.Act) (*IndTask.Act, error) {
	if err := s.repo.AppAct.CheckDuplicateBook(act); err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			return nil, &myErrors.MyError{Err: fmt.Errorf("checkDuplicateBook:%s", e.Error()), Code: e.Status()}
		default:
			return nil, fmt.Errorf("checkDuplicateBook:%w", err)
		}
	}

	if err := s.CheckAct(act.ListBookId, 0); err != nil {
		return nil, err
	}
	var err error
	act.PreCost, err = s.CalcRentPreCost(act.ListBookId, act.RentalTime)
	if err != nil {
		return nil, err
	}
	if err := s.CheckAmountActsByUser(act.UserId); err != nil {
		return nil, err
	}
	act.Status = "open"
	newAct, err := s.repo.AppAct.CreateIssueAct(act)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while creating act in database:%w", err), Code: 500}
	}
	return newAct, nil
}

func (s *ActService) GetActsByUser(userId int, page int) ([]IndTask.Act, int, error) {
	if err := s.repo.GetUserById(userId); err != nil {
		logger.Errorf("Such user with id = %d does not exist", userId)
		switch e := err.(type) {
		case myErrors.Error:
			return nil, 0, &myErrors.MyError{Err: fmt.Errorf("getActsByUser: %w", err), Code: e.Status()}
		default:
			return nil, 0, fmt.Errorf("getActsByUser: %w", err)
		}
	}
	acts, pages, err := s.repo.AppAct.GetActsByUser(userId, false, page)
	if err != nil {
		return nil, 0, &myErrors.MyError{Err: fmt.Errorf("error while getting act by usrId=%d in database:%w", userId, err), Code: 500}
	}
	return acts, pages, nil
}

func (s *ActService) ChangeAct(act *IndTask.Act, actId int, method string) (*IndTask.Act, error) {
	if err := s.repo.Validation.GetActById(actId, true); err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			logger.Errorf("ChangeAct: such act with id=%d does not exists:%s", actId, err)
			return nil, &myErrors.MyError{Err: fmt.Errorf("changeAct:%s", e.Error()), Code: e.Status()}
		default:
			logger.Errorf("ChangeAct:%s", err)
			return nil, fmt.Errorf("changeAct:%w", err)
		}
	}
	var oneAct *IndTask.Act
	var err error
	if method == "PUT" {
		if err := s.CheckAct(act.ListBookId, actId); err != nil {
			switch e := err.(type) {
			case myErrors.Error:
				return nil, &myErrors.MyError{Err: fmt.Errorf("changeAct:%s", e.Error()), Code: e.Status()}
			default:
				return nil, fmt.Errorf("changeAct:%w", err)
			}
		}
		if act.PreCost == 0 {
			preCost, err := s.CalcRentPreCost(act.ListBookId, act.RentalTime)
			if err != nil {
				return nil, err
			}
			act.PreCost = preCost
		}
		if err := s.CheckAmountActsByUser(act.UserId); err != nil {
			switch e := err.(type) {
			case myErrors.Error:
				return nil, &myErrors.MyError{Err: fmt.Errorf("checkAmoutActsByUser:%s", e.Error()), Code: e.Status()}
			default:
				return nil, fmt.Errorf("checkAmoutActsByUser:%w", err)
			}
		}
		oneAct, err = s.repo.AppAct.ChangeAct(act, actId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while updating act in database:%w", err), Code: 500}
		}
	} else {
		oneAct, err = s.repo.AppAct.GetOneAct(actId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting act by id=%d in database:%w", actId, err), Code: 500}
		}
	}
	return oneAct, nil
}

func (s *ActService) AddReturnAct(returnAct *IndTask.ReturnAct) (*IndTask.Act, error) {
	if err := s.CheckReturnAct(returnAct); err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			return nil, &myErrors.MyError{Err: fmt.Errorf("addReturnAct:%s", e.Error()), Code: e.Status()}
		default:
			return nil, fmt.Errorf("addReturnAct:%w", err)
		}
	}
	act, err := s.repo.AppAct.AddReturnAct(returnAct)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while adding return act in database:%w", err), Code: 500}
	}
	return act, nil
}

func (s *ActService) CheckAct(listBookId int, actId int) error {
	book, err := s.repo.AppBook.GetOneListBook(listBookId)
	if err != nil {
		logger.Errorf("Error while getting instance of book whith id=%d", listBookId)
		return &myErrors.MyError{Err: fmt.Errorf("error while getting instance of book whith id=%d", listBookId), Code: 500}
	}
	if actId != 0 {
		act, err := s.repo.AppAct.GetOneAct(actId)
		if err != nil {
			logger.Errorf("Error while getting act whith id=%d", actId)
			return &myErrors.MyError{Err: fmt.Errorf("error while getting act whith id=%d", actId), Code: 500}
		}
		if act.ListBookId == listBookId {
			if book.Scrapped == true {
				logger.Errorf("That instance of book with id=%d is already scrapped", listBookId)
				return &myErrors.MyError{Err: fmt.Errorf("that instance of book with id=%d is already scrapped", listBookId), Code: 400}
			}
		}
	} else {
		if book.Scrapped == true || book.Issued == true {
			logger.Errorf("That instance of book with id=%d is already issued or scrapped", listBookId)
			return &myErrors.MyError{Err: fmt.Errorf("that instance of book with id=%d is already issued or scrapped", listBookId), Code: 400}
		}
	}
	return nil
}

func (s *ActService) CalcRentPreCost(listBookId int, rentalTime IndTask.MyDuration) (float64, error) {
	book, err := s.repo.AppBook.GetOneListBook(listBookId)
	if err != nil {
		logger.Errorf("Error while getting instance of book by ListBookId=%d:%s", listBookId, err)
		return 0, &myErrors.MyError{Err: fmt.Errorf("error while getting instance of book whith id=%d", listBookId), Code: 500}
	}
	rentTime := rentalTime.Hours() / 24
	rentPreCost := rentTime * book.RentCost
	return rentPreCost, nil
}

func (s *ActService) CheckAmountActsByUser(userId int) error {
	listAct, _, err := s.repo.GetActsByUser(userId, true, 0)
	if err != nil {
		logger.Errorf("Error getting instance of book by userId = %d:%s", userId, err)
		return &myErrors.MyError{Err: fmt.Errorf("error getting instance of book by userId = %d:%w", userId, err), Code: 500}
	}
	if len(listAct) > 4 {
		logger.Errorf("This user (id = %d) has too many acts (%d)", userId, len(listAct))
		return &myErrors.MyError{Err: fmt.Errorf("this user (id = %d) has too many acts (%d)", userId, len(listAct)), Code: 400}
	}
	return nil
}

func (s *ActService) CheckReturnAct(returnAct *IndTask.ReturnAct) error {
	act, err := s.repo.AppAct.GetOneAct(returnAct.ActId)
	if err != nil {
		logger.Errorf("Error getting act by id = %d:%s", returnAct.ActId, err)
		return &myErrors.MyError{Err: fmt.Errorf("error getting act by id = %d:%s", returnAct.ActId, err), Code: 500}
	}
	if act.UserId != returnAct.UserId {
		logger.Errorf("Incorrectly specified userId in return act (expected:%d, received:%d)", act.UserId, returnAct.UserId)
		return &myErrors.MyError{Err: fmt.Errorf("incorrectly specified userId in return act (expected:%d, received:%d)", act.UserId, returnAct.UserId), Code: 400}
	}
	if act.ListBookId != returnAct.ListBookId {
		logger.Errorf("Incorrectly specified ListBookId in return act (expected:%d, received:%d)", act.ListBookId, returnAct.ListBookId)
		return &myErrors.MyError{Err: fmt.Errorf("incorrectly specified ListBookId in return act (expected:%d, received:%d)", act.ListBookId, returnAct.ListBookId), Code: 400}
	}
	err = s.setCost(returnAct)
	if err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			return &myErrors.MyError{Err: fmt.Errorf("setCost:%s", e.Error()), Code: e.Status()}
		default:
			return fmt.Errorf("setCost:%w", err)
		}
	}
	return nil
}

func (s *ActService) setCost(returnAct *IndTask.ReturnAct) error {
	var discount float64
	acts, _, err := s.repo.AppAct.GetActsByUser(returnAct.UserId, true, 0)
	if err != nil {
		return &myErrors.MyError{Err: fmt.Errorf("error while getting acts from database:%w", err), Code: 500}
	}
	if len(acts) == 0 {
		logger.Errorf("User with id=%d has already returned all books", returnAct.UserId)
		return &myErrors.MyError{Err: fmt.Errorf("user with id=%d has already returned all books", returnAct.UserId), Code: 400}
	}
	var rawActs []IndTask.Act
	for _, act := range acts {
		if act.Cost == 0 {
			rawActs = append(rawActs, act)
		}
	}
	act, err := s.repo.AppAct.GetOneAct(returnAct.ActId)
	if err != nil {
		return &myErrors.MyError{Err: fmt.Errorf("error while getting one act from database:%w", err), Code: 500}
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
		_, err := s.repo.AppAct.ChangeAct(&act, act.Id)
		if err != nil {
			logger.Errorf("Error while updating preCost and cost in act (id=%d):%s", act.Id, err)
			return &myErrors.MyError{Err: fmt.Errorf("error updating preCost and cost in act (id=%d):%s", act.Id, err), Code: 500}
		}
	}
	return nil
}

func (s *ActService) InputFineFoto(req *http.Request, actId int) ([]string, error) {
	var filePathes []string
	m := req.MultipartForm
	files := m.File["file"]
	for i, headers := range files {
		reqfile, err := files[i].Open()
		if err != nil {
			logger.Errorf("InputFineFoto: error while getting file from multipart form:%s", err)
			return nil, &myErrors.MyError{Err: fmt.Errorf("inputFineFoto: error while getting file from multipart form:%w", err), Code: 400}
		}
		defer reqfile.Close()
		fileBytes, err := ioutil.ReadAll(reqfile)
		if err != nil {
			logger.Errorf("InputFineFoto: error while reading file from request:%s", err)
			return nil, &myErrors.MyError{Err: fmt.Errorf("inputFineFoto: error while reading file from request:%w", err), Code: 500}
		}
		directoryPath := fmt.Sprintf("%simages/fines/issueActId%d", s.cfg.FilePath, actId)
		filePath := fmt.Sprintf("%s/%s", directoryPath, translate.Translate(headers.Filename))
		err = os.MkdirAll(directoryPath, 0777)
		if err != nil {
			logger.Errorf("InputFineFoto: error while creating directory (%s):%s", directoryPath, err)
			return nil, &myErrors.MyError{Err: fmt.Errorf("inputFineFoto: error while creating directory (%s):%w", directoryPath, err), Code: 500}
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			logger.Errorf("InputFineFoto: error while opening file %s:%s", filePath, err)
			return nil, &myErrors.MyError{Err: fmt.Errorf("inputFineFoto: error while opening file %s:%w", filePath, err), Code: 500}
		}
		defer file.Close()
		_, err = file.Write(fileBytes)
		if err != nil {
			logger.Errorf("InputFineFoto: error while writing file:%s", err)
			return nil, &myErrors.MyError{Err: fmt.Errorf("inputFineFoto: error while writing file:%w", err), Code: 500}
		}
		filePath = strings.Replace(filePath, s.cfg.FilePath, "", 1)
		filePathes = append(filePathes, filePath)
	}
	return filePathes, nil
}
