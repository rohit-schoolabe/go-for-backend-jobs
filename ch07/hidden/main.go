package main

import (
	"encoding/json"
	"fmt"
)

type quiet struct {
	code  string
	title string
}

type loud struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Shelf  string `json:"shelf,omitempty"`
	Secret string `json:"-"`
}

func main() {
	q, _ := json.Marshal(quiet{"CS-101", "Go for Backend"})
	fmt.Println(string(q)) // {}

	l, _ := json.Marshal(loud{Code: "CS-101",
		Title: "Go for Backend", Secret: "staff only"})
	fmt.Println(string(l))
}
