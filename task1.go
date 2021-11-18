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

func main() {
	str1 := "aabaa"
	str2 := "abac"
	str3 := "a"
	checkPalindrome(str1)
	checkPalindrome(str2)
	checkPalindrome(str3)
}
