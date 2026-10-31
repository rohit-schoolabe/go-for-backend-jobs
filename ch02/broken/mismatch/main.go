//go:build ignore

package main

import "fmt"

func main() {
	var depositAmount = 50000
	var interestRate = 7.0
	var maturity = depositAmount * (1 + interestRate/100)
	fmt.Println(maturity)
}
