package main

import (
	"fmt"
	"math"
)

// inflationRate is an assumption for the whole program.
const inflationRate = 6.0

func main() {
	var depositAmount, interestRate, years float64

	fmt.Print("Deposit amount (₹): ")
	fmt.Scan(&depositAmount)
	fmt.Print("Yearly interest rate (%): ")
	fmt.Scan(&interestRate)
	fmt.Print("Number of years: ")
	fmt.Scan(&years)

	growth := math.Pow(1+interestRate/100, years)
	maturity := depositAmount * growth
	realValue := maturity / math.Pow(1+inflationRate/100, years)

	fmt.Printf("Maturity amount: ₹%.2f\n", maturity)
	fmt.Printf("In today's money (%.1f%% inflation): ₹%.2f\n",
		inflationRate, realValue)
}
