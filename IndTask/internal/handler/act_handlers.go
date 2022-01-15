package handler

import (
	"bytes"
	"encoding/json"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) getActs(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getIssueActs")
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	var issueActs []IndTask.IssueAct
	issueActs, err = h.services.AppMove.GetIssueActs(page)
	if err != nil {
		h.logger.Errorf("Error while getting issueActs list from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&issueActs)
	if err != nil {
		h.logger.Errorf("Marshal error:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("Error while writting response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) createIssueAct(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createIssueAct")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.IssueAct
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("Error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	validationErrors := validateStruct(h, input)
	if len(validationErrors) != 0 {
		for _, err := range validationErrors {
			http.Error(w, err, 400)
		}
		h.logger.Error("erroneous data in the request")
		return
	}
	issueActId, err := h.services.AppMove.CreateIssueAct(&input, req.Method)
	if err != nil {
		h.logger.Errorf("Error while creating issueAct in the database:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(issueActId))
}

func (h *Handler) getIssueActsByUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getIssueActsByUser")
	userId, err := strconv.Atoi(req.URL.Query().Get("user_id"))
	if err != nil || userId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	var issueActs []IndTask.IssueAct
	issueActs, err = h.services.AppMove.GetIssueActsByUser(userId, page)
	if err != nil {
		h.logger.Errorf("Error while getting issueActsByUser from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&issueActs)
	if err != nil {
		h.logger.Errorf("Marshal error:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("Error while writting response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) changeIssueAct(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeIssueAct")
	actId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || actId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" {
		h.logger.Error("Method Not Allowed")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeIssueAct")
		var input IndTask.IssueAct
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request: %s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		validationErrors := validateStruct(h, input)
		if len(validationErrors) != 0 {
			for _, err := range validationErrors {
				http.Error(w, err, 400)
			}
			h.logger.Error("erroneous data in the request")
			return
		}
		_, err = h.services.AppMove.ChangeIssueAct(&input, actId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while updating issue act: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeIssueAct")
		issueAct, err := h.services.AppMove.ChangeIssueAct(nil, actId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while getting issueAct id=%d from database:%s", actId, err)
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(issueAct)
		if err != nil {
			h.logger.Errorf("Can not marshal issueAct:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("Can not write output into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

//handlers for return acts

func (h *Handler) getReturnActs(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getReturnActs")
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	var returnActs []IndTask.ReturnAct
	returnActs, err = h.services.AppMove.GetReturnActs(page)
	if err != nil {
		h.logger.Errorf("Error while getting returnActs list from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&returnActs)
	if err != nil {
		h.logger.Errorf("Marshal error:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("Error while writting response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) createReturnAct(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createReturnAct")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.ReturnAct
	if req.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
			h.logger.Errorf("Error while parsing multipart form:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		body := bytes.NewBufferString(req.PostFormValue("body"))
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request:%s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		if err := service.InputFineFoto(req, &input); err != nil {
			return
		}
	}
	validationErrors := validateStruct(h, input)
	if len(validationErrors) != 0 {
		for _, err := range validationErrors {
			http.Error(w, err, 400)
		}
		h.logger.Error("erroneous data in the request")
		return
	}
	returnActId, err := h.services.AppMove.CreateReturnAct(&input)
	if err != nil {
		h.logger.Errorf("Error while creating returnAct in the database:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(returnActId))
}

func (h *Handler) getReturnActsByUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getReturnActsByUser")
	userId, err := strconv.Atoi(req.URL.Query().Get("user_id"))
	if err != nil || userId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	CheckMethod(w, req, "GET", h.logger)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	var returnActs []IndTask.ReturnAct
	returnActs, err = h.services.AppMove.GetReturnActsByUser(userId, page)
	if err != nil {
		h.logger.Errorf("Error while getting returnActsByUser from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&returnActs)
	if err != nil {
		h.logger.Errorf("Marshal error:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("Error while writting response:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) changeReturnAct(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeReturnAct")
	actId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || actId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" {
		h.logger.Error("Method Not Allowed")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeReturnAct")
		var input IndTask.ReturnAct
		if req.Header.Get("Content-Type") == "application/json" {
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("Error while decoding request: %s", err)
				http.Error(w, err.Error(), 400)
				return
			}
		} else {
			body := bytes.NewBufferString(req.PostFormValue("body"))
			decoder := json.NewDecoder(body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("Error while decoding request:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
			if err := service.InputFineFoto(req, &input); err != nil {
				return
			}
		}
		validationErrors := validateStruct(h, input)
		if len(validationErrors) != 0 {
			for _, err := range validationErrors {
				http.Error(w, err, 400)
			}
			h.logger.Error("erroneous data in the request")
			return
		}
		_, err = h.services.AppMove.ChangeReturnAct(&input, actId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while updating return act: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeReturnAct")
		returnAct, err := h.services.AppMove.ChangeReturnAct(nil, actId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while getting returnAct from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(returnAct)
		if err != nil {
			h.logger.Errorf("Can not marshal returnAct:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("Can not write output into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
}
