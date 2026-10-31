// Command timezones shows why Shelf stores UTC
// and shows Indian time only at the edges.
package main

import (
	"fmt"
	"time"
	_ "time/tzdata" // zone data inside the binary
)

func main() {
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println(err)
		return
	}

	// The desk types a due date. It means midnight in India.
	wrong, _ := time.Parse("2006-01-02", "2026-03-15")
	due, _ := time.ParseInLocation("2006-01-02",
		"2026-03-15", ist)
	fmt.Println("Parse:          ", wrong)
	fmt.Println("ParseInLocation:", due)
	fmt.Println("difference:     ", wrong.Sub(due))

	// Store and compare in UTC; show IST to people.
	stored := due.UTC()
	fmt.Println("stored:", stored.Format(time.RFC3339))
	fmt.Println("shown: ",
		stored.In(ist).Format("02 Jan 2006, 3:04 PM MST"))

	// The same instant compares equal in any zone.
	fmt.Println("equal:", stored.Equal(due))
}
