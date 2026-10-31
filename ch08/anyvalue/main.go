package main

import "fmt"

func describe(value any) {
	switch v := value.(type) {
	case int:
		fmt.Println("whole number:", v+1)
	case float64:
		fmt.Printf("decimal: %.2f\n", v)
	case string:
		fmt.Println("text of length", len(v))
	default:
		fmt.Printf("something else: %T\n", v)
	}
}

func main() {
	describe(41)
	describe(18.5)
	describe("Go for Backend")
	describe(true)

	var box any = "CS-101"
	code, ok := box.(string)
	fmt.Println(code, ok) // CS-101 true

	n, ok := box.(int)
	fmt.Println(n, ok) // 0 false
}
