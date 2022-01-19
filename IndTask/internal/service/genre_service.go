package service

import (
	"fmt"
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
	genres, err := g.repo.GetGenres()
	if err != nil {
		return nil, fmt.Errorf("error while getting genres from database:%w", err)
	}
	return genres, nil
}

func (g *GenreService) CreateGenre(genre *IndTask.Genre) (*IndTask.Genre, error) {
	listGenres, err := g.repo.GetGenres()
	if err != nil {
		return nil, fmt.Errorf("error while getting genres from database:%w", err)
	}
	for _, bdGenre := range listGenres {
		if bdGenre.GenreName == genre.GenreName {
			logger.Errorf("Genre with that name:%s already exists", genre.GenreName)
			return nil, fmt.Errorf("genre with that name:%s already exists", genre.GenreName)
		}
	}
	newGenre, err := g.repo.CreateGenre(genre)
	if err != nil {
		return nil, fmt.Errorf("error while creating genre in database:%w", err)
	}
	return newGenre, nil
}

func (g *GenreService) ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error) {
	listGenres, err := g.repo.GetGenres()
	if err != nil {
		return nil, fmt.Errorf("error while getting genres from database:%w", err)
	}
	var genreExist = false
	for _, bdGenre := range listGenres {
		if bdGenre.Id == genreId {
			genreExist = true
		}
	}
	if genreExist == false {
		logger.Errorf("Such a genre:%d does not exist", genreId)
		return nil, fmt.Errorf("such a genre:%d does not exist", genreId)
	}
	if method == "GET" {
		oneGenre, err := g.repo.GetOneGenre(genreId)
		if err != nil {
			return nil, fmt.Errorf("error while getting one genre from database:%w", err)
		}
		return oneGenre, nil
	}
	if method == "PUT" {
		for _, bdGenre := range listGenres {
			if bdGenre.GenreName == genre.GenreName {
				logger.Errorf("Genre with that name:%s already exists", genre.GenreName)
				return nil, fmt.Errorf("genre with that name:%s already exists", genre.GenreName)
			}
		}
		upGenre, err := g.repo.ChangeGenre(genre, genreId)
		if err != nil {
			return nil, fmt.Errorf("error while changing genre in database:%w", err)
		}
		return upGenre, nil
	}
	if method == "DELETE" {
		err := g.repo.DeleteGenre(genreId)
		if err != nil {
			return nil, fmt.Errorf("error while deleting one genre from database:%w", err)
		}
		return nil, nil
	}
	return nil, nil
}
