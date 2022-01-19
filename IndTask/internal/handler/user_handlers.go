package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getUsers")
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 0 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	sorting := service.SortTypeUser(req.URL.Query().Get("sort"))
	var listUsers []IndTask.UserResponse
	listUsers, pages, err := h.services.AppUser.GetUsers(page, sorting)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listUsers)
	if err != nil {
		h.logger.Errorf("UserHandler: error while marshaling list users:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pages", strconv.Itoa(pages))
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getUsers: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createUser")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.User
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("createUser: error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	validationErrors := validateStruct(h, input)
	if len(validationErrors) != 0 {
		h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
		errors, err := json.Marshal(validationErrors)
		if err != nil {
			h.logger.Errorf("UserHandler: error while marshaling list errors:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("UserHandler: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	user, err := h.services.AppUser.CreateUser(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&user)
	if err != nil {
		h.logger.Errorf("UserHandler: error while marshaling user:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("createUser: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}

func (h *Handler) changeUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeUser")
	userId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || userId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Errorf("Method %s Not Allowed", req.Method)
		http.Error(w, fmt.Sprintf("Method %s Not Allowed", req.Method), 405)
		return
	}
	var input IndTask.User
	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeUser")
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
				h.logger.Errorf("UserHandler: error while marshaling list errors:%s", err)
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
	h.logger.Infof("Method %s, changeGenre", req.Method)
	user, err := h.services.AppUser.ChangeUser(&input, userId, req.Method)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if user != nil {
		output, err := json.Marshal(&user)
		if err != nil {
			h.logger.Errorf("UserHandler: error while marshaling user:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("changeUser: error while writing response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
