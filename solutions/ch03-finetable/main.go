package main

import "fmt"

func fineFor(daysLate int) float64 {
	if daysLate <= 0 {
		return 0
	}
	fine := float64(daysLate) * 2
	if fine > 50 {
		return 50
	}
	return fine
}

func main() {
	for day := 0; day <= 30; day++ {
		fmt.Printf("%2d days: Rs %5.2f\n",
			day, fineFor(day))
	}
}
