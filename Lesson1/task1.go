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
func MatrixSum(m [][]int) (s int, e error) {
	var sum int
	var err error
	var index = make(map[int]string)
	for i := range m {
		if err != nil {
			break
		}
		switch {
		case i == 0:
			if len(m[0]) == 0 {
				err = fmt.Errorf("Срез не должен быть пустым!")
				fmt.Println(err)
				break
			} else {
				for n, a := range m[0] {
					index[n] = ""
					sum += a
					if a == 0 {
						delete(index, n)
					}
				}
			}
		case i != 0:
			if len(m[0]) != len(m[i]) && len(m[i]) != 0 {
				err = fmt.Errorf("Передайте срез срезов с одинаковой длиной!")
				fmt.Println(err)
				break
			} else if len(m[i]) == 0 {
				err = fmt.Errorf("Срез не должен быть пустым!")
				fmt.Println(err)
				break
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
	return sum, err
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
			fmt.Println(par)
			switch val, ok := par[x]; val == 0 && ok {
			case false:
				break
			case true:
				par[x] = n
			}
		}
	}
	for start, stop := range par {
		if stop == 0 {
			continue
		}
		var reverse []rune
		for i := stop - 1; i > start; i-- {
			reverse = append(reverse, rune(s[i]))
		}
		s = strings.Replace(s, string(s[start+1:stop]), string(reverse), 1)
	}
	return s
}
func main() {

	str1 := "abac"
	fmt.Println(CheckPalindrome(str1))

	//________________________________________________
	matrix := [][]int{
		{0, 2, 3, 4},
		{0, 0, 0, 1},
		{1, 1, 4, 8},
	}
	fmt.Println(MatrixSum(matrix))

	//________________________________________________
	s := "qwe)rty(hjkkll(123)"
	fmt.Println(ReversChar(s))
}
