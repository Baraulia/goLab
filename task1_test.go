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
		Reset(func() {
			t.Log("FINISH")
		})

	})
}
