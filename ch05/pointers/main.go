package main

import "fmt"

func main() {
	marks := 72
	fmt.Println(marks)  // 72
	fmt.Println(&marks) // an address such as 0xc000012078

	var marksPointer *int = &marks
	samePointer := &marks
	fmt.Println(marksPointer == samePointer) // true
	fmt.Println(*marksPointer)               // 72

	final := withGrace(marks)
	fmt.Println(marks, final) // 72 77

	addGrace(&marks)
	fmt.Println(marks) // 77

	var nothing *int
	fmt.Println(nothing == nil) // true
}

func withGrace(m int) int {
	return m + 5
}

func addGrace(m *int) {
	*m = *m + 5
}
