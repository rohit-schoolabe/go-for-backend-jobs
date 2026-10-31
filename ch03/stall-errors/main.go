package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	sales, err := readNumber("Total sales (₹)")
	if err != nil {
		fmt.Println(err)
		return
	}
	expenses, err := readNumber("Total expenses (₹)")
	if err != nil {
		fmt.Println(err)
		return
	}
	taxRate, err := readNumber("Tax rate (%)")
	if err != nil {
		fmt.Println(err)
		return
	}

	before, after, ratio :=
		calculateProfit(sales, expenses, taxRate)
	fmt.Printf("Profit before tax: %.2f\n", before)
	fmt.Printf("Profit after tax: %.2f\n", after)
	fmt.Printf("Ratio: %.3f\n", ratio)
	storeResults(before, after, ratio)
}

func readNumber(prompt string) (float64, error) {
	var value float64
	fmt.Print(prompt, ": ")
	if _, err := fmt.Scan(&value); err != nil {
		return 0, errors.New("not a number")
	}
	if value <= 0 {
		return 0, errors.New("value must be a positive number")
	}
	return value, nil
}

func calculateProfit(sales, expenses, taxRate float64) (
	float64, float64, float64) {
	before := sales - expenses
	after := before * (1 - taxRate/100)
	return before, after, before / after
}

func storeResults(before, after, ratio float64) {
	results := fmt.Sprintf(
		"Before: %.1f\nAfter: %.1f\nRatio: %.3f\n",
		before, after, ratio)
	err := os.WriteFile("results.txt", []byte(results), 0644)
	if err != nil {
		panic(err)
	}
}
