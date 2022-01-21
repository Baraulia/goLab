package repository

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/lib/pq"
	"time"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

var bookLimit = 10

func (r *BookPostgres) GetThreeBooks() ([]IndTask.MostPopularBook, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetThreeBooks: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getThreeBooks: can not starts transaction:%w", err)
	}
	var listBooks []IndTask.MostPopularBook
	query := "SELECT books.id, books.cover, SUM(list_books.rent_number) AS readers, SUM(act.rating)/COUNT(act.id) AS rating  FROM list_books JOIN act " +
		"ON act.rating>0 AND act.listbook_id=list_books.id " +
		"JOIN books ON books.id=list_books.book_id GROUP BY books.id ORDER BY readers LIMIT 3"
	rows, err := transaction.Query(query)
	if err != nil {
		logger.Errorf("GetThreeBooks: can not executes a query:%s", err)
		return nil, fmt.Errorf("getThreeBooks: repository error:%w", err)
	}
	for rows.Next() {
		var book IndTask.MostPopularBook
		if err := rows.Scan(&book.Id, &book.Cover, &book.Readers, &book.Rating); err != nil {
			logger.Errorf("GetThreeBooks: error while scanning for book:%s", err)
			return nil, fmt.Errorf("getThreeBooks: repository error:%w", err)
		}
		listBooks = append(listBooks, book)
	}
	return listBooks, transaction.Commit()
}

func (r *BookPostgres) GetBooks(page int, sorting string) ([]*IndTask.BookResponse, int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetBooks: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getBooks: can not starts transaction:%w", err)
	}
	var listBooks []*IndTask.BookResponse
	var rows *sql.Rows
	var pages int
	var booksId []int
	var genres = make(map[int][]IndTask.Genre)
	if page == 0 {
		query := fmt.Sprintf("SELECT books.id, books.book_name, books.published, books.amount, count(list_books.id) AS av_books FROM books "+
			"JOIN list_books ON books.id=list_books.book_id AND list_books.issued='false' GROUP BY books.id ORDER BY %s", sorting)
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetBooks: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getBooks: repository error:%w", err)
		}
	} else {
		query := fmt.Sprintf("SELECT books.id, books.book_name, books.published, books.amount, count(list_books.id) AS av_books FROM books "+
			"JOIN list_books ON books.id=list_books.book_id AND list_books.issued='false' GROUP BY books.id ORDER BY %s LIMIT $1 OFFSET $2", sorting)
		rows, err = transaction.Query(query, bookLimit, (page-1)*bookLimit)
		if err != nil {
			logger.Errorf("GetBooks: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getBooks: repository error:%w", err)
		}
	}
	for rows.Next() {
		var book IndTask.BookResponse
		if err := rows.Scan(&book.Id, &book.BookName, &book.Published, &book.Number, &book.AvailableNumber); err != nil {
			logger.Errorf("GetBooks: error while scanning for book:%s", err)
			return nil, 0, fmt.Errorf("getBooks: repository error:%w", err)
		}
		booksId = append(booksId, book.Id)
		listBooks = append(listBooks, &book)
	}
	query := "SELECT book_genre.book_id, genre.id, genre.genre_name FROM genre JOIN book_genre ON genre.id = book_genre.genre_id AND book_genre.book_id = ANY ($1)"
	rows, err = transaction.Query(query, pq.Array(booksId))
	if err != nil {
		logger.Errorf("GetBooks: can not executes a query:%s", err)
		return nil, 0, fmt.Errorf("getBooks: repository error:%w", err)
	}
	for rows.Next() {
		var genre IndTask.Genre
		var bookId int
		if err := rows.Scan(&bookId, &genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("GetBooks: error while scanning for genre:%s", err)
			return nil, 0, fmt.Errorf("getBooks: repository error:%w", err)
		}
		genres[bookId] = append(genres[bookId], genre)
	}
	for _, book := range listBooks {
		book.Genre = genres[book.Id]
	}

	query = "SELECT CEILING(COUNT(id)/$1::float) FROM books"
	row := transaction.QueryRow(query, bookLimit)
	if err := row.Scan(&pages); err != nil {
		logger.Errorf("Error while scanning for pages:%s", err)
		return nil, 0, fmt.Errorf("getBooks: error while scanning for pages:%w", err)
	}
	return listBooks, pages, transaction.Commit()
}

func (r *BookPostgres) CreateBook(book *IndTask.Book, bookRentCost float64) (*IndTask.OneBookResponse, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("createBook: can not starts transaction:%w", err)
	}
	var newBook IndTask.OneBookResponse
	defer transaction.Rollback()
	createBookQuery := "INSERT INTO books (book_name, cost, cover, published, pages, amount) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, book_name, cost, cover, published, pages, amount"
	row := transaction.QueryRow(createBookQuery, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, book.Amount)
	if err := row.Scan(&newBook.Id, &newBook.BookName, &newBook.Cost, &newBook.Cover, &newBook.Published, &newBook.Pages, &newBook.Amount); err != nil {
		logger.Errorf("Error while scanning for book:%s", err)
		return nil, fmt.Errorf("createBook: error while scanning for book:%w", err)
	}
	for _, reqAuthor := range book.AuthorsId {
		var author IndTask.Author
		createBookAuthorQuery := "INSERT INTO book_author (book_id, author_id) VALUES ($1, $2)"
		_, err = transaction.Exec(createBookAuthorQuery, newBook.Id, reqAuthor)
		if err != nil {
			logger.Errorf("Error while execution query for insert into book_author:%s", err)
			return nil, fmt.Errorf("CreateBook: error while execution query for insert into book_author:%w", err)
		}
		query := "SELECT id, author_name, author_foto FROM authors WHERE id=$1"
		row := transaction.QueryRow(query, reqAuthor)
		if err != nil {
			logger.Errorf("CreateBook: can not executes a query:%s", err)
			return nil, fmt.Errorf("createBook: repository error:%w", err)
		}
		if err := row.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("CreateBook: error while scanning for author:%s", err)
			return nil, fmt.Errorf("createBook: repository error:%w", err)
		}
		newBook.Authors = append(newBook.Authors, author)
	}
	for _, reqGenre := range book.GenreId {
		var genre IndTask.Genre
		createBookAuthorQuery := "INSERT INTO book_genre (book_id, genre_id) VALUES ($1, $2)"
		_, err = transaction.Exec(createBookAuthorQuery, newBook.Id, reqGenre)
		if err != nil {
			logger.Errorf("Error while execution query for insert into book_genre:%s", err)
			return nil, fmt.Errorf("CreateBook: error while execution query for insert into book_genre:%w", err)
		}
		query := "SELECT id, genre_name FROM genre WHERE id=$1"
		row := transaction.QueryRow(query, reqGenre)
		if err != nil {
			logger.Errorf("CreateBook: can not executes a query:%s", err)
			return nil, fmt.Errorf("createBook: repository error:%w", err)
		}
		if err := row.Scan(&genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("CreateBook: error while scanning for author:%s", err)
			return nil, fmt.Errorf("createBook: repository error:%w", err)
		}
		newBook.Genre = append(newBook.Genre, genre)
	}
	for i := 0; i < book.Amount; i++ {
		createListBookQuery := "INSERT INTO list_books (book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err := transaction.Exec(createListBookQuery, newBook.Id, "false", 0, bookRentCost, time.Now(), 100, false)
		if err != nil {
			logger.Errorf("Error while execution query for insert into list_books:%s", err)
			return nil, fmt.Errorf("CreateBook: error while execution query for insert into list_books:%w", err)
		}
	}
	return &newBook, transaction.Commit()
}

func (r *BookPostgres) GetOneBook(bookId int) (*IndTask.OneBookResponse, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneBook: can not starts transaction:%w", err)
	}
	var book IndTask.OneBookResponse
	query := "SELECT id, book_name, cost, cover, published, pages, amount FROM books WHERE id = $1"
	row := transaction.QueryRow(query, bookId)
	if err := row.Scan(&book.Id, &book.BookName, &book.Cost, &book.Cover, &book.Published, &book.Pages, &book.Amount); err != nil {
		logger.Errorf("GetOneBook: error while scanning for book:%s", err)
		return nil, fmt.Errorf("getOneBook: repository error:%w", err)
	}
	book.Authors, err = r.ReturnBindAuthors(book.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting bound authors:%w", err)
	}
	book.Genre, err = r.ReturnBindGenres(book.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting bound genres:%w", err)
	}
	return &book, transaction.Commit()
}

func (r *BookPostgres) ChangeBook(book *IndTask.Book, bookId int, bookRentCost float64) (*IndTask.OneBookResponse, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("changeBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var query string
	var upBook IndTask.OneBookResponse
	var row *sql.Row
	if book.Cover == "" {
		query = "UPDATE books SET book_name=$1, cost=$2, published=$3, pages=$4 WHERE id = $5 RETURNING id, book_name, cost, cover, published, pages, amount"
		row = transaction.QueryRow(query, book.BookName, book.Cost, book.Published, book.Pages, bookId)
	} else {
		query = "UPDATE books SET book_name=$1, cost=$2, cover=$3, published=$4, pages=$5 WHERE id = $6 RETURNING id, book_name, cost, cover, published, pages, amount"
		row = transaction.QueryRow(query, book.BookName, book.Cost, book.Cover, book.Published, book.Pages, bookId)
	}
	if err := row.Scan(&upBook.Id, &upBook.BookName, &upBook.Cost, &upBook.Cover, &upBook.Published, &upBook.Pages, &upBook.Amount); err != nil {
		logger.Errorf("ChangeBook: error while scanning for book:%s", err)
		return nil, fmt.Errorf("changeBook: repository error:%w", err)
	}
	query = "UPDATE list_books SET rent_cost = $1 WHERE book_id=$2"
	_, err = transaction.Exec(query, bookRentCost, bookId)
	if err != nil {
		logger.Errorf("Repository error while updating rent_cost in list_books:%s", err)
		return nil, fmt.Errorf("changeBook: repository error:%w", err)
	}
	upBook.Authors, err = r.ReturnBindAuthors(upBook.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting bound authors:%w", err)
	}
	upBook.Genre, err = r.ReturnBindGenres(upBook.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting bound genres:%w", err)
	}
	return &upBook, transaction.Commit()
}

func (r *BookPostgres) DeleteBook(bookId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeBook: can not starts transaction:%s", err)
		return fmt.Errorf("changeBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM books WHERE id = $1"
	_, err = transaction.Exec(query, bookId)
	if err != nil {
		logger.Errorf("Repository error while deleting book:%s", err)
		return fmt.Errorf("deleteBook: repository error:%w", err)
	}
	query = "DELETE FROM list_books WHERE book_id = $1"
	_, err = transaction.Exec(query, bookId)
	if err != nil {
		logger.Errorf("Repository error while deleting listBook:%s", err)
		return fmt.Errorf("deleteBook: repository error:%w", err)
	}
	return transaction.Commit()
}

func (r *BookPostgres) GetListBooks(page int) ([]IndTask.ListBooksResponse, int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetListBooks: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getListBooks: can not starts transaction:%w", err)
	}
	var listBooks []IndTask.ListBooksResponse
	var rows *sql.Rows
	var booksId []int
	var genres = make(map[int][]IndTask.Genre)
	var authors = make(map[int][]IndTask.Author)
	var pages int
	if page == 0 {
		query := "SELECT list_books.id, list_books.book_id, list_books.issued, list_books.rent_number, list_books.rent_cost, list_books.reg_date, list_books.condition, list_books.scrapped, " +
			"books.id, books.book_name, books.cost, books.cover, books.published, books.pages, books.amount FROM list_books JOIN books ON list_books.book_id = books.id " +
			"WHERE list_books.issued='false' AND list_books.scrapped='false' ORDER BY list_books.Id"
		rows, err = transaction.Query(query)
		if err != nil {
			logger.Errorf("GetListBooks: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getListBooks: repository error:%w", err)
		}
	} else {
		query := "SELECT list_books.id, list_books.book_id, list_books.issued, list_books.rent_number, list_books.rent_cost, list_books.reg_date, list_books.condition, list_books.scrapped, " +
			"books.id, books.book_name, books.cost, books.cover, books.published, books.pages, books.amount FROM list_books JOIN books ON list_books.book_id = books.id " +
			"WHERE list_books.issued='false' AND list_books.scrapped='false' ORDER BY list_books.Id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, bookLimit, (page-1)*bookLimit)
		if err != nil {
			logger.Errorf("GetListBooks: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getListBooks: repository error:%w", err)
		}
	}
	for rows.Next() {
		var book IndTask.ListBooksResponse
		var unBook IndTask.OneBookResponse
		var bookId int
		if err := rows.Scan(&book.Id, &bookId, &book.Issued, &book.RentNumber, &book.RentCost, &book.RegDate, &book.Condition, &book.Scrapped,
			&unBook.Id, &unBook.BookName, &unBook.Cost, &unBook.Cover, &unBook.Published, &unBook.Pages, &unBook.Amount); err != nil {
			logger.Errorf("GetListBooks: error while scanning for book:%s", err)
			return nil, 0, fmt.Errorf("getListBooks: repository error:%w", err)
		}
		booksId = append(booksId, bookId)
		book.Book = &unBook
		listBooks = append(listBooks, book)
	}
	query := "SELECT book_genre.book_id, genre.id, genre.genre_name FROM genre JOIN book_genre ON genre.id = book_genre.genre_id AND book_genre.book_id = ANY ($1)"
	rows, err = transaction.Query(query, pq.Array(booksId))
	if err != nil {
		logger.Errorf("GetBooks: can not executes a query:%s", err)
		return nil, 0, fmt.Errorf("getBooks: repository error:%w", err)
	}
	for rows.Next() {
		var genre IndTask.Genre
		var bookId int
		if err := rows.Scan(&bookId, &genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("GetBooks: error while scanning for genre:%s", err)
			return nil, 0, fmt.Errorf("getBooks: repository error:%w", err)
		}
		genres[bookId] = append(genres[bookId], genre)
	}
	query = "SELECT book_author.book_id, authors.id, authors.author_name, authors.author_foto FROM authors JOIN book_author ON authors.id = book_author.author_id AND book_author.book_id = ANY ($1)"
	rows, err = transaction.Query(query, pq.Array(booksId))
	if err != nil {
		logger.Errorf("GetListBooks: can not executes a query:%s", err)
		return nil, 0, fmt.Errorf("getListBooks: repository error:%w", err)
	}
	for rows.Next() {
		var author IndTask.Author
		var authorId int
		if err := rows.Scan(&authorId, &author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("GetListBooks: error while scanning for author:%s", err)
			return nil, 0, fmt.Errorf("getListBooks: repository error:%w", err)
		}
		authors[authorId] = append(authors[authorId], author)
	}
	for _, book := range listBooks {
		book.Book.Genre = genres[book.Book.Id]
		book.Book.Authors = authors[book.Book.Id]
	}
	query = "SELECT CEILING(COUNT(id)/$1::float) FROM list_books"
	row := transaction.QueryRow(query, bookLimit)
	if err := row.Scan(&pages); err != nil {
		logger.Errorf("Error while scanning for pages:%s", err)
		return nil, 0, fmt.Errorf("getListBooks: error while scanning for pages:%w", err)
	}
	return listBooks, pages, transaction.Commit()
}

func (r *BookPostgres) GetOneListBook(listBookId int) (*IndTask.ListBooksResponse, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetOneListBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getOneListBook: can not starts transaction:%w", err)
	}
	var listBook IndTask.ListBooksResponse
	var bookId int
	query := "SELECT id, book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped FROM list_books WHERE id = $1"
	row := transaction.QueryRow(query, listBookId)
	if err := row.Scan(&listBook.Id, &bookId, &listBook.Issued, &listBook.RentNumber, &listBook.RentCost, &listBook.RegDate, &listBook.Condition, &listBook.Scrapped); err != nil {
		logger.Errorf("GetOneListBook: error while scanning for listBook:%s", err)
		return nil, fmt.Errorf("getOneListBook: repository error:%w", err)
	}
	listBook.Book, err = r.GetOneBook(bookId)
	if err != nil {
		logger.Errorf("GetOneListBook: error while getting book:%s", err)
		return nil, fmt.Errorf("getOneListBook: repository error:%w", err)
	}
	return &listBook, transaction.Commit()
}

func (r *BookPostgres) CreateListBook(bookId int, bookRentCost float64) (*IndTask.ListBooksResponse, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("CreateListBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("createListBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var listBook IndTask.ListBooksResponse
	createListBookQuery := "INSERT INTO list_books (book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped) VALUES ($1, $2, $3, $4, $5, $6, $7)" +
		"RETURNING id, book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped"
	row := transaction.QueryRow(createListBookQuery, bookId, "false", 0, bookRentCost, time.Now(), 100, false)
	if err := row.Scan(&listBook.Id, &bookId, &listBook.Issued, &listBook.RentNumber, &listBook.RentCost, &listBook.RegDate, &listBook.Condition, &listBook.Scrapped); err != nil {
		logger.Errorf("CreateListBook: error while scanning for listBook:%s", err)
		return nil, fmt.Errorf("createListBook: repository error:%w", err)
	}
	query := "UPDATE books SET amount=amount+1 WHERE id = $1"
	_, err = transaction.Exec(query, bookId)
	if err != nil {
		logger.Errorf("Error while updating books.amount:%s", err)
		return nil, fmt.Errorf("createListBook: Error while updating books.amount:%w", err)
	}
	listBook.Book, err = r.GetOneBook(bookId)
	if err != nil {
		logger.Errorf("CreateListBook: error while getting book:%s", err)
		return nil, fmt.Errorf("createListBook: repository error:%w", err)
	}
	listBook.Book.Amount = listBook.Book.Amount + 1
	return &listBook, transaction.Commit()
}

func (r *BookPostgres) ChangeListBook(listBook *IndTask.ListBook, listBookId int) (*IndTask.ListBooksResponse, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ChangeListBook: can not starts transaction:%s", err)
		return nil, fmt.Errorf("changeListBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var upBook IndTask.ListBooksResponse
	var bookId int
	query := "UPDATE list_books SET issued=$1, rent_cost=$2, condition=$3, scrapped=$4 WHERE id = $5 " +
		"RETURNING id, book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped"
	row := transaction.QueryRow(query, listBook.Issued, listBook.RentCost, listBook.Condition, listBook.Scrapped, listBookId)
	if err := row.Scan(&upBook.Id, &bookId, &upBook.Issued, &upBook.RentNumber, &upBook.RentCost, &upBook.RegDate, &upBook.Condition, &upBook.Scrapped); err != nil {
		logger.Errorf("ChangeListBook: error while scanning for listBook:%s", err)
		return nil, fmt.Errorf("changeListBook: repository error:%w", err)
	}
	upBook.Book, err = r.GetOneBook(bookId)
	if err != nil {
		logger.Errorf("ChangeListBook: error while scanning for listBook:%s", err)
		return nil, fmt.Errorf("changeListBook: repository error:%w", err)
	}
	return &upBook, transaction.Commit()
}

func (r *BookPostgres) DeleteListBook(listBookId int) error {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("DeleteListBook: can not starts transaction:%s", err)
		return fmt.Errorf("deleteListBook: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM list_books WHERE id = $1"
	_, err = transaction.Exec(query, listBookId)
	if err != nil {
		logger.Errorf("Repository error while deleting listBook:%s", err)
		return fmt.Errorf("deleteListBook: repository error:%w", err)
	}
	query = "UPDATE books SET amount=amount-1 WHERE id = (SELECT book_id FROM list_books WHERE id = $1)"
	_, err = transaction.Exec(query, listBookId)
	if err != nil {
		logger.Errorf("Error while updating books.amount:%s", err)
		return fmt.Errorf("deleteListBook: Error while updating books.amount:%w", err)
	}
	return transaction.Commit()
}

func (r *BookPostgres) GetAuthorsByBookId(bookId int) ([]int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetAuthorsByBookId: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getAuthorsByBookId: can not starts transaction:%s", err)
	}
	var authorsId []int
	query := "SELECT author_id FROM book_author WHERE book_id=$1"
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		logger.Errorf("GetAuthorsByBookId: can not executes a query:%s", err)
		return nil, fmt.Errorf("getAuthorsByBookId: repository error:%w", err)
	}
	for rows.Next() {
		var authorId int
		if err := rows.Scan(&authorId); err != nil {
			logger.Errorf("GetAuthorsByBookId: error while scanning for authorId:%s", err)
			return nil, fmt.Errorf("getAuthorsByBookId: repository error:%w", err)
		}
		authorsId = append(authorsId, authorId)
	}
	return authorsId, transaction.Commit()
}

func (r *BookPostgres) ReturnBindAuthors(bookId int) ([]IndTask.Author, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ReturnBindAuthors: can not starts transaction:%s", err)
		return nil, fmt.Errorf("returnBindAuthors: can not starts transaction:%s", err)
	}
	var authors []IndTask.Author
	query := "SELECT id, author_name, author_foto FROM authors JOIN book_author ON authors.id = book_author.author_id AND book_author.book_id = $1"
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		logger.Errorf("ReturnBindAuthors: can not executes a query:%s", err)
		return nil, fmt.Errorf("returnBindAuthors: repository error:%w", err)
	}
	for rows.Next() {
		var author IndTask.Author
		if err := rows.Scan(&author.Id, &author.AuthorName, &author.AuthorFoto); err != nil {
			logger.Errorf("ReturnBindAuthors: error while scanning for author:%s", err)
			return nil, fmt.Errorf("returnBindAuthors: repository error:%w", err)
		}
		authors = append(authors, author)
	}
	return authors, transaction.Commit()
}

func (r *BookPostgres) ReturnBindGenres(bookId int) ([]IndTask.Genre, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("ReturnBindGenres: can not starts transaction:%s", err)
		return nil, fmt.Errorf("returnBindGenres: can not starts transaction:%s", err)
	}
	var genres []IndTask.Genre
	query := "SELECT id, genre_name FROM genre JOIN book_genre ON genre.id = book_genre.genre_id AND book_genre.book_id = $1"
	rows, err := transaction.Query(query, bookId)
	if err != nil {
		logger.Errorf("ReturnBindGenres: can not executes a query:%s", err)
		return nil, fmt.Errorf("returnBindGenres: repository error:%w", err)
	}
	for rows.Next() {
		var genre IndTask.Genre
		if err := rows.Scan(&genre.Id, &genre.GenreName); err != nil {
			logger.Errorf("ReturnBindGenres: error while scanning for genre:%s", err)
			return nil, fmt.Errorf("returnBindGenres: repository error:%w", err)
		}
		genres = append(genres, genre)
	}
	return genres, transaction.Commit()
}

func (r *BookPostgres) GetBooksId() ([]int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetListBooksId: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getListBooksId: can not starts transaction:%w", err)
	}
	var BooksId []int
	query := "SELECT id FROM books"
	rows, err := transaction.Query(query)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.Errorf("GetBooksId: error while scanning for book:%s", err)
			return nil, fmt.Errorf("getBooksId: repository error:%w", err)
		}
		BooksId = append(BooksId, id)
	}
	return BooksId, transaction.Commit()
}

func (r *BookPostgres) GetListBooksId() ([]int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		logger.Errorf("GetListBooksId: can not starts transaction:%s", err)
		return nil, fmt.Errorf("getListBooksId: can not starts transaction:%w", err)
	}
	var listBooksId []int
	query := "SELECT id FROM list_books"
	rows, err := transaction.Query(query)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.Errorf("GetListBooksId: error while scanning for book:%s", err)
			return nil, fmt.Errorf("getListBooksId: repository error:%w", err)
		}
		listBooksId = append(listBooksId, id)
	}
	return listBooksId, transaction.Commit()
}
