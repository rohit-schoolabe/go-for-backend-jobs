package main

import "fmt"

func main() {
	var sales, expenses, taxRate float64
	fmt.Print("Sales: ")
	fmt.Scan(&sales)
	fmt.Print("Expenses: ")
	fmt.Scan(&expenses)
	fmt.Print("Tax rate (%): ")
	fmt.Scan(&taxRate)

	before := sales - expenses
	after := before * (1 - taxRate/100)
	ratio := before / after

	fmt.Printf("Before tax: %.2f\n", before)
	fmt.Printf("After tax:  %.2f\n", after)
	fmt.Printf("Ratio:      %.2f\n", ratio)
}
