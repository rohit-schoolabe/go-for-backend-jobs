//go:build ignore

package main

func withGrace(m *int) int {
	return m + 5
}

func main() {
	marks := 72
	_ = withGrace(&marks)
}
