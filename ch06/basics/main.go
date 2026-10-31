package main

import (
	"fmt"
	"time"
)

type book struct {
	code    string
	title   string
	author  string
	copies  int
	addedAt time.Time
}

func main() {
	gopl := book{
		code:    "CS-101",
		title:   "Go for Backend",
		author:  "R. Srivastava",
		copies:  3,
		addedAt: time.Now(),
	}

	dbms := book{"CS-205", "Databases", "K. Rao", 2, time.Now()}

	var empty book
	partial := book{title: "Untitled"}

	printBook(gopl)
	printBook(dbms)
	fmt.Printf("%q %d\n", empty.title, empty.copies) // "" 0
	fmt.Printf("%q %d\n", partial.title, partial.copies)

	addCopy(&gopl)
	fmt.Println(gopl.copies) // 4
}

func printBook(b book) {
	fmt.Println(b.code, b.title, "by", b.author)
}

func addCopy(b *book) {
	b.copies++ // short for (*b).copies++
}
