package main

import "fmt"

func main() {
	var sales, expenses, taxRate float64

	fmt.Print("Total sales (₹): ")
	fmt.Scan(&sales)
	fmt.Print("Total expenses (₹): ")
	fmt.Scan(&expenses)
	fmt.Print("Tax rate (%): ")
	fmt.Scan(&taxRate)

	profitBeforeTax := sales - expenses
	profitAfterTax := profitBeforeTax * (1 - taxRate/100)
	ratio := profitBeforeTax / profitAfterTax

	fmt.Println("Profit before tax:", profitBeforeTax)
	fmt.Println("Profit after tax:", profitAfterTax)
	fmt.Println("Ratio:", ratio)
}
