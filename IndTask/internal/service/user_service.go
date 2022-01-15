package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type UserService struct {
	repo repository.AppUser
}

func NewUserService(repo repository.AppUser) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) GetUsers(page int) ([]IndTask.User, error) {
	users, err := u.repo.GetUsers(page)
	if err != nil {
		return nil, fmt.Errorf("error while getting users from database:%w", err)
	}
	return users, nil
}

func (u *UserService) CreateUser(user *IndTask.User) (int, error) {
	listUsers, err := u.repo.GetUsers(0)
	if err != nil {
		return 0, fmt.Errorf("error while getting users from database:%w", err)
	}
	for _, bdUser := range listUsers {
		if bdUser.Email == user.Email {
			logger.Errorf("User with that email:%s already exists", user.Email)
			return bdUser.Id, fmt.Errorf("user with that email:%s already exists", user.Email)
		}
	}
	userId, err := u.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("error while creating user in database:%w", err)
	}
	return userId, nil
}

func (u *UserService) ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error) {
	listUsers, err := u.repo.GetUsers(0)
	if err != nil {
		return nil, fmt.Errorf("error while getting users from database:%w", err)
	}
	var userExist = false
	for _, bdUser := range listUsers {
		if bdUser.Id == userId {
			userExist = true
		}
	}
	if userExist == false {
		logger.Errorf("Such a user:%d does not exist", userId)
		return nil, fmt.Errorf("such a user:%d does not exist", userId)
	}
	if method == "GET" {
		user, err := u.repo.GetOneUser(userId)
		if err != nil {
			return nil, fmt.Errorf("error while getting one user from database:%w", err)
		}
		return user, nil
	}
	if method == "PUT" {
		err := u.repo.ChangeUser(user, userId)
		if err != nil {
			return nil, fmt.Errorf("error while changing user in database:%w", err)
		}
	}
	if method == "DELETE" {
		err := u.repo.DeleteUser(userId)
		if err != nil {
			return nil, fmt.Errorf("error while deleting one user from database:%w", err)
		}
		return nil, nil
	}
	return nil, nil
}
