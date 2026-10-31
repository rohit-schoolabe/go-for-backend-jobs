package main

import "fmt"

func maxOf[T int | float64 | string](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println(maxOf(3, 7))             // 7
	fmt.Println(maxOf(2.5, 1.5))         // 2.5
	fmt.Println(maxOf("Zebra", "apple")) // apple
	fmt.Println(maxOf("Go 10", "Go 9"))  // Go 9
}
