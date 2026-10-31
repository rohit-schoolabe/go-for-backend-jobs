//go:build ignore

package main

func addGrace(m *int) {
	*m = *m + 5
}

func main() {
	marks := 72
	addGrace(marks)
}
