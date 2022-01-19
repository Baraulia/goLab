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

	router.HandleFunc("/", h.getThreeBooks)
	router.HandleFunc("/books", h.getBooks)
	router.HandleFunc("/books/create", h.createBook)
	router.HandleFunc("/books/change", h.changeBook)
	router.HandleFunc("/list_books", h.getListBooks)
	router.HandleFunc("/list_books/create", h.createListBook)
	router.HandleFunc("/list_books/change", h.changeListBooks)

	router.HandleFunc("/users/", h.getUsers)
	router.HandleFunc("/users/create", h.createUser)
	router.HandleFunc("/users/change", h.changeUser)

	router.HandleFunc("/act/getAll", h.getActs)
	router.HandleFunc("/act/create", h.createIssueAct)
	router.HandleFunc("/act/getByUser", h.getActsByUser)
	router.HandleFunc("/act/change", h.changeAct)
	router.HandleFunc("/act/add", h.addReturnAct)

	router.HandleFunc("/authors/", h.getAuthors)
	router.HandleFunc("/authors/create", h.createAuthor)
	router.HandleFunc("/authors/change", h.changeAuthor)

	router.HandleFunc("/genres/", h.getGenres)
	router.HandleFunc("/genres/create", h.createGenre)
	router.HandleFunc("/genres/change", h.changeGenre)

	return router
}
