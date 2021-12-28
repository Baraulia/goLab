package IndTask

import (
	"database/sql/driver"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"strings"
	"time"
)

var logger logging.Logger

const layout = "2006-01-02"

type MyTime struct {
	time.Time
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

type Book struct {
	Id        int     `json:"id"`
	BookName  string  `json:"book_name"`
	GenreId   []int   `json:"genre_id"`
	Cost      float32 `json:"cost"`
	AuthorsId []int   `json:"authors"`
	Cover     string  `json:"cover"`
	Published int     `json:"published"`
	Pages     int     `json:"pages"`
	Amount    int     `json:"amount"`
	RentCost  int     `json:"rent_cost"`
}

type ListBooks struct {
	Id         int     `json:"id"`
	BookId     int     `json:"book_id"`
	Issued     bool    `json:"issued"`
	RentNumber int     `json:"rent_number"`
	RentCost   float32 `json:"rent_cost"`
	Condition  int     `json:"condition"`
}

type Author struct {
	Id         int    `json:"id"`
	AuthorName string `json:"author_name"`
	AuthorFoto string `json:"author_foto"`
}

type User struct {
	Id         int     `json:"id"`
	Surname    string  `json:"surname" binding:"required"`
	UserName   string  `json:"user_name" binding:"required"`
	Patronymic string  `json:"patronymic"`
	PaspNumber string  `json:"pasp_number"`
	Email      string  `json:"email"`
	Adress     string  `json:"adress"`
	BirthDate  *MyTime `json:"birth_date" binding:"required"`
}

type Genre struct {
	Id        int    `json:"id" db:"id"`
	GenreName string `json:"genre_name" db:"genre_name"`
}

type IssueAct struct {
	Id         int           `json:"id"`
	UserId     int           `json:"user_id"`
	BookId     int           `json:"book_id"`
	RentalTime time.Duration `json:"rental_time"`
	ReturnDate time.Time     `json:"return_date"`
	PreCost    float32       `json:"pre_cost"`
	Status     bool          `json:"status"`
}

type ReturnAct struct {
	Id               int       `json:"id"`
	UserId           int       `json:"user_id"`
	BookId           int       `json:"book_id"`
	Cost             float32   `json:"cost"`
	ReturnDate       time.Time `json:"return_date"`
	Foto             []string  `json:"foto"`
	Fine             float32   `json:"fine"`
	ConditionDecrese int       `json:"condition_decrese"`
	Rating           int       `json:"rating"`
}
