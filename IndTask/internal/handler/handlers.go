package handler

import (
	"encoding/json"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

//handlers for books

func (h *Handler) getBooks(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) createBook(w http.ResponseWriter, req *http.Request) {
	var input IndTask.Book
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		logrus.Error(err.Error())
	}
	id, err := h.services.AppBook.CreateBook(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(id))
	w.WriteHeader(200)
}

func (h *Handler) changeBook(w http.ResponseWriter, req *http.Request) {
	if req.Method == "PUT" {

	}

}

//handlers for users

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) changeUser(w http.ResponseWriter, req *http.Request) {

}

//handlers for movement books

func (h *Handler) moveInBook(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) moveOutBook(w http.ResponseWriter, req *http.Request) {

}

//handlers for authors

func (h *Handler) getAuthors(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) createAuthor(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) changeAuthor(w http.ResponseWriter, req *http.Request) {

}

//handlers for genres

func (h *Handler) getGenres(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) createGenre(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) changeGenre(w http.ResponseWriter, req *http.Request) {

}
