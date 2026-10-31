package main

import "fmt"

func showMenu() {
	fmt.Println()
	fmt.Println("1. Check copies")
	fmt.Println("2. Return copies")
	fmt.Println("3. Issue copies")
	fmt.Println("4. Exit")
}

func readCount(prompt string) (int, error) {
	var n int
	fmt.Print(prompt)
	if _, err := fmt.Scan(&n); err != nil {
		return 0, fmt.Errorf("please type a whole number")
	}
	if n <= 0 {
		return 0, fmt.Errorf("the count must be positive")
	}
	return n, nil
}
