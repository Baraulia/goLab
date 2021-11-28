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

func main() {
	r := httprouter.New()
	port := "0.0.0.0:8000"
	r.GET("/", index)
	r.GET("/students", allStudents)
	r.GET("/students?:grade", allStudents)
	r.GET("/students/:id", oneStudents)
	fmt.Printf("Server listen on port:%s", port)
	listener, err := net.Listen("tcp", port)
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

func CreateJsonFile() {
	var students Students
	students.ListStudents = []Student{{1, "Bacилий", "Иванов", 11},
		{2, "Петр", "Петров", 10},
		{3, "Надежда", "Сидорова", 11},
	}
	buf, err := json.Marshal(students)
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}
	err = ioutil.WriteFile(Filename, buf, 0777)
	if err != nil {
		log.Fatal("Cannot write updated file:", err)
	}

}

func allStudents(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var students Students
	data, err := ioutil.ReadFile(Filename)
	if err != nil {
		log.Fatal("Cannot load settings:", err)
	}
	err = json.Unmarshal(data, &students)
	if err != nil {
		log.Fatal("Invalid data format:", err)
	}
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
		fmt.Println(err)
		if err := tmpl.Execute(w, students.ListStudents); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var students Students
	data, err := ioutil.ReadFile(Filename)
	if err != nil {
		log.Fatal("Cannot load settings:", err)
	}
	err = json.Unmarshal(data, &students)
	if err != nil {
		log.Fatal("Invalid data format:", err)
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
	var students Students
	data, err := ioutil.ReadFile(Filename)
	if err != nil {
		log.Fatal("Cannot load settings:", err)
	}
	err = json.Unmarshal(data, &students)
	if err != nil {
		log.Fatal("Invalid data format:", err)
	}
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
