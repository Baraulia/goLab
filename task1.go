package main

import (
	"fmt"
	"strings"
)

//function checks if the string is a palindrome
func checkPalindrome(s string) {
	var reverse []byte
	for i := len(s) - 1; i >= 0; i-- {
		reverse = append(reverse, s[i])
	}
	switch s == string(reverse) {
	case true:
		fmt.Printf("Строка %s является палиндромом\n", s)
	case false:
		fmt.Printf("Строка %s не является палиндромом\n", s)
	}
}

//function sums the elements of the matrix
func matrixSum(m [][]int) int {
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
func reversChar(s string) string {
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

func main() {
	str1 := "aabaa"
	str2 := "abac"
	str3 := "a"
	checkPalindrome(str1)
	checkPalindrome(str2)
	checkPalindrome(str3)
	//________________________________________________
	matrix := [][]int{{0, 2, 3, 4},
		{0, 0, 0, 1},
		{1, 1, 4, 8},
	}
	matrix2 := [][]int{{1, 1, 1, 0},
		{0, 5, 0, 1},
		{2, 1, 3, 10}}
	fmt.Println(matrixSum(matrix))
	fmt.Println(matrixSum(matrix2))
	//________________________________________________
	s := "asdf(qwerty)jjkk(kldsf)jgj(fg)"
	fmt.Println(reversChar(s))

}
