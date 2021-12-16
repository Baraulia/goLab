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
	router.HandleFunc("/create", h.createUser)
	router.HandleFunc("/:user_id", h.changeUser) //реализовать методы GET, PUT, DELETE

	router.HandleFunc("/in", h.moveInBook)
	router.HandleFunc("/out", h.moveOutBook)

	//books := router.Group("/")
	//{
	//	books.GET("/", h.getBooks)
	//	books.POST("/new_book", h.createBook)
	//	books.PUT("/:book_id", h.updateBook)
	//	books.DELETE("/:book_id", h.deleteBook)
	//}
	//
	//users := router.Group("/users")
	//{
	//	users.GET("/", h.getUsers)
	//	users.POST("/create", h.createUser)
	//	users.GET("/:user_id", h.getUser)
	//	users.PUT("/:user_id", h.updateUser)
	//	users.DELETE("/:user_id", h.deleteUser)
	//}
	//
	//moveBook := router.Group("/move")
	//{
	//	moveBook.POST("/in", h.moveInBook)
	//	moveBook.POST("/out", h.moveOutBook)
	//}

	return router
}
