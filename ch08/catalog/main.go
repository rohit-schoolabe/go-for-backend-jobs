package main

import (
	"fmt"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch08/catalog/book"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch08/catalog/member"
)

// saver is anything that can store itself.
type saver interface {
	Save() error
}

// displayer is anything that can show itself.
type displayer interface {
	Display()
}

// outputter can do both.
type outputter interface {
	saver
	displayer
}

func main() {
	b, err := book.New("CS 101", "Go for Backend", "Asha Nair")
	if err != nil {
		fmt.Println(err)
		return
	}
	m, err := member.New("S2291", "Rahul Verma")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, item := range []outputter{b, m} {
		if err := output(item); err != nil {
			fmt.Println("Saving failed:", err)
			return
		}
	}
}

func output(data outputter) error {
	data.Display()
	return saveData(data)
}

func saveData(data saver) error {
	if err := data.Save(); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Println("Saved.")
	return nil
}
