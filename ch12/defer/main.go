package main

import "fmt"

func closeDesk() {
	defer fmt.Println("lights off")
	defer fmt.Println("lock the till")
	fmt.Println("count the cash")
}

func main() {
	closeDesk()

	copies := 3
	defer fmt.Println("copies at defer time:", copies)
	copies = 5
	fmt.Println("copies now:", copies)
}
