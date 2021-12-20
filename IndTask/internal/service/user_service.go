package service

import "github.com/Baraulia/goLab/IndTask.git/internal/repository"

type UserService struct {
	repo repository.AppUser
}

func NewUserService(repo repository.AppUser) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) GetUsers() {

}

func (u *UserService) CreateUser() {

}

func (u *UserService) ChangeUser() {

}
