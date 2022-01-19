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

func (h *Handler) getActs(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getActs")
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 0 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	var acts []IndTask.Act
	acts, pages, err := h.services.AppAct.GetActs(page)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&acts)
	if err != nil {
		h.logger.Errorf("ActHandler: error while marshaling list acts:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pages", strconv.Itoa(pages))
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getActs: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}

func (h *Handler) createIssueAct(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createIssueAct")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.Act
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("createIssueAct: error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	validationErrors := validateStruct(h, input)
	if len(validationErrors) != 0 {
		h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
		errors, err := json.Marshal(validationErrors)
		if err != nil {
			h.logger.Errorf("ActHandler: error while marshaling list errors:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("ActHandler: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	act, err := h.services.AppAct.CreateIssueAct(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&act)
	if err != nil {
		h.logger.Errorf("ActHandler: error while marshaling act:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("createIssueAct: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}

func (h *Handler) getActsByUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getActsByUser")
	CheckMethod(w, req, "GET", h.logger)
	userId, err := strconv.Atoi(req.URL.Query().Get("user_id"))
	if err != nil || userId < 0 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.Error(w, fmt.Sprintf("No url request:%s", err), 400)
		return
	}
	var acts []IndTask.Act
	acts, pages, err := h.services.AppAct.GetActsByUser(userId, page)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&acts)
	if err != nil {
		h.logger.Errorf("ActHandler: error while marshaling list users:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pages", strconv.Itoa(pages))
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getActsByUser: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) changeAct(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeAct")
	actId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || actId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" {
		h.logger.Errorf("Method %s Not Allowed", req.Method)
		http.Error(w, fmt.Sprintf("Method %s Not Allowed", req.Method), 405)
		return
	}
	var input IndTask.Act
	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeIssueAct")
		if req.Header.Get("Content-Type") == "application/json" {
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("Error while decoding request(%s): %s", req.Body, err)
				http.Error(w, err.Error(), 400)
				return
			}
		} else {
			if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
				h.logger.Errorf("changeAct: error while parsing multipart form:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
			body := bytes.NewBufferString(req.PostFormValue("body"))
			decoder := json.NewDecoder(body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("changeAct: error while decoding request:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
			photos, err := service.InputFineFoto(req, input.Id)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			input.Foto = photos
		}
		validationErrors := validateStruct(h, input)
		if len(validationErrors) != 0 {
			h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
			errors, err := json.Marshal(validationErrors)
			if err != nil {
				h.logger.Errorf("ActHandler: error while marshaling list errors:%s", err)
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(errors)
			if err != nil {
				h.logger.Errorf("ActHandler: can not write errors into response:%s", err)
				http.Error(w, err.Error(), 500)
				return
			}
			return
		}
	}
	h.logger.Infof("Method %s, changeAct", req.Method)
	oneAct, err := h.services.AppAct.ChangeAct(&input, actId, req.Method)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&oneAct)
	if err != nil {
		h.logger.Errorf("ActHandler: error while marshaling act:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("changeAct: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) addReturnAct(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createReturnAct")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.ReturnAct
	if req.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("addReturnAct: error while decoding request:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
			h.logger.Errorf("addReturnAct: error while parsing multipart form:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		body := bytes.NewBufferString(req.PostFormValue("body"))
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("addReturnAct: error while decoding request:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		photos, err := service.InputFineFoto(req, input.ActId)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		input.Foto = photos
	}
	validationErrors := validateStruct(h, input)
	if len(validationErrors) != 0 {
		h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
		errors, err := json.Marshal(validationErrors)
		if err != nil {
			h.logger.Errorf("ActHandler: error while marshaling list errors:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("ActHandler: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	act, err := h.services.AppAct.AddReturnAct(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(act)
	if err != nil {
		h.logger.Errorf("ActHandler: error while marshaling act:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("addReturnAct: error while writing response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}
