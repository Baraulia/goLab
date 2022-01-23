package IndTask

import (
	"database/sql/driver"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"strconv"
	"strings"
	"time"
)

var logger logging.Logger

const layout = "2006-01-02"

type MyTime struct {
	time.Time
}

type MyDuration struct {
	time.Duration
}

func (c *MyTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" || string(data) == "" {
		logger.Error("date of birth is not specified")
		return fmt.Errorf("date of birth is not specified")
	} else {
		s := strings.Trim(string(data), "\"")
		// Fractional seconds are handled implicitly by Parse.
		tt, err := time.Parse(layout, s)
		*c = MyTime{tt}
		return err
	}
}

func (c MyTime) Value() (driver.Value, error) {
	return driver.Value(c.Time), nil
}

func (c *MyTime) Scan(src interface{}) error {
	switch t := src.(type) {
	case time.Time:
		c.Time = t
		return nil
	default:
		return fmt.Errorf("column type not supported")
	}
}
func (c MyTime) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(layout))), nil
}

// Value converts Duration to a primitive value ready to written to a database.
func (d MyDuration) Value() (driver.Value, error) {
	return driver.Value(int64(d.Duration.Hours() / 24)), nil
}

// Scan reads a Duration value from database driver type.
func (d *MyDuration) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case int64:
		*d = MyDuration{time.Duration(v * 24 * 1e9 * 60 * 60)}
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Duration from: %#v", v)
	}
	return nil
}

func (d *MyDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, int64(d.Hours()/24))), nil
}

func (d *MyDuration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "" {
		logger.Error("rent duration is not specified")
		return fmt.Errorf("rent duration is not specified")
	} else {
		days, err := strconv.Atoi(string(data))
		if err != nil {
			logger.Errorf("Error converting json data : %s into integer:%s", string(data), err)
			return fmt.Errorf("error converting json data : %s into integer:%s", string(data), err)
		}
		nanosec := days * 1e9 * 24 * 60 * 60
		*d = MyDuration{time.Duration(nanosec)}
		return nil
	}
}

type Book struct {
	Id        int     `json:"id"`
	BookName  string  `json:"book_name" validate:"string,min=1,max=255"`
	GenreId   []int   `json:"genres" validate:"genreExist"`
	Cost      float32 `json:"cost"`
	AuthorsId []int   `json:"authors" validate:"authorExist"`
	Cover     string  `json:"cover"`
	Published int     `json:"published" validate:"number,min=1700,max=2022"`
	Pages     int     `json:"pages" validate:"number,min=5,max=10000"`
	Amount    int     `json:"amount" validate:"number,min=0,max=100"`
}

type MostPopularBook struct {
	Id      int     `json:"id"`
	Cover   string  `json:"cover"`
	Readers int     `json:"readers"`
	Rating  float32 `json:"rating"`
}

type BookResponse struct {
	Id              int     `json:"id"`
	BookName        string  `json:"book_name"`
	Genre           []Genre `json:"genres"`
	Published       int     `json:"published"`
	Number          int     `json:"number"`
	AvailableNumber int     `json:"available_number"`
}

type OneBookResponse struct {
	Id        int      `json:"id"`
	BookName  string   `json:"book_name"`
	Genre     []Genre  `json:"genres"`
	Cost      float32  `json:"cost"`
	Authors   []Author `json:"authors"`
	Cover     string   `json:"cover"`
	Published int      `json:"published"`
	Pages     int      `json:"pages" validate:"number,min=5,max=10000"`
	Amount    int      `json:"amount" validate:"number,min=1,max=100"`
}

type ListBook struct {
	Id         int       `json:"id"`
	BookId     int       `json:"book_id" validate:"bookExist"`
	Issued     bool      `json:"issued"`
	RentNumber int       `json:"rent_number"`
	RentCost   float64   `json:"rent_cost"`
	RegDate    time.Time `json:"reg_date"`
	Condition  int       `json:"condition"`
	Scrapped   bool      `json:"scrapped"`
}

type ListBooksResponse struct {
	Id         int              `json:"id"`
	Book       *OneBookResponse `json:"book"`
	Issued     bool             `json:"issued"`
	RentNumber int              `json:"rent_number"`
	RentCost   float64          `json:"rent_cost"`
	RegDate    time.Time        `json:"reg_date"`
	Condition  int              `json:"condition"`
	Scrapped   bool             `json:"scrapped"`
}

type Author struct {
	Id         int    `json:"id"`
	AuthorName string `json:"author_name" validate:"string,min=2,max=255"`
	AuthorFoto string `json:"author_foto"`
}

type User struct {
	Id         int    `json:"id"`
	Surname    string `json:"surname" validate:"string,min=2,max=255"`
	UserName   string `json:"user_name" validate:"string,min=2,max=255"`
	Patronymic string `json:"patronymic" validate:"string,min=2,max=255"`
	PaspNumber string `json:"pasp_number" validate:"string,min=6,max=50"`
	Email      string `json:"email" validate:"email"`
	Address    string `json:"address" validate:"string,min=2,max=255"`
	BirthDate  MyTime `json:"birth_date" validate:"birthDay"`
}

type UserResponse struct {
	Id        int    `json:"id"`
	Surname   string `json:"surname"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	BirthDate MyTime `json:"birth_date"`
}

type Genre struct {
	Id        int    `json:"id" db:"id"`
	GenreName string `json:"genre_name" db:"genre_name" validate:"string,min=3,max=255"`
}

type ReturnAct struct {
	ActId            int       `json:"act_id" validate:"actExist"`
	UserId           int       `json:"user_id" validate:"userExist"`
	ListBookId       int       `json:"list_book_id" validate:"listBookExist"`
	ReturnDate       time.Time `json:"return_date"`
	Foto             []string  `json:"foto"`
	Fine             float64   `json:"fine"`
	ConditionDecrese int       `json:"condition_decrese" validate:"number,min=1,max=100"`
	Rating           int       `json:"rating" validate:"number,min=0,max=10"`
}

type Act struct {
	Id               int        `json:"id"`
	UserId           int        `json:"user_id" validate:"userExist"`
	ListBookId       int        `json:"list_book_id" validate:"listBookExist"`
	RentalTime       MyDuration `json:"rental_time" validate:"rentalTime"`
	ReturnDate       time.Time  `json:"return_date"`
	PreCost          float64    `json:"pre_cost"`
	Cost             float64    `json:"cost"`
	Status           string     `json:"status"`
	ActualReturnDate time.Time  `json:"actual_return_date"`
	Foto             []string   `json:"foto"`
	Fine             float64    `json:"fine"`
	ConditionDecrese int        `json:"condition_decrese" validate:"number,min=0,max=100"`
	Rating           int        `json:"rating" validate:"number,min=0,max=10"`
}

type Debtor struct {
	Email string
	Name  string
	Book  []string
}
