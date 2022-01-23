package handler

import (
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"net/http"
)

func CheckMethod(w http.ResponseWriter, req *http.Request, allowMethod string, logger logging.Logger) {
	if req.Method != allowMethod {
		logger.Errorf("Method is not allowed.Eexpected method: %s, given method:%s", allowMethod, req.Method)
		http.Error(w, "METHOD IS NOT ALLOWED!!!", 405)
		return
	}
}
