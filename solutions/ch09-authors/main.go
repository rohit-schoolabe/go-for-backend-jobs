package main

import (
	"fmt"
	"slices"
)

type book struct {
	title, author string
}

func main() {
	books := []book{
		{"Wings of Fire", "Kalam"},
		{"Ignited Minds", "Kalam"},
		{"Go for Backend", "Asha Nair"},
	}

	counts := make(map[string]int)
	for _, b := range books {
		counts[b.author]++
	}

	authors := make([]string, 0, len(counts))
	for a := range counts {
		authors = append(authors, a)
	}
	slices.Sort(authors)
	for _, a := range authors {
		fmt.Println(a, counts[a])
	}
}
