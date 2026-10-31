package main

import "fmt"

type book struct {
	code  string
	title string
}

func main() {
	waitlist := []string{"Asha", "Rahul"}
	waitlist = append(waitlist, "Meera")
	fmt.Println(waitlist, len(waitlist)) // [Asha Rahul Meera] 3

	// Asha got her book: drop the first name.
	waitlist = waitlist[1:]
	fmt.Println(waitlist) // [Rahul Meera]

	// Remove index 0 of the remaining list, keeping order.
	i := 0
	waitlist = append(waitlist[:i], waitlist[i+1:]...)
	fmt.Println(waitlist) // [Meera]

	morning := []string{"Kiran", "Dev"}
	waitlist = append(waitlist, morning...)
	fmt.Println(waitlist) // [Meera Kiran Dev]

	shelf := []book{
		{"CS-101", "Go for Backend"},
		{"DB-210", "SQL Basics"},
	}
	shelf = append(shelf, book{"NW-330", "Networks"})
	fmt.Println(len(shelf), shelf[2].title) // 3 Networks

	codes := make([]string, 0, len(shelf))
	for _, b := range shelf {
		codes = append(codes, b.code)
	}
	fmt.Println(codes) // [CS-101 DB-210 NW-330]

	slots := make([]string, 2)
	slots = append(slots, "late")
	fmt.Printf("%q\n", slots) // ["" "" "late"]

	for i, name := range waitlist {
		fmt.Println(i+1, name)
	}
}
