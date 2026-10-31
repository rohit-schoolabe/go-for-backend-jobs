package main

import "fmt"

func factorial(n int) int {
	if n <= 1 {
		return 1 // base case
	}
	return n * factorial(n-1)
}

type section struct {
	name     string
	books    int
	children []section
}

func countBooks(s section) int {
	total := s.books
	for _, child := range s.children {
		total += countBooks(child)
	}
	return total
}

func main() {
	fmt.Println(factorial(5)) // 120

	library := section{"Library", 0, []section{
		{"Computing", 40, []section{
			{"Go", 12, nil},
			{"Databases", 9, nil},
		}},
		{"Fiction", 75, nil},
	}}
	fmt.Println(countBooks(library)) // 136
}
