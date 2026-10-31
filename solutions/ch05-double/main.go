package main

import "fmt"

func double(n *int) {
	*n *= 2
}

func main() {
	copies := 4
	double(&copies)
	fmt.Println(copies) // 8
}
