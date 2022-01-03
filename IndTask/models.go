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
	BookName  string  `json:"book_name"`
	GenreId   []int   `json:"genre_id"`
	Cost      float32 `json:"cost"`
	AuthorsId []int   `json:"authors"`
	Cover     string  `json:"cover"`
	Published int     `json:"published"`
	Pages     int     `json:"pages"`
	Amount    int     `json:"amount"`
}

type ListBooks struct {
	Id         int       `json:"id"`
	BookId     int       `json:"book_id"`
	Issued     bool      `json:"issued"`
	RentNumber int       `json:"rent_number"`
	RentCost   float64   `json:"rent_cost"`
	RegDate    time.Time `json:"reg_date"`
	Condition  int       `json:"condition"`
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
	Id         int        `json:"id"`
	UserId     int        `json:"user_id"`
	ListBookId int        `json:"list_book_id"`
	RentalTime MyDuration `json:"rental_time"`
	ReturnDate time.Time  `json:"return_date"`
	PreCost    float64    `json:"pre_cost"`
	Cost       float64    `json:"cost"`
	Status     string     `json:"status"`
}

type ReturnAct struct {
	Id               int       `json:"id"`
	IssueActId       int       `json:"issue_act_id"`
	ReturnDate       time.Time `json:"return_date"`
	Foto             []string  `json:"foto"`
	Fine             float64   `json:"fine"`
	ConditionDecrese int       `json:"condition_decrese"`
	Rating           int       `json:"rating"`
}
