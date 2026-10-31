package main

import (
	"fmt"
	"math"
)

const inflationRate = 6.0 // package scope: calculate can see it

func main() {
	printPrompt("Deposit amount (₹)")
	fmt.Println()

	maturity, realValue := calculate(50000, 7, 5)
	fmt.Printf("%.2f %.2f\n", maturity, realValue)

	onlyMaturity, _ := calculate(50000, 7, 5)
	fmt.Printf("%.2f\n", onlyMaturity)

	m, r := calculateNamed(50000, 7, 5)
	fmt.Printf("%.2f %.2f\n", m, r)
}

func printPrompt(text string) {
	fmt.Print(text, ": ")
}

func calculate(amount, rate, years float64) (float64, float64) {
	maturity := amount * math.Pow(1+rate/100, years)
	realValue := maturity / math.Pow(1+inflationRate/100, years)
	return maturity, realValue
}

func calculateNamed(amount, rate, years float64) (
	maturity float64, realValue float64) {
	maturity = amount * math.Pow(1+rate/100, years)
	realValue = maturity / math.Pow(1+inflationRate/100, years)
	return
}
