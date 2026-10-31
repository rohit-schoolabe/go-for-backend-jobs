package fines

import "testing"

// FuzzParseRupees checks a property rather than examples:
// whatever parses must survive a round trip.
func FuzzParseRupees(f *testing.F) {
	for _, seed := range []string{"12", "12.5", "Rs 3.05",
		"₹0", "1.", "-1"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		p, err := ParseRupees(s)
		if err != nil {
			return
		}
		if p < 0 {
			t.Fatalf("%q parsed to negative %d", s, p)
		}
		back, err := ParseRupees(FormatRupees(p))
		if err != nil || back != p {
			t.Fatalf("%q -> %d -> %q -> %d, %v",
				s, p, FormatRupees(p), back, err)
		}
	})
}
