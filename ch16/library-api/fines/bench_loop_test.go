//go:build go1.24

package fines

import "testing"

// b.Loop needs Go 1.24. It times only the loop body, so
// setup before it needs no ResetTimer.
func BenchmarkParseRupeesLoop(b *testing.B) {
	inputs := []string{"12", "Rs 12.50", "₹0.05"}
	b.ReportAllocs()
	for b.Loop() {
		for _, s := range inputs {
			ParseRupees(s)
		}
	}
}
