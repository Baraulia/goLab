package service

import (
	"errors"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
)

var GenreAlreadyExists = errors.New("genre with that name already exists")
var GenreDoesNotExists = errors.New("genre with that id does not exists")

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
	listGenres, err := g.repo.GetGenres()
	if err != nil {
		logger.Errorf("Error when getting genres:%s", err)
		return 0, err
	}
	for _, bdGenre := range listGenres {
		if bdGenre.GenreName == genre.GenreName {
			logger.Error("Genre with the same name already exists")
			return bdGenre.Id, GenreAlreadyExists
		}
	}
	return g.repo.CreateGenre(genre)
}

func (g *GenreService) ChangeGenre(genre *IndTask.Genre, genreId int, method string) (*IndTask.Genre, error) {
	listGenres, err := g.repo.GetGenres()
	if err != nil {
		logger.Errorf("Error when getting genres:%s", err)
		return nil, err
	}
	var genreExist = false
	var genreNameDuplicate = false
	for _, bdGenre := range listGenres {
		if bdGenre.Id == genreId {
			genreExist = true
		}
	}
	if genreExist == false {
		logger.Error("Such a genre does not exist")
		return nil, GenreDoesNotExists
	}
	if method == "PUT" {
		for _, bdGenre := range listGenres {
			if bdGenre.GenreName == genre.GenreName {
				genreNameDuplicate = true
			}
		}
	}
	if genreNameDuplicate {
		logger.Error("Genre with the same name already exists")
		return nil, GenreAlreadyExists
	}
	return g.repo.ChangeGenre(genre, genreId, method)

}
