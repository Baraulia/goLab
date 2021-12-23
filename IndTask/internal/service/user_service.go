package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

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
	return u.repo.CreateUser(user)
}

func (u *UserService) ChangeUser(user *IndTask.User, userId int, method string) (*IndTask.User, error) {
	return u.repo.ChangeUser(user, userId, method)

}
