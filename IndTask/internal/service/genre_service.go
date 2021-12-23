package service

import (
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

type GenreService struct {
	repo repository.AppGenre
}

func NewGenreService(repo repository.AppGenre) *GenreService {
	return &GenreService{repo: repo}
}

func (g *GenreService) GetGenres() ([]IndTask.Genre, error) {
	return g.repo.GetGenres()
}

func (g *GenreService) CreateGenre(genre *IndTask.Genre) (int, error) {
	return g.repo.CreateGenre(genre)
}

func (g *GenreService) ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error) {
	return g.repo.ChangeGenre(genre, genreId, method)
}
