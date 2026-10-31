package main

import (
	"fmt"
	"slices"
	"strings"
)

type book struct {
	Code   string
	Title  string
	Copies int
}

type catalog struct {
	books  []book
	byCode map[string]int // code -> index in books
}

func newCatalog() *catalog {
	return &catalog{byCode: make(map[string]int)}
}

func (c *catalog) add(b book) error {
	if _, exists := c.byCode[b.Code]; exists {
		return fmt.Errorf("duplicate code %s", b.Code)
	}
	c.books = append(c.books, b)
	c.byCode[b.Code] = len(c.books) - 1
	return nil
}

func (c *catalog) find(code string) (book, bool) {
	i, ok := c.byCode[strings.ToUpper(code)]
	if !ok {
		return book{}, false
	}
	return c.books[i], true
}

func (c *catalog) titles() []string {
	out := make([]string, 0, len(c.books))
	for _, b := range c.books {
		out = append(out, b.Title)
	}
	slices.Sort(out)
	return out
}

func main() {
	c := newCatalog()
	c.add(book{"CS-101", "Go for Backend", 3})
	c.add(book{"DB-210", "SQL Basics", 1})
	if err := c.add(book{"CS-101", "Copy", 1}); err != nil {
		fmt.Println(err) // duplicate code CS-101
	}

	if b, ok := c.find("db-210"); ok {
		fmt.Println(b.Title, b.Copies) // SQL Basics 1
	}
	_, ok := c.find("ZZ-000")
	fmt.Println(ok)         // false
	fmt.Println(c.titles()) // [Go for Backend SQL Basics]
}
