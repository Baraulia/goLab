package main

import (
	"fmt"
	"strings"
)

//function checks if the string is a palindrome
func CheckPalindrome(s string) string {
	var reverse []byte
	for i := len(s) - 1; i >= 0; i-- {
		reverse = append(reverse, s[i])
	}
	var response string
	switch s == string(reverse) {
	case true:
		response = fmt.Sprintf("Строка %s является палиндромом\n", s)
	case false:
		response = fmt.Sprintf("Строка %s не является палиндромом\n", s)
	}
	return response
}

//function sums the elements of the matrix
func MatrixSum(m [][]int) int {
	var sum int
	var index = make(map[int]string)
	for i := range m {
		switch {
		case i == 0:
			for n, a := range m[0] {
				index[n] = ""
				sum += a
				if a == 0 {
					delete(index, n)
				}
			}
		case i != 0:
			if len(m[0]) != len(m[i]) {
				fmt.Println("Передайте срез срезов с одинаковой длиной!")
			} else {
				for n := range index {
					if m[i][n] == 0 {
						delete(index, n)
					} else {
						sum += m[i][n]
					}

				}
			}
		}
	}
	return sum
}

//function reverses characters in parentheses
func ReversChar(s string) string {
	par := make(map[int]int)
	var x int
	for n, i := range s {
		if string(i) == "(" {
			x = n
			par[x] = 0
		} else if string(i) == ")" {
			par[x] = n
		}
	}
	for start, stop := range par {
		var reverse []rune
		for i := stop - 1; i > start; i-- {
			reverse = append(reverse, rune(s[i]))
		}
		s = strings.Replace(s, string(s[start+1:stop]), string(reverse), 1)
	}
	return s
}
