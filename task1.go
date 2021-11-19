package main

import "fmt"

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
}
