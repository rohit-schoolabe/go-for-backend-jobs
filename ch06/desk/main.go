package main

import (
	"fmt"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch06/desk/loan"
)

func main() {
	bookCode := readLine("Book code: ")
	member := readLine("Member name: ")

	issued := time.Date(2026, 3, 2, 10, 0, 0, 0, time.Local)
	l, err := loan.New(bookCode, member, issued)
	if err != nil {
		fmt.Println("Cannot issue:", err)
		return
	}
	fmt.Println(l.Describe())

	l.Renew(7)
	fmt.Println(l.Describe())

	returned := issued.AddDate(0, 0, 30)
	fmt.Printf("Fine if returned on day 30: ₹%.2f\n",
		l.Fine(returned))
}

// readLine accepts an empty answer, unlike fmt.Scan,
// which would keep waiting for a word.
func readLine(prompt string) string {
	var value string
	fmt.Print(prompt)
	fmt.Scanln(&value)
	return value
}
