package handler

import (
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"net/http"
)

type Handler struct {
	logger   logging.Logger
	services *service.Service
}

func NewHandler(logger logging.Logger, services *service.Service) *Handler {
	return &Handler{logger: logger, services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", h.getBooks)
	router.HandleFunc("/new_book", h.createBook)
	router.HandleFunc("/:book_id", h.changeBook)

	router.HandleFunc("/users/", h.getUsers)
	router.HandleFunc("/users/create", h.createUser)
	router.HandleFunc("/users/change", h.changeUser) //реализовать методы GET, PUT, DELETE

	router.HandleFunc("/in", h.moveInBook)
	router.HandleFunc("/in/:user_id", h.getMoveInBook)
	router.HandleFunc("/out", h.moveOutBook)

	router.HandleFunc("/authors/", h.getAuthors)
	router.HandleFunc("/authors/create", h.createAuthor)
	router.HandleFunc("/authors/change", h.changeAuthor) //реализовать методы GET, PUT, DELETE

	router.HandleFunc("/genres/", h.getGenres)
	router.HandleFunc("/genres/create", h.createGenre)
	router.HandleFunc("/genres/change", h.changeGenre) //реализовать методы GET, PUT, DELETE

	return router
}
