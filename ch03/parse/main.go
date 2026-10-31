package main

import (
	"fmt"
	"strconv"
)

func main() {
	fine, err := strconv.ParseFloat("18.50", 64)
	fmt.Println(fine, err) // 18.5 <nil>

	_, err = strconv.ParseFloat("eighteen", 64)
	fmt.Println(err)
	// strconv.ParseFloat: parsing "eighteen": invalid syntax

	average := 254 / 3.0 // the .0 makes this a float division
	fmt.Printf("%.2f\n", average)
}
