package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/pkg/translate"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//handlers for books

func (h *Handler) getBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getBooks")
	CheckMethod(w, req, "GET", h.logger)
	var listBooks []IndTask.Book
	listBooks, err := h.services.AppBook.GetBooks()
	if err != nil {
		h.logger.Errorf("Error while getting books list from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(listBooks)
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

func (h *Handler) getListBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getListBooks")
	CheckMethod(w, req, "GET", h.logger)
	var listBooks []IndTask.ListBooks
	listBooks, err := h.services.AppBook.GetListBooks()
	if err != nil {
		h.logger.Errorf("Error while getting listBooks list from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(listBooks)
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

func (h *Handler) createBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createBook")
	CheckMethod(w, req, "POST", h.logger)
	if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		h.logger.Errorf("Error while parsing multipart form:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	var input IndTask.Book
	body := bytes.NewBufferString(req.PostFormValue("body"))
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("Error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	if err := InputCoverFoto(w, req, h, &input); err != nil {
		return
	}
	bookId, err := h.services.AppBook.CreateBook(&input)
	if err != nil {
		h.logger.Errorf("Error while creating book in the database:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	header := w.Header()
	header.Add("id", strconv.Itoa(bookId))
}

func (h *Handler) changeBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeBook")
	bookId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || bookId < 1 {
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
		h.logger.Info("Method PUT, changeBook")
		var input IndTask.Book
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
			if err := InputCoverFoto(w, req, h, &input); err != nil {
				return
			}
		}
		_, err = h.services.AppBook.ChangeBook(&input, bookId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while updating book: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeBook")
		book, err := h.services.AppBook.ChangeBook(nil, bookId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while getting book from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(book)
		if err != nil {
			h.logger.Errorf("Can not marshal book:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("Can not write output into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

func (h *Handler) changeListBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeListBooks")
	listBookId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || listBookId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Error("Method Not Allowed")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeListBooks")
		var input IndTask.ListBooks
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request: %s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		_, err = h.services.AppBook.ChangeListBook(&input, listBookId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while updating listBook: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeListBooks")
		book, err := h.services.AppBook.ChangeListBook(nil, listBookId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while getting listBook from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(book)
		if err != nil {
			h.logger.Errorf("Can not marshal listBook:%s", err)
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

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeListBooks")
		_, err = h.services.AppBook.ChangeListBook(nil, listBookId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while deleting listBook from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

//handlers for users

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getUsers")
	CheckMethod(w, req, "GET", h.logger)
	var listUsers []IndTask.User
	listUsers, err := h.services.AppUser.GetUsers()
	if err != nil {
		h.logger.Errorf("Error while getting users list from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listUsers)
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

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createUser")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.User
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("Error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}

	userId, err := h.services.AppUser.CreateUser(&input)
	if err != nil {
		h.logger.Errorf("Error while creating user in the database:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	header := w.Header()
	header.Add("id", strconv.Itoa(userId))
}

func (h *Handler) changeUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeUser")
	userId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || userId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Error("Method Not Allowed")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeUser")
		var input IndTask.User
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request: %s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		_, err = h.services.AppUser.ChangeUser(&input, userId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while updating user: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeUser")
		user, err := h.services.AppUser.ChangeUser(nil, userId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while getting user from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(user)
		if err != nil {
			h.logger.Errorf("Can not marshal user:%s", err)
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

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeUser")
		_, err = h.services.AppUser.ChangeUser(nil, userId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while deleting genre from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

}

//handlers for issue acts

func (h *Handler) getIssueActs(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getIssueActs")
	CheckMethod(w, req, "GET", h.logger)
	var issueActs []IndTask.IssueAct
	issueActs, err := h.services.AppMove.GetIssueActs()
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
	var issueActs []IndTask.IssueAct
	issueActs, err = h.services.AppMove.GetIssueActsByUser(userId)
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
	var returnActs []IndTask.ReturnAct
	returnActs, err := h.services.AppMove.GetReturnActs()
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
		if err := InputFineFoto(w, req, h, &input); err != nil {
			return
		}
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
	var returnActs []IndTask.ReturnAct
	returnActs, err = h.services.AppMove.GetReturnActsByUser(userId)
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
			if err := InputFineFoto(w, req, h, &input); err != nil {
				return
			}
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

//handlers for authors

func (h *Handler) getAuthors(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getAuthors")
	CheckMethod(w, req, "GET", h.logger)
	var listAuthors []IndTask.Author
	listAuthors, err := h.services.AppAuthor.GetAuthors()
	if err != nil {
		h.logger.Errorf("Error while getting authors list from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	var output []byte
	output, err = json.Marshal(listAuthors)
	if err != nil {
		h.logger.Errorf("Marshal error:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("Error while writting response:%s", err)
		http.Error(w, err.Error(), 501)
		return
	}

}

func (h *Handler) createAuthor(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createAuthor")
	CheckMethod(w, req, "POST", h.logger)
	if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		h.logger.Errorf("Error while parsing maxMemoty of the uploaded file:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	var input IndTask.Author
	input.AuthorName = req.PostFormValue("author_name")
	if err := InputAuthorFoto(w, req, h, &input); err != nil {
		return
	}
	authorId, err := h.services.AppAuthor.CreateAuthor(&input)
	if err != nil {
		h.logger.Errorf("Error while creating author in the database:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(authorId))
}

func (h *Handler) changeAuthor(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeAuthor")
	authorId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || authorId < 1 {
		h.logger.Errorf("No url request:%s", err)
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Error("Method Not Allowed")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeAuthor")
		var input IndTask.Author

		if req.Header.Get("Content-Type") == "application/json" {
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&input); err != nil {
				h.logger.Errorf("Error while decoding request: %s", err)
				http.Error(w, err.Error(), 400)
				return
			}
		} else {
			if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
				h.logger.Errorf("Error while setting maxMemoty of the uploaded file:%s", err)
				http.Error(w, err.Error(), 400)
				return
			}
			input.AuthorName = req.PostFormValue("author_name")
			if err := InputAuthorFoto(w, req, h, &input); err != nil {
				return
			}

		}

		_, err = h.services.AppAuthor.ChangeAuthor(&input, authorId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while updating genre: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeAuthor")
		author, err := h.services.AppAuthor.ChangeAuthor(nil, authorId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while getting author from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(author)
		if err != nil {
			h.logger.Errorf("Can not marshal user:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("Can not write output into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeAuthor")
		_, err = h.services.AppAuthor.ChangeAuthor(nil, authorId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while deleting genre from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

//handlers for genres

func (h *Handler) getGenres(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getGenres")
	CheckMethod(w, req, "GET", h.logger)
	var listGenre []IndTask.Genre
	listGenre, err := h.services.AppGenre.GetGenres()
	if err != nil {
		h.logger.Errorf("Error while getting genre list from database: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&listGenre)
	if err != nil {
		h.logger.Errorf("Marshal error:%s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("Error while writting response:%s", err)
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
	GenreId, err := h.services.AppGenre.CreateGenre(&input)
	if err != nil {
		h.logger.Errorf("Error while creating genre in the database:%s", err)
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
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		h.logger.Error("Method Not Allowed")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeGenre")
		var input IndTask.Genre
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			h.logger.Errorf("Error while decoding request: %s", err)
			http.Error(w, err.Error(), 400)
			return
		}
		_, err = h.services.AppGenre.ChangeGenre(&input, genreId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while updating genre: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeGenre")
		genre, err := h.services.AppGenre.ChangeGenre(nil, genreId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while getting genre from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(genre)
		if err != nil {
			h.logger.Errorf("Can not marshal genre:%s", err)
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

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeGenre")
		_, err = h.services.AppGenre.ChangeGenre(nil, genreId, req.Method)
		if err != nil {
			h.logger.Errorf("Error while deleting genre from database:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
}
func InputCoverFoto(w http.ResponseWriter, req *http.Request, h *Handler, input *IndTask.Book) error {
	reqFile, fileHeader, err := req.FormFile("file")
	if err != nil {
		h.logger.Errorf("Error while parsing multipart form:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	defer reqFile.Close()
	filePath := fmt.Sprintf("images/book_covers/%s_%d.%s", translate.Translate(input.BookName), input.Published, (strings.Split(fileHeader.Filename, "."))[1])
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		h.logger.Errorf("Error while opening file:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(reqFile)
	if err != nil {
		h.logger.Errorf("Error while reading request file:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		h.logger.Errorf("Error while writting into file:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	input.Cover = filePath
	return nil
}

func InputAuthorFoto(w http.ResponseWriter, req *http.Request, h *Handler, input *IndTask.Author) error {
	reqFile, fileHeader, err := req.FormFile("file")
	if err != nil {
		h.logger.Errorf("Error while parsing multipart form:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	defer reqFile.Close()
	filePath := fmt.Sprintf("images/authors/%s.%s", translate.Translate(input.AuthorName), (strings.Split(fileHeader.Filename, "."))[1])
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		h.logger.Errorf("Error while opening file:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(reqFile)
	if err != nil {
		h.logger.Errorf("Error while reading request file:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		h.logger.Errorf("Error while writting into file:%s", err)
		http.Error(w, err.Error(), 400)
		return err
	}
	input.AuthorFoto = filePath
	return nil
}

func InputFineFoto(w http.ResponseWriter, req *http.Request, h *Handler, input *IndTask.ReturnAct) error {
	m := req.MultipartForm
	files := m.File["file"]
	for i, headers := range files {
		reqfile, err := files[i].Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		defer reqfile.Close()
		fileBytes, err := ioutil.ReadAll(reqfile)
		if err != nil {
			h.logger.Errorf("Error while reading request file:%s", err)
			http.Error(w, err.Error(), 400)
			return err
		}
		directoryPath := fmt.Sprintf("images/fines/issueActId%d", input.IssueActId)
		filePath := fmt.Sprintf("%s/%s", directoryPath, translate.Translate(headers.Filename))
		err = os.MkdirAll(directoryPath, 0777)
		if err != nil {
			h.logger.Errorf("Error while creating directories:%s", err)
			http.Error(w, err.Error(), 400)
			return err
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			h.logger.Errorf("Error while opening file:%s", err)
			http.Error(w, err.Error(), 400)
			return err
		}
		defer file.Close()
		_, err = file.Write(fileBytes)
		if err != nil {
			h.logger.Errorf("Error while writting into file:%s", err)
			http.Error(w, err.Error(), 400)
			return err
		}
		input.Foto = append(input.Foto, filePath)
	}
	return nil
}
