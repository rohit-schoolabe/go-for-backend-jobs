// Package fines computes late fines in paise, so that
// money is never a float.
package fines

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

// Rule is a fine rate and a cap, both in paise.
type Rule struct {
	PerDay int64
	Cap    int64
}

// Standard is the library's rule: Rs 2 a day, at most Rs 50.
var Standard = Rule{PerDay: 200, Cap: 5000}

var ist = mustLoad("Asia/Kolkata")

func mustLoad(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

// DaysLate counts calendar days in India between the due
// date and the return. Returning on the due day is free.
func DaysLate(due, returned time.Time) int {
	d := day(due)
	r := day(returned)
	if !r.After(d) {
		return 0
	}
	return int(r.Sub(d).Hours() / 24)
}

// day is midnight UTC on the Indian calendar date of t.
func day(t time.Time) time.Time {
	y, m, dd := t.In(ist).Date()
	return time.Date(y, m, dd, 0, 0, 0, 0, time.UTC)
}

// Amount is the fine for a number of days late.
func (r Rule) Amount(days int) int64 {
	if days <= 0 {
		return 0
	}
	return min(int64(days)*r.PerDay, r.Cap)
}

// FormatRupees turns 1250 paise into "Rs 12.50".
func FormatRupees(paise int64) string {
	return fmt.Sprintf("Rs %d.%02d", paise/100, paise%100)
}

var errAmount = errors.New("not an amount in rupees")

// ParseRupees reads what a librarian types into a form,
// such as "12", "12.5", "Rs 12.50" or "₹12.50".
func ParseRupees(s string) (int64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "₹")
	s = strings.TrimPrefix(s, "Rs")
	s = strings.TrimSpace(s)
	whole, frac, hasDot := strings.Cut(s, ".")
	if whole == "" || len(frac) > 2 || (hasDot && frac == "") ||
		strings.ContainsAny(s, "+-_") {
		return 0, errAmount
	}
	for len(frac) < 2 {
		frac += "0"
	}
	r, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || r < 0 || r > 1_000_000 {
		return 0, errAmount
	}
	p, err := strconv.ParseInt(frac, 10, 64)
	if err != nil || p < 0 {
		return 0, errAmount
	}
	return r*100 + p, nil
}
