package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	books := router.Group("/")
	{
		books.GET("/", h.getBooks)
		books.POST("/new_book", h.createBook)
		books.PUT("/:book_id", h.updateBook)
		books.DELETE("/:book_id", h.deleteBook)
	}

	users := router.Group("/users")
	{
		users.GET("/", h.getUsers)
		users.POST("/create", h.createUser)
		users.GET("/:user_id", h.getUser)
		users.PUT("/:user_id", h.updateUser)
		users.DELETE("/:user_id", h.deleteUser)
	}

	moveBook := router.Group("/move")
	{
		moveBook.POST("/in", h.moveInBook)
		moveBook.POST("/out", h.moveOutBook)
	}

	return router
}
