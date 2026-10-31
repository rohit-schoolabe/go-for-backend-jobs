package fines

import (
	"testing"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/fines"
)

func TestFormatRupees(t *testing.T) {
	tests := []struct {
		paise int64
		want  string
	}{
		{0, "Rs 0.00"}, {5, "Rs 0.05"}, {100, "Rs 1.00"},
		{1250, "Rs 12.50"}, {123456, "Rs 1234.56"},
		{-150, "Rs -1.50"},
	}
	for _, tt := range tests {
		if got := formatRupees(tt.paise); got != tt.want {
			t.Errorf("%d: got %q, want %q",
				tt.paise, got, tt.want)
		}
	}
}

func FuzzDaysLate(f *testing.F) {
	f.Add(int64(1773513000), int64(1773601200))
	f.Add(int64(0), int64(-86400))
	f.Fuzz(func(t *testing.T, due, back int64) {
		d := time.Unix(due, 0)
		r := time.Unix(back, 0)
		if n := fines.DaysLate(d, r); n < 0 {
			t.Fatalf("negative: %d", n)
		}
		if n := fines.DaysLate(d, d); n != 0 {
			t.Fatalf("due instant gave %d", n)
		}
	})
}
