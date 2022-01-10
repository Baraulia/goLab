package handler

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const tagName = "validate"

var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

type Validator interface {
	Validate(interface{}) error
}

type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}) error {
	return nil
}

type NumberValidator struct {
	Min int
	Max int
}

func (v NumberValidator) Validate(val interface{}) error {
	num := val.(int)
	if num < v.Min {
		return fmt.Errorf("should be greater than %v", v.Min)
	}
	if v.Max >= v.Min && num > v.Max {
		return fmt.Errorf("should be less than %v", v.Max)
	}
	return nil
}

type StringValidator struct {
	Min int
	Max int
}

func (v StringValidator) Validate(val interface{}) error {
	l := len(val.(string))

	if l == 0 {
		return fmt.Errorf("the value of the field cannot be empty")
	}

	if l < v.Min {
		return fmt.Errorf("should be at least %v chars long", v.Min)
	}

	if v.Max >= v.Min && l > v.Max {
		return fmt.Errorf("should be less than %v chars long", v.Max)
	}

	return nil
}

type EmailValidator struct {
}

func (v EmailValidator) Validate(val interface{}) error {
	if !mailRe.MatchString(val.(string)) {
		return fmt.Errorf("it is not a valid email address")
	}
	return nil
}

type GenreExistValidator struct {
	handler *Handler
}

func (v GenreExistValidator) Validate(val interface{}) error {
	for _, genre := range val.([]int) {
		err := v.handler.services.Validation.GetGenreById(genre)
		if err != nil {
			return fmt.Errorf("such genre with id = %d does not exist:%s", genre, err)
		}
	}
	return nil
}

type AuthorExistValidator struct {
	handler *Handler
}

func (v AuthorExistValidator) Validate(val interface{}) error {
	for _, author := range val.([]int) {
		err := v.handler.services.Validation.GetAuthorById(author)
		if err != nil {
			return fmt.Errorf("such author with id = %d does not exist:%s", author, err)
		}
	}
	return nil
}

type BirthDayValidator struct {
}

func (v BirthDayValidator) Validate(val interface{}) error {
	if val.(IndTask.MyTime).Time.After(time.Now()) {
		return fmt.Errorf("the birthday is incorrect")
	} else {
		age := time.Now().Sub(val.(IndTask.MyTime).Time)
		if age.Hours() < 24*365*10 {
			return fmt.Errorf("user is too young: he is %d", uint(age.Hours()/(24*365)))
		} else if age.Hours() > 24*365*100 {
			return fmt.Errorf("people don't live that much:%d years", uint(age.Hours()/(24*365)))
		}
	}
	return nil
}

type UserExistValidator struct {
	handler *Handler
}

func (v UserExistValidator) Validate(val interface{}) error {
	err := v.handler.services.Validation.GetUserById(val.(int))
	if err != nil {
		return fmt.Errorf("such user with id = %d does not exist:%s", val.(int), err)
	}

	return nil
}

type ListBookExistValidator struct {
	handler *Handler
}

func (v ListBookExistValidator) Validate(val interface{}) error {
	err := v.handler.services.Validation.GetListBookById(val.(int))
	if err != nil {
		return fmt.Errorf("such listbook with id = %d does not exist:%s", val.(int), err)
	}

	return nil
}

type RentalTimeValidator struct {
}

func (v RentalTimeValidator) Validate(val interface{}) error {
	if val.(IndTask.MyDuration).Duration > 30*24*time.Hour {
		return fmt.Errorf("cannot lend a book for such a long period:%f", val.(IndTask.MyDuration).Duration.Hours()/24)
	}
	return nil
}

type IssueActExistValidator struct {
	handler *Handler
}

func (v IssueActExistValidator) Validate(val interface{}) error {
	err := v.handler.services.Validation.GetIssueActById(val.(int))
	if err != nil {
		return fmt.Errorf("such issue act with id = %d does not exist:%s", val.(int), err)
	}

	return nil
}

func getValidatorFromTag(tag string, h *Handler) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number":
		validator := NumberValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "email":
		return EmailValidator{}
	case "genreExist":
		return GenreExistValidator{handler: h}
	case "authorExist":
		return AuthorExistValidator{handler: h}
	case "birthDay":
		return BirthDayValidator{}
	case "userExist":
		return UserExistValidator{handler: h}
	case "listBookExist":
		return ListBookExistValidator{handler: h}
	case "rentalTime":
		return RentalTimeValidator{}
	case "issueActExist":
		return IssueActExistValidator{handler: h}
	}
	return DefaultValidator{}
}

func validateStruct(h *Handler, s interface{}) []error {
	var errs []error
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		validator := getValidatorFromTag(tag, h)
		err := validator.Validate(v.Field(i).Interface())
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err))
		}
	}
	return errs
}
