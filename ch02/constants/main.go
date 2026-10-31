package main

import "fmt"

const secondsPerDay = 24 * 60 * 60

func main() {
	const big = 1 << 40 // untyped: adapts to where it is used
	var f float32 = big
	var balance = 1000.0
	balance = balance + 250 // now 1250
	balance = 0             // now 0
	var count int           // zero value: 0
	var name string         // zero value: ""
	var open bool           // zero value: false
	fmt.Println(secondsPerDay, f, balance, count, name, open)
}
