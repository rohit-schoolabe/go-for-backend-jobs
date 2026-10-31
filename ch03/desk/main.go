package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const copiesFile = "copies.txt"

func main() {
	copies, err := readCopies()
	if err != nil {
		fmt.Println("Warning:", err)
		fmt.Println("Starting with 0 copies.")
	}

	fmt.Println("Shelf desk: copies of 'Go for Backend'")
	for {
		fmt.Println()
		fmt.Println("1. Check copies")
		fmt.Println("2. Return copies")
		fmt.Println("3. Issue copies")
		fmt.Println("4. Exit")

		var choice int
		fmt.Print("Your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("On the shelf:", copies)
		case 2:
			n, err := readCount("How many returned? ")
			if err != nil {
				fmt.Println(err)
				continue
			}
			copies += n
			writeCopies(copies)
			fmt.Println("On the shelf:", copies)
		case 3:
			n, err := readCount("How many to issue? ")
			if err != nil {
				fmt.Println(err)
				continue
			}
			if n > copies {
				fmt.Println("Only", copies, "on the shelf.")
				continue
			}
			copies -= n
			writeCopies(copies)
			fmt.Println("On the shelf:", copies)
		default:
			fmt.Println("Goodbye.")
			return
		}
	}
}

func readCount(prompt string) (int, error) {
	var n int
	fmt.Print(prompt)
	if _, err := fmt.Scan(&n); err != nil {
		return 0, errors.New("please type a whole number")
	}
	if n <= 0 {
		return 0, errors.New("the count must be positive")
	}
	return n, nil
}

func writeCopies(copies int) {
	text := strconv.Itoa(copies)
	os.WriteFile(copiesFile, []byte(text), 0644)
}

func readCopies() (int, error) {
	data, err := os.ReadFile(copiesFile)
	if err != nil {
		return 0, errors.New("no saved count found")
	}
	copies, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, errors.New("saved count is not a number")
	}
	return copies, nil
}
