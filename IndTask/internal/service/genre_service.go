package service

import "github.com/Baraulia/goLab/IndTask.git/internal/repository"

type GenreService struct {
	repo repository.AppGenre
}

func NewGenreService(repo repository.AppGenre) *GenreService {
	return &GenreService{repo: repo}
}

func (u *GenreService) GetGenres() {

}

func (u *GenreService) CreateGenre() {

}

func (u *GenreService) ChangeGenre() {

}
