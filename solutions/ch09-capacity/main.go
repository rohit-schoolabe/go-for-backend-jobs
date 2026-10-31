package main

import "fmt"

func main() {
	s := make([]int, 3, 10)
	t := s[:5]
	fmt.Println(s, len(s), cap(s)) // [0 0 0] 3 10
	fmt.Println(t, len(t), cap(t)) // [0 0 0 0 0] 5 10

	t[4] = 9
	s = append(s, 1)
	fmt.Println(t) // [0 0 0 1 9]
}
