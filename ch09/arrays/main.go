package main

import "fmt"

func main() {
	fines := [4]float64{18, 4, 50, 12}
	fmt.Println(fines) // [18 4 50 12]

	var titles [3]string
	titles[0] = "Go for Backend"
	titles[2] = "SQL Basics"
	fmt.Printf("%q\n", titles)
	// ["Go for Backend" "" "SQL Basics"]
	fmt.Println(fines[2]) // 50

	featured := fines[1:3]
	fmt.Println(featured) // [4 50]

	fmt.Println(fines[:2], fines[2:]) // [18 4] [50 12]
	firstOfFeatured := featured[:1]
	fmt.Println(firstOfFeatured) // [4]

	featured[0] = 5
	fmt.Println(fines) // [18 5 50 12]

	fmt.Println(len(featured), cap(featured)) // 2 3
	featured = featured[:3]
	fmt.Println(featured) // [5 50 12]
}
