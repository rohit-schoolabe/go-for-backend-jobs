package main

import "fmt"

type ratingMap map[string]float64

func (r ratingMap) print() {
	for title, stars := range r {
		fmt.Printf("%-16s %.1f\n", title, stars)
	}
}

func main() {
	shelfOf := map[string]string{
		"CS-101": "A3",
		"DB-210": "B1",
	}
	fmt.Println(shelfOf["CS-101"]) // A3

	shelfOf["NW-330"] = "C2" // add
	shelfOf["DB-210"] = "B4" // overwrite
	delete(shelfOf, "CS-101")
	fmt.Println(len(shelfOf)) // 2

	fmt.Printf("%q\n", shelfOf["XX-999"]) // ""
	shelf, ok := shelfOf["XX-999"]
	fmt.Printf("%q %v\n", shelf, ok) // "" false

	ratings := make(ratingMap, 3)
	ratings["Go for Backend"] = 4.7
	ratings["SQL Basics"] = 4.2
	ratings.print()

	var missing map[string]int
	fmt.Println(missing["x"]) // 0, reading a nil map is fine
}
