// Package fines holds the Chapter 16 exercise answers.
package fines

import "fmt"

// formatRupees handles negative amounts, such as refunds:
// -150 paise is "Rs -1.50", not "Rs -1.-50".
func formatRupees(paise int64) string {
	sign := ""
	if paise < 0 {
		sign = "-"
		paise = -paise
	}
	return fmt.Sprintf("Rs %s%d.%02d", sign,
		paise/100, paise%100)
}
