package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) getAuthors(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getAuthors")
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 0 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	var listAuthors []IndTask.Author
	listAuthors, pages, err := h.services.AppAuthor.GetAuthors(page)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listAuthors)
	if err != nil {
		h.logger.Errorf("AuthorHandler: error while marshaling list authors:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pages", strconv.Itoa(pages))
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getAuthors: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}

func (h *Handler) createAuthor(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createAuthor")
	CheckMethod(w, req, "POST", h.logger)
	if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		h.logger.Errorf("createAuthor: error while parsing multipart form:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	var input IndTask.Author
	input.AuthorName = req.PostFormValue("author_name")
	if err := service.InputAuthorFoto(req, &input); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	validationErrors := validateStruct(h, input)
	if len(validationErrors) != 0 {
		h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
		errors, err := json.Marshal(validationErrors)
		if err != nil {
			h.logger.Errorf("AuthorHandler: error while marshaling list errors:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("AuthorHandler: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	author, err := h.services.AppAuthor.CreateAuthor(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&author)
	if err != nil {
		h.logger.Errorf("AuthorHandler: error while marshaling author:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("createAuthor: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}

func (h *Handler) changeAuthor(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeAuthor")
	authorId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || authorId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Errorf("Method %s Not Allowed", req.Method)
		http.Error(w, fmt.Sprintf("Method %s Not Allowed", req.Method), 405)
		return
	}
	var input IndTask.Author
	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeAuthor")
		if req.Header.Get("Content-Type") == "application/json" {
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("Error while decoding request(%s): %s", req.Body, err)
				http.Error(w, err.Error(), 400)
				return
			}
		} else {
			if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
				h.logger.Errorf("changeAuthor: error while parsing multipart form:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
			input.AuthorName = req.PostFormValue("author_name")
			if err := service.InputAuthorFoto(req, &input); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		validationErrors := validateStruct(h, input)
		if len(validationErrors) != 0 {
			h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
			errors, err := json.Marshal(validationErrors)
			if err != nil {
				h.logger.Errorf("GenreHandler: error while marshaling list errors:%s", err)
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
	h.logger.Infof("Method %s, changeAuthor", req.Method)
	author, err := h.services.AppAuthor.ChangeAuthor(&input, authorId, req.Method)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if author != nil {
		output, err := json.Marshal(&author)
		if err != nil {
			h.logger.Errorf("AuthorHandler: error while marshaling author:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("changeAuthor: error while writing response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
