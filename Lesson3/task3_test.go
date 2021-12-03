package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	//"net/http/httptest"
	"testing"
)

func TestCreateJsonFile(t *testing.T) {
	Convey("Create Json file", t, func() {
		var st = Students{}
		filename := "mock/test_db.json"
		Convey("Empty ListStudents", func() {
			err := CreateJsonFile(filename, st)
			So(err, ShouldEqual, nil)
		})
		Convey("Denied permission", func() {
			filename := "mock/test2_db.json"
			ioutil.WriteFile(filename, []byte(""), 0444)
			err := CreateJsonFile(filename, st)
			So(err, ShouldBeError, "open mock/test2_db.json: permission denied")
			os.Remove(filename)
		})
		os.Remove(filename)

	})
}

func TestReadJson(t *testing.T) {
	Convey("Read Json file", t, func() {
		filename := "mock/test_db.json"
		Convey("Read non-existent Json file", func() {
			_, err := ReadJson(filename)
			So(err, ShouldBeError, "open mock/test_db.json: no such file or directory")
		})
		Convey("Read existent Json file", func() {
			var st Students
			st.ListStudents = []Student{{1, "Sergey", "Tumakov", 12}}
			CreateJsonFile(filename, st)
			_, err := ReadJson(filename)
			So(err, ShouldEqual, nil)
		})

		os.Remove(filename)

	})
}
