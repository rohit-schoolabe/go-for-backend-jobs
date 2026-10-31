package main

import "fmt"

func add[T int | float64 | string](a, b T) T {
	return a + b
}

// Number lists the types sum accepts.
type Number interface {
	int | int64 | float64
}

func sum[T Number](values []T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

func main() {
	fmt.Println(add(2, 3))           // 5
	fmt.Println(add(2.5, 0.25))      // 2.75
	fmt.Println(add("Go ", "rocks")) // Go rocks
	fmt.Println(add[float64](2, 3))  // 5, T chosen explicitly

	fines := []float64{18, 4, 50}
	fmt.Println(sum(fines)) // 72
}
