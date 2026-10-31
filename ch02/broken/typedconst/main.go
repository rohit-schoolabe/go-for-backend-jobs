//go:build ignore

package main

import "fmt"

func main() {
	const big = 1 << 40   // untyped
	var f float32 = big   // fine: converted at compile time
	const limit int = 100 // typed
	var g float64 = limit // error: limit has type int
	fmt.Println(f, g)
}
