package handler

import (
	"encoding/json"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

//handlers for books

func (h *Handler) getBooks(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getBooks")
	CheckMethod(w, req, "GET", h.logger)
	var listBooks []IndTask.Book
	listBooks, err := h.services.AppBook.GetBooks()
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.logger.Error(err.Error())
		return
	}
	var output []byte
	err = json.Unmarshal(output, &listBooks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func (h *Handler) createBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createBook")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.Book
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		logrus.Error(err.Error())
	}
	bookId, err := h.services.AppBook.CreateBook(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(bookId))
	w.WriteHeader(200)
}

func (h *Handler) changeBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeBook")
	bookId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || bookId < 1 {
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		http.Error(w, "Method Not Allowed", 405)
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeBook")
		var input IndTask.Book
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			http.Error(w, err.Error(), 400)
			logrus.Error(err.Error())
		}
		_, err = h.services.AppBook.ChangeBook(&input, bookId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeBook")
		var book *IndTask.Book
		book, err = h.services.AppBook.ChangeBook(nil, bookId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var output []byte
		err = json.Unmarshal(output, book)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(output)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeBook")
		_, err = h.services.AppBook.ChangeBook(nil, bookId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}

}

//handlers for users

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getUsers")
	CheckMethod(w, req, "GET", h.logger)
	var listUsers []IndTask.User
	listUsers, err := h.services.AppUser.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		logrus.Error(err.Error())
		return
	}

	var output []byte
	err = json.Unmarshal(output, &listUsers)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(output)
	if err != nil {
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
		http.Error(w, err.Error(), 400)
		logrus.Error(err.Error())
	}
	userId, err := h.services.AppUser.CreateUser(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(userId))
	w.WriteHeader(200)

}

func (h *Handler) changeUser(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeUser")
	userId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || userId < 1 {
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		http.Error(w, "Method Not Allowed", 405)
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeUser")
		var input IndTask.User
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			http.Error(w, err.Error(), 400)
			logrus.Error(err.Error())
		}
		_, err = h.services.AppUser.ChangeUser(&input, userId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeUser")
		var user *IndTask.User
		user, err = h.services.AppUser.ChangeUser(nil, userId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var output []byte
		err = json.Unmarshal(output, user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(output)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeUser")
		_, err = h.services.AppUser.ChangeUser(nil, userId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}

}

//handlers for movement books

func (h *Handler) moveInBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working moveInBook")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.IssueAct
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		logrus.Error(err.Error())
	}
	issueActId, err := h.services.AppMove.MoveInBook(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(issueActId))
	w.WriteHeader(200)
}

func (h *Handler) getMoveInBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getMoveInBook")
	userId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || userId < 1 {
		http.NotFound(w, req)
		return
	}
	CheckMethod(w, req, "GET", h.logger)
	var issueActs []IndTask.IssueAct
	issueActs, err = h.services.AppMove.GetMoveInBooks(userId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var output []byte
	err = json.Unmarshal(output, &issueActs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) moveOutBook(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working moveOutBook")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.ReturnAct
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		logrus.Error(err.Error())
	}
	returnActId, err := h.services.AppMove.MoveOutBook(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(returnActId))
	w.WriteHeader(200)
}

//handlers for authors

func (h *Handler) getAuthors(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getAuthors")
	CheckMethod(w, req, "GET", h.logger)
	var listAuthors []IndTask.Author
	listAuthors, err := h.services.AppAuthor.GetAuthors()
	if err != nil {
		http.Error(w, err.Error(), 500)
		logrus.Error(err.Error())
		return
	}

	var output []byte
	err = json.Unmarshal(output, &listAuthors)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func (h *Handler) createAuthor(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working createAuthor")
	CheckMethod(w, req, "POST", h.logger)
	var input IndTask.Author
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		logrus.Error(err.Error())
	}
	authorId, err := h.services.AppAuthor.CreateAuthor(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(authorId))
	w.WriteHeader(200)

}

func (h *Handler) changeAuthor(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working changeAuthor")
	authorId, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || authorId < 1 {
		http.NotFound(w, req)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		http.Error(w, "Method Not Allowed", 405)
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeAuthor")
		var input IndTask.Author
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			http.Error(w, err.Error(), 400)
			logrus.Error(err.Error())
		}
		_, err = h.services.AppAuthor.ChangeAuthor(&input, authorId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeAuthor")
		var user *IndTask.Author
		user, err = h.services.AppAuthor.ChangeAuthor(nil, authorId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var output []byte
		err = json.Unmarshal(output, user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(output)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeAuthor")
		_, err = h.services.AppAuthor.ChangeAuthor(nil, authorId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}

}

//handlers for genres

func (h *Handler) getGenres(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Working getGenres")
	CheckMethod(w, req, "GET", h.logger)
	var listGenre []IndTask.Genre
	listGenre, err := h.services.AppGenre.GetGenres()
	if err != nil {
		http.Error(w, err.Error(), 500)
		logrus.Error(err.Error())
		return
	}
	var output []byte
	output, err = json.Marshal(&listGenre)
	if err != nil {
		http.Error(w, err.Error(), 500)
		logrus.Errorf("Marshal error:%s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
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
		http.Error(w, err.Error(), 400)
		logrus.Error(err.Error())
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
	genreId, err := strconv.Atoi(req.URL.Query().Get("genre_id"))
	if err != nil || genreId < 1 {
		http.NotFound(w, req)
		logrus.Errorf("No url request:%s", err)
		return
	}
	if req.Method != "PUT" && req.Method != "GET" && req.Method != "DELETE" {
		http.Error(w, "Method Not Allowed", 405)
	}

	if req.Method == "PUT" {
		h.logger.Info("Method PUT, changeGenre")
		var input IndTask.Genre
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&input); err != nil {
			http.Error(w, err.Error(), 400)
			logrus.Error(err.Error())
		}
		_, err = h.services.AppGenre.ChangeGenre(&input, genreId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "GET" {
		h.logger.Info("Method GET, changeGenre")
		genre, err := h.services.AppGenre.ChangeGenre(nil, genreId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var output []byte
		err = json.Unmarshal(output, genre)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte("jjjkgu"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if req.Method == "DELETE" {
		h.logger.Info("Method DELETE, changeGenre")
		_, err = h.services.AppGenre.ChangeGenre(nil, genreId, req.Method)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}
