package main

import "fmt"

type Member struct {
	Name  string
	Phone string
}

func (m Member) Greeting() string {
	return "Hello, " + m.Name
}

// Staff embeds Member: its fields and methods are promoted.
type Staff struct {
	Member
	Department string
}

// Guest has a named field instead: nothing is promoted.
type Guest struct {
	host Member
	Note string
}

func main() {
	s := Staff{
		Member:     Member{Name: "Meena", Phone: "98450 00000"},
		Department: "Physics",
	}
	fmt.Println(s.Name)              // promoted field
	fmt.Println(s.Greeting())        // promoted method
	fmt.Println(s.Member.Greeting()) // the long way still works

	g := Guest{host: s.Member, Note: "visiting scholar"}
	fmt.Println(g.host.Name, "-", g.Note)
}
