package handler

import (
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", h.getBooks)
	router.HandleFunc("/new_book", h.createBook)
	router.HandleFunc("/:book_id", h.changeBook)

	router.HandleFunc("/users", h.getUsers)
	router.HandleFunc("/users/create_user", h.createUser)
	router.HandleFunc("/users/:user_id", h.changeUser) //реализовать методы GET, PUT, DELETE

	router.HandleFunc("/in", h.moveInBook)
	router.HandleFunc("/out", h.moveOutBook)

	router.HandleFunc("/authors", h.getAuthors)
	router.HandleFunc("/authors/create", h.createAuthor)
	router.HandleFunc("/authors/:user_id", h.changeAuthor) //реализовать методы GET, PUT, DELETE

	router.HandleFunc("/genres", h.getGenres)
	router.HandleFunc("/genres/create", h.createGenre)
	router.HandleFunc("/genres/:user_id", h.changeGenre) //реализовать методы GET, PUT, DELETE

	return router
}
