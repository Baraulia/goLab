package main

import "fmt"

func main() {
	// TASK 1
	str1 := "aabaa"
	str2 := "abac"
	str3 := "a"
	fmt.Println(CheckPalindrome(str1))
	fmt.Println(CheckPalindrome(str2))
	fmt.Println(CheckPalindrome(str3))
	//________________________________________________
	matrix := [][]int{{0, 2, 3, 4},
		{0, 0, 0, 1},
		{1, 1, 4, 8},
	}
	matrix2 := [][]int{{1, 1, 1, 0},
		{0, 5, 0, 1},
		{2, 1, 3, 10}}
	fmt.Println(MatrixSum(matrix))
	fmt.Println(MatrixSum(matrix2))
	//________________________________________________
	s := "asdf(qwerty)jjkk(kldsf)jgj(fg)"
	fmt.Println(ReversChar(s))

	//TASK 2

}
