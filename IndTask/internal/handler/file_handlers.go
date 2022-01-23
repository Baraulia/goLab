package handler

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git/internal/myErrors"
	"net/http"
)

func (h *Handler) getFile(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getFile")
	CheckMethod(w, req, "GET", h.logger)
	path := req.URL.Query().Get("file")
	if path == "" {
		h.logger.Error("No url request")
		http.Error(w, fmt.Sprint("No url request"), 400)
		return
	}
	file, err := h.services.AppFile.GetFile(path)
	if err != nil {
		switch e := err.(type) {
		case myErrors.Error:
			http.Error(w, e.Error(), e.Status())
			return
		default:
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
	}
	_, err = w.Write(file)
	if err != nil {
		h.logger.Errorf("getFile: error while writing response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}
}
