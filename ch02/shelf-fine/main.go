package main

import "fmt"

const (
	finePerDay = 2.0  // ₹ per day late
	maxFine    = 50.0 // a fine never exceeds this
)

func main() {
	var title string
	var daysLate int

	fmt.Print("Book code: ")
	fmt.Scan(&title)
	fmt.Print("Days late: ")
	fmt.Scan(&daysLate)

	fine := float64(daysLate) * finePerDay
	fine = min(fine, maxFine)

	receipt := fmt.Sprintf("%-10s %3d days  ₹%6.2f",
		title, daysLate, fine)
	fmt.Println(receipt)
	fmt.Printf("Type of fine: %T\n", fine)
}
