//go:build ignore

package main

import "fmt"

func total(prices []float64) (sum float64) {
	for _, p := range prices {
		sum := sum + p // a new sum, local to the loop
		_ = sum
	}
	return
}

func main() {
	x := 1
	x := 2 // no new variables on left side of :=
	fmt.Println(x, total([]float64{10, 20}))
}
