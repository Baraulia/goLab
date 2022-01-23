package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/myErrors"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type UserService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) GetUsers(page int, sorting string) ([]IndTask.UserResponse, int, error) {
	users, pages, err := u.repo.GetUsers(page, sorting)
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting users from database:%w", err)
	}
	return users, pages, nil
}

func (u *UserService) CreateUser(user *IndTask.User) (*IndTask.User, error) {
	listUsers, _, err := u.repo.GetUsers(0, "email")
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting users from database:%w", err), Code: 500}
	}
	for _, bdUser := range listUsers {
		if bdUser.Email == user.Email {
			logger.Errorf("User with that email:%s already exists", user.Email)
			return nil, &myErrors.MyError{Err: fmt.Errorf("user with that email:%s already exists", user.Email), Code: 400}
		}
	}
	newUser, err := u.repo.CreateUser(user)
	if err != nil {
		return nil, &myErrors.MyError{Err: fmt.Errorf("error while creating user in database:%w", err), Code: 500}
	}
	return newUser, nil
}

func (u *UserService) ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error) {
	if err := u.repo.Validation.GetUserById(userId); err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			logger.Errorf("Such a user:%d does not exist", userId)
			return nil, &myErrors.MyError{Err: fmt.Errorf("changeUser:%s", e.Error()), Code: e.Status()}
		default:
			logger.Errorf("ChangeUser:%s", err)
			return nil, fmt.Errorf("changeUser:%w", err)
		}
	}
	if method == "GET" {
		oneUser, err := u.repo.GetOneUser(userId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while getting one user from database:%w", err), Code: 500}
		}
		return oneUser, nil
	}
	if method == "PUT" {
		upUser, err := u.repo.ChangeUser(user, userId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while changing user in database:%w", err), Code: 500}
		}
		return upUser, nil
	}
	if method == "DELETE" {
		err := u.repo.DeleteUser(userId)
		if err != nil {
			return nil, &myErrors.MyError{Err: fmt.Errorf("error while deleting one user from database:%w", err), Code: 500}
		}
		return nil, nil
	}
	return nil, nil
}

func (u *UserService) FoundUser(userSurname string) (*IndTask.User, error) {
	oneUser, err := u.repo.FoundUser(userSurname)
	if err != nil {
		return nil, fmt.Errorf("error while getting one user by userSurname: %s from database:%w", userSurname, err)
	}
	return oneUser, nil
}

func (u *UserService) SortTypeUser(sorting string) string {
	switch sorting {
	case "userSurnameDesc":
		return userSurnameDesc
	case "userSurnameAsc":
		return userSurnameAsc
	case "userNameDesc":
		return userNameDesc
	case "userNameAsc":
		return userNameAsc
	case "emailDesc":
		return emailDesc
	case "emailAsc":
		return emailAsc
	case "addressDesc":
		return addressDesc
	case "addressAsc":
		return addressAsc
	case "birthDateDesc":
		return birthDateDesc
	case "birthDateAsc":
		return birthDateAsc
	}
	return "surname"
}
