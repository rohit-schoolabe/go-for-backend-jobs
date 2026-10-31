package main

import "fmt"

const renewDays = 7

// renewCopy returns a new due day; the caller's is untouched.
func renewCopy(dueDay int) int {
	dueDay += renewDays
	return dueDay
}

// extendDue moves the caller's due day forward in place.
func extendDue(dueDay *int, days int) {
	if dueDay == nil {
		return
	}
	*dueDay += days
}

func main() {
	dueDay := 120 // day of the year the book is due

	renewCopy(dueDay)
	fmt.Println("after renewCopy:", dueDay) // still 120

	dueDay = renewCopy(dueDay)
	fmt.Println("after assigning:", dueDay) // 127

	extendDue(&dueDay, renewDays)
	fmt.Println("after extendDue:", dueDay) // 134
}
