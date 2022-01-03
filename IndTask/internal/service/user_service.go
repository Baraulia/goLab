package service

import (
	"errors"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

var UserAlreadyExists = errors.New("user with that email already exists")
var UserDoesNotExists = errors.New("user with that id does not exists")

type UserService struct {
	repo repository.AppUser
}

func NewUserService(repo repository.AppUser) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) GetUsers() ([]IndTask.User, error) {
	return u.repo.GetUsers()
}

func (u *UserService) CreateUser(user *IndTask.User) (int, error) {
	listUsers, err := u.repo.GetUsers()
	if err != nil {
		logger.Errorf("Error when getting genres:%s", err)
		return 0, err
	}

	if user.UserName == "" || user.Surname == "" {
		logger.Error("necessary to fill in the Fields with the user_name and surname")
		return 0, fmt.Errorf("necessary to fill in the Fields with the user_name and surname")
	}

	for _, bdUser := range listUsers {
		if bdUser.Email == user.Email {
			logger.Error("User with the same email already exists")
			return bdUser.Id, UserAlreadyExists
		}
	}
	return u.repo.CreateUser(user)
}

func (u *UserService) ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error) {
	listUsers, err := u.repo.GetUsers()
	if err != nil {
		logger.Errorf("Error when getting users:%s", err)
		return nil, err
	}
	var userExist = false
	for _, bdUser := range listUsers {
		if bdUser.Id == userId {
			userExist = true
		}
	}
	if userExist == false {
		logger.Error("Such a user does not exist")
		return nil, UserDoesNotExists
	}
	return u.repo.ChangeUser(user, userId, method)

}
