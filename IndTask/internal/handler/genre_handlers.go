package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"net/http"
	"strconv"
)

func (h *Handler) getGenres(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getGenres")
	CheckMethod(w, req, "GET", h.logger)
	var listGenre []IndTask.Genre
	listGenre, err := h.services.AppGenre.GetGenres()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listGenre)
	if err != nil {
		h.logger.Errorf("GenreHandler: error while marshaling list genres:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getGenres: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}

func (h *Handler) createGenre(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createGenre")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.Genre
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("Error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
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
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("GenreHandler: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	GenreId, err := h.services.AppGenre.CreateGenre(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(GenreId))
}

func (h *Handler) changeGenre(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeGenre")
	genreId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || genreId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Errorf("Method %s Not Allowed", req.Method)
		http.Error(w, fmt.Sprintf("Method %s Not Allowed", req.Method), 405)
		return
	}
	var input IndTask.Genre
	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeGenre")
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request(%s): %s", req.Body, err)
			http.Error(w, err.Error(), 400)
			return
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
	h.logger.Infof("Method %s, changeGenre", req.Method)
	genre, err := h.services.AppGenre.ChangeGenre(&input, genreId, req.Method)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if genre != nil {
		output, err := json.Marshal(genre)
		if err != nil {
			h.logger.Errorf("GenreHandler: error while marshaling genre:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("changeGenre: error while writing response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
