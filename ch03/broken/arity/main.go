//go:build ignore

package main

import "fmt"

func area(length, width float64) float64 {
	return length * width
}

func pair() (int, int) { return 1, 2 }

func main() {
	fmt.Println(area(5))       // not enough arguments
	fmt.Println(area(5, 3, 2)) // too many arguments
	x := pair()                // assignment mismatch
	fmt.Println(x)
}
