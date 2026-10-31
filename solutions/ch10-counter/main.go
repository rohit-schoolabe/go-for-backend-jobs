package main

import "fmt"

func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func main() {
	a := counter()
	b := counter()
	fmt.Println(a(), a(), a()) // 1 2 3
	fmt.Println(b())           // 1
}
