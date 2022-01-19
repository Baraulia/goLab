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
		return nil, fmt.Errorf("error while getting users from database:%w", err)
	}
	for _, bdUser := range listUsers {
		if bdUser.Email == user.Email {
			logger.Errorf("User with that email:%s already exists", user.Email)
			return nil, fmt.Errorf("user with that email:%s already exists", user.Email)
		}
	}
	newUser, err := u.repo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("error while creating user in database:%w", err)
	}
	return newUser, nil
}

func (u *UserService) ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error) {
	listUsers, _, err := u.repo.GetUsers(0, "id")
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
		oneUser, err := u.repo.GetOneUser(userId)
		if err != nil {
			return nil, fmt.Errorf("error while getting one user from database:%w", err)
		}
		return oneUser, nil
	}
	if method == "PUT" {
		upUser, err := u.repo.ChangeUser(user, userId)
		if err != nil {
			return nil, fmt.Errorf("error while changing user in database:%w", err)
		}
		return upUser, nil
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
