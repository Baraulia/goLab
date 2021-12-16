package handler

import (
	"net/http"
)

//handlers for books

func (h *Handler) getBooks(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) createBook(w http.ResponseWriter, req *http.Request) {

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
