package fines

import (
	"testing"
	"time"
)

// The b.N loop works on every Go version.
func BenchmarkFormatRupees(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FormatRupees(int64(i))
	}
}

func BenchmarkDaysLate(b *testing.B) {
	due := time.Now()
	returns := make([]time.Time, 1000)
	for i := range returns {
		returns[i] = due.Add(time.Duration(i) * time.Hour)
	}
	b.ResetTimer() // don't count the setup above
	for i := 0; i < b.N; i++ {
		DaysLate(due, returns[i%len(returns)])
	}
}
