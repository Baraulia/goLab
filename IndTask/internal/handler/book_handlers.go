package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) getThreeBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getThreeBooks")
	CheckMethod(w, req, "GET", h.logger)
	var listBooks []IndTask.MostPopularBook
	listBooks, err := h.services.AppBook.GetThreeBooks()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listBooks)
	if err != nil {
		h.logger.Errorf("BookHandler: error while marshaling three books:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getThreeBooks: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) getBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getBooks")
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 0 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	sorting := service.SortTypeBook(req.URL.Query().Get("sort"))
	CheckMethod(w, req, "GET", h.logger)
	var listBooks []*IndTask.BookResponse
	listBooks, pages, err := h.services.AppBook.GetBooks(page, sorting)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listBooks)
	if err != nil {
		h.logger.Errorf("BookHandler: error while marshaling list books:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pages", strconv.Itoa(pages))
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getBooks: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) getListBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getListBooks")
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 0 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	var listBooks []IndTask.ListBooksResponse
	listBooks, pages, err := h.services.AppBook.GetListBooks(page)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listBooks)
	if err != nil {
		h.logger.Errorf("BookHandler: error while marshaling list instances of books:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pages", strconv.Itoa(pages))
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getListBooks: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) createBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createBook")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.Book
	if req.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("createBook: error while decoding request:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
			h.logger.Errorf("createBook: error while parsing multipart form:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		body := bytes.NewBufferString(req.PostFormValue("body"))
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("createBook: error while decoding request:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		if err := h.services.AppBook.InputCoverFoto(req, &input); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	validationErrors := validateStruct(h, input)
	if len(validationErrors) != 0 {
		h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
		errors, err := json.Marshal(validationErrors)
		if err != nil {
			h.logger.Errorf("BookHandler: error while marshaling list errors:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("BookHandler: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	book, err := h.services.AppBook.CreateBook(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&book)
	if err != nil {
		h.logger.Errorf("BookHandler: error while marshaling book:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("createBook: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}

func (h *Handler) changeBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeBook")
	bookId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || bookId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Errorf("Method %s Not Allowed", req.Method)
		http.Error(w, fmt.Sprintf("Method %s Not Allowed", req.Method), 405)
		return
	}
	var input IndTask.Book
	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeBook")
		if req.Header.Get("Content-Type") == "application/json" {
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("changeBook: error while decoding request:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
		} else {
			if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
				h.logger.Errorf("createBook: error while parsing multipart form:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
			body := bytes.NewBufferString(req.PostFormValue("body"))
			decoder := json.NewDecoder(body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("changeBook: error while decoding request:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
			if err := h.services.AppBook.InputCoverFoto(req, &input); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		validationErrors := validateStruct(h, input)
		if len(validationErrors) != 0 {
			h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
			errors, err := json.Marshal(validationErrors)
			if err != nil {
				h.logger.Errorf("BookHandler: error while marshaling list errors:%s", err)
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(errors)
			if err != nil {
				h.logger.Errorf("Can not write errors into response:%s", err)
				http.Error(w, err.Error(), 500)
				return
			}
			return
		}
	}
	h.logger.Infof("Method %s, changeBook", req.Method)
	book, err := h.services.AppBook.ChangeBook(&input, bookId, req.Method)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if book != nil {
		output, err := json.Marshal(&book)
		if err != nil {
			h.logger.Errorf("BookHandler: error while marshaling book:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("changeBook: error while writing response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) createListBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createListBook")
	CheckMethod(w, req, "POST", h.logger)
	var input struct {
		BookId int `json:"book_id"`
	}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("createListBook: error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	book, err := h.services.AppBook.CreateListBook(input.BookId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&book)
	if err != nil {
		h.logger.Errorf("BookHandler: error while marshaling book:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("createListBooks: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) changeListBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeListBooks")
	listBookId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || listBookId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Errorf("Method %s Not Allowed", req.Method)
		http.Error(w, fmt.Sprintf("Method %s Not Allowed", req.Method), 405)
		return
	}
	var input IndTask.ListBook
	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeListBooks")
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request(%s): %s", req.Body, err)
			http.Error(w, err.Error(), 400)
			return
		}
		input.BookId = 1
		validationErrors := validateStruct(h, input)
		if len(validationErrors) != 0 {
			h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
			errors, err := json.Marshal(validationErrors)
			if err != nil {
				h.logger.Errorf("BookHandler: error while marshaling list errors:%s", err)
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(errors)
			if err != nil {
				h.logger.Errorf("Can not write errors into response:%s", err)
				http.Error(w, err.Error(), 500)
				return
			}
			return
		}
	}
	h.logger.Infof("Method %s, changeListBooks", req.Method)
	book, err := h.services.AppBook.ChangeListBook(&input, listBookId, req.Method)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if book != nil {
		output, err := json.Marshal(&book)
		if err != nil {
			h.logger.Errorf("BookHandler: error while marshaling book:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("changeListBooks: error while writing response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
