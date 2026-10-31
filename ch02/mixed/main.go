package main

import (
	"fmt"
	"math"
)

func main() {
	var depositAmount = 50000
	var interestRate = 7.0
	var years = 5

	var maturity = float64(depositAmount) *
		math.Pow(1+interestRate/100, float64(years))
	fmt.Println(maturity)

	var marks = 78 + 85 + 91
	var average = marks / 3 // 84: integer division
	fmt.Println(marks, average, 7/2, 2^3)
}
