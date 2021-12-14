package IndTask

import (
	"os"
	"time"
)

type Book struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	GenreName   string    `json:"genre_name"`
	Cost        float32   `json:"cost"`
	Authors     string    `json:"authors"`
	Cover       os.File   `json:"cover"`
	AuthorsFoto []os.File `json:"authors_foto"`
	RentCost    float32   `json:"rent_cost"`
	Published   time.Time `json:"published"`
	RegDate     time.Time `json:"reg_date"`
	Pages       int       `json:"pages"`
	Condition   int       `json:"condition"`
	Rating      []int     `json:"rating"`
}

type User struct {
	Surname    string    `json:"surname"`
	Name       string    `json:"name"`
	Patronymic string    `json:"patronymic"`
	PaspNumber string    `json:"pasp_number"`
	Email      string    `json:"email"`
	Adress     string    `json:"adress"`
	BirthDate  time.Time `json:"birth_date"`
}

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ListBooks struct {
	Id         int  `json:"id"`
	BookId     int  `json:"book_id"`
	Issued     bool `json:"issued"`
	RentNumber int  `json:"rent_number"`
}

type IssueAct struct {
	UserId     int           `json:"user_id"`
	RentalTime time.Duration `json:"rental_time"`
	ReturnDate time.Time     `json:"return_date"`
	PreCost    float32       `json:"pre_cost"`
	BookId     int           `json:"book_id"`
	Status     bool          `json:"status"`
}

type ReturnAct struct {
	UserId           int       `json:"user_id"`
	BookId           int       `json:"book_id"`
	Cost             float32   `json:"cost"`
	ReturnDate       time.Time `json:"return_date"`
	Foto             []os.File `json:"foto"`
	Fine             float32   `json:"fine"`
	ConditionDecrese int       `json:"condition_decrese"`
}
