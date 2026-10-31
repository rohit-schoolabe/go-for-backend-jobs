//go:build ignore

package main

func add(a, b any) any {
	return a + b
}

func main() {
	_ = add(2, 3)
}
