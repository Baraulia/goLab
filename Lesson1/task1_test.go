package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCheckPalindrome(t *testing.T) {
	Convey("Check if the string is a palindrome", t, func() {
		Convey("Check string with len=1", func() {
			So(CheckPalindrome("a"), ShouldEqual, "Строка a является палиндромом\n")
		})
		Convey("Check empty string", func() {
			So(CheckPalindrome(""), ShouldEqual, "Строка  является палиндромом\n")

		})
		Convey("Check not palendrome", func() {
			So(CheckPalindrome("qwerty"), ShouldEqual, "Строка qwerty не является палиндромом\n")

		})
		Reset(func() {
			t.Log("FINISH")
		})

	})
}
func TestMatrixSum(t *testing.T) {
	Convey("Different slices", t, func() {
		_, err := MatrixSum([][]int{{0, 1, 2, 3}, {5, 0, 0, 8}, {5, 6, 7, 8, 10}})
		So(err, ShouldBeError, ("Передайте срез срезов с одинаковой длиной!"))
	})
	Convey("Empty slice", t, func() {
		_, err := MatrixSum([][]int{{}, {0, 1, 2, 3}, {}, {5, 6, 7, 8, 10}})
		So(err, ShouldBeError, ("Срез не должен быть пустым!"))
		_, err = MatrixSum([][]int{{0, 1, 2, 3}, {}, {5, 6, 7, 8, 10}})
		So(err, ShouldBeError, ("Срез не должен быть пустым!"))
	})

}

func TestReversChar(t *testing.T) {
	Convey("First custom case with brackets", t, func() {
		So(ReversChar(")(asdf(123)"), ShouldEqual, ")(asdf(321)")
	})
	Convey("Second custom case with brackets", t, func() {
		So(ReversChar(")()asdf)"), ShouldEqual, ")()asdf)")
	})
}
