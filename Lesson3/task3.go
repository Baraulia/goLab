package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

type Student struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Grade   int    `json:"grade"`
}
type Students struct {
	ListStudents []Student
}

const Filename = "mock/db.json"

var students = []Student{{1, "Bacилий", "Иванов", 11},
	{2, "Петр", "Петров", 10},
	{3, "Надежда", "Сидорова", 11},
}

func main() {
	r := httprouter.New()
	PORT := "0.0.0.0:8000"
	r.GET("/", index)
	r.GET("/students", allStudents)
	r.GET("/students?:grade", allStudents)
	r.GET("/students/:id", oneStudents)
	r.POST("/create", createStudent)
	fmt.Printf("Server listen on port:%s", PORT)
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("net.Listen", err)
	}
	server := &http.Server{
		Handler: r,
	}
	err = server.Serve(listener)
	if err != nil {
		log.Fatal("server", err)
	}
}

func CreateJsonFile(filename string, s Students) error {
	buf, err := json.Marshal(s.ListStudents)
	if err != nil {
		return err
	} else {
		err = ioutil.WriteFile(filename, buf, 0777)
		if err != nil {
			return err
		}
	}
	return err
}

func ReadJson(filename string) (Students, error) {
	var students Students
	f, err := os.Open(filename)
	if err != nil {
		return students, err
	} else {
		defer f.Close()
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return students, err
		} else {
			err = json.Unmarshal(data, &students.ListStudents)
			if err != nil {
				return students, err
			}
		}
	}
	return students, err
}

func createStudent(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var student Student
	r.ParseForm()
	var st = make(map[string]interface{})
	for key, value := range r.Form {
		if key == "grade" || key == "id" {
			v, _ := strconv.Atoi(value[0])
			st[key] = v
		} else {
			st[key] = value[0]
		}
	}
	js, err := json.Marshal(st)
	if err != nil {
		log.Fatal("Cannot decode Json:", err)
	}
	err = json.Unmarshal(js, &student)
	if err != nil {
		log.Fatal("Cannot decode Json:", err)
	}
	students, _ := ReadJson(Filename)
	students.ListStudents = append(students.ListStudents, student)
	CreateJsonFile(Filename, students)
	allStudents(w, r, params)
}

func allStudents(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	students, _ := ReadJson(Filename)
	tmpl, err := template.ParseFiles("templates/allstudents.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var GradeStudents Students
	if par, err := strconv.Atoi(r.URL.Query().Get("grade")); err == nil {
		for _, i := range students.ListStudents {
			if i.Grade == par {
				GradeStudents.ListStudents = append(GradeStudents.ListStudents, i)
			}
		}
		if err := tmpl.Execute(w, GradeStudents.ListStudents); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		if err := tmpl.Execute(w, students.ListStudents); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	students, err := ReadJson(Filename)
	if err != nil {
		log.Fatal("It is impossible to Read Json:", err)
	}
	var grades = make(map[int]int)
	for _, i := range students.ListStudents {
		grades[i.Grade] = i.Grade
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := tmpl.Execute(w, grades); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func oneStudents(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	students, _ := ReadJson(Filename)
	id, _ := strconv.Atoi(params.ByName("id"))
	tmpl, err := template.ParseFiles("templates/onestudents.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var student Student
	for _, i := range students.ListStudents {
		if i.Id == id {
			student = i
		}
	}
	if err := tmpl.Execute(w, student); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}
