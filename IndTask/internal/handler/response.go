package handler

type errorResponce struct {
	Message string `json:"message"`
}

//func newErrorResponse(w http.ResponseWriter, req *http.Request, statusCode int, message string) {
//	logrus.Error(message)
//	req.Response.
//}
