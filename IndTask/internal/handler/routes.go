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

	router.HandleFunc("/books", h.getBooks)
	router.HandleFunc("/books/create", h.createBook)
	router.HandleFunc("/books/change", h.changeBook) //реализовать методы GET, PUT
	router.HandleFunc("/list_books", h.getListBooks)
	router.HandleFunc("/list_books/change", h.changeListBooks)

	router.HandleFunc("/users/", h.getUsers)
	router.HandleFunc("/users/create", h.createUser)
	router.HandleFunc("/users/change", h.changeUser) //реализовать методы GET, PUT, DELETE

	router.HandleFunc("/in/getAll", h.getReturnActs)
	router.HandleFunc("/in/create", h.createReturnAct)
	router.HandleFunc("/in/getByUser", h.getReturnActsByUser)
	router.HandleFunc("/in/change", h.changeReturnAct) //реализовать методы GET, PUT

	router.HandleFunc("/out/getAll", h.getIssueActs)
	router.HandleFunc("/out/create", h.createIssueAct)
	router.HandleFunc("/out/getByUser", h.getIssueActsByUser)
	router.HandleFunc("/out/change", h.changeIssueAct) //реализовать методы GET, PUT

	router.HandleFunc("/authors/", h.getAuthors)
	router.HandleFunc("/authors/create", h.createAuthor)
	router.HandleFunc("/authors/change", h.changeAuthor) //реализовать методы GET, PUT, DELETE

	router.HandleFunc("/genres/", h.getGenres)
	router.HandleFunc("/genres/create", h.createGenre)
	router.HandleFunc("/genres/change", h.changeGenre) //реализовать методы GET, PUT, DELETE

	return router
}
