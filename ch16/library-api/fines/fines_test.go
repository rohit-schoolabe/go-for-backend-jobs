package fines

import (
	"testing"
	"time"
)

func TestAmount(t *testing.T) {
	tests := []struct {
		name string
		days int
		want int64
	}{
		{"on time", 0, 0},
		{"returned early", -3, 0},
		{"one day", 1, 200},
		{"just under the cap", 24, 4800},
		{"at the cap", 25, 5000},
		{"far past the cap", 90, 5000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Standard.Amount(tt.days)
			assertPaise(t, got, tt.want)
		})
	}
}

// assertPaise reports the caller's line, not this one,
// because of t.Helper.
func assertPaise(t *testing.T, got, want int64) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, want %s",
			FormatRupees(got), FormatRupees(want))
	}
}

func TestDaysLate(t *testing.T) {
	// Due 15 March in India: 14 March 18:30 UTC.
	due := time.Date(2026, 3, 14, 18, 30, 0, 0, time.UTC)
	tests := []struct {
		name     string
		returned time.Time
		want     int
	}{
		{"11:30 pm IST on the due day",
			time.Date(2026, 3, 15, 18, 0, 0, 0, time.UTC), 0},
		{"12:30 am IST the next day",
			time.Date(2026, 3, 15, 19, 0, 0, 0, time.UTC), 1},
		{"a week later",
			time.Date(2026, 3, 22, 6, 0, 0, 0, time.UTC), 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DaysLate(due, tt.returned)
			if got != tt.want {
				t.Errorf("got %d days, want %d", got, tt.want)
			}
		})
	}
}

func TestParseRupees(t *testing.T) {
	good := map[string]int64{
		"12": 1200, "12.5": 1250, "Rs 12.50": 1250,
		"₹0.05": 5, " 7 ": 700,
	}
	for in, want := range good {
		got, err := ParseRupees(in)
		if err != nil {
			t.Errorf("ParseRupees(%q): %v", in, err)
			continue
		}
		assertPaise(t, got, want)
	}
	for _, in := range []string{
		"", "12.", "1.234", "-3", "1.+5", "x",
	} {
		if _, err := ParseRupees(in); err == nil {
			t.Errorf("ParseRupees(%q) should fail", in)
		}
	}
}
