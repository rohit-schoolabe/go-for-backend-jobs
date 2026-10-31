package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch07/catalog/book"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	code, _ := readLine("Book code: ")
	title, _ := readLine("Title: ")
	author, err := readLine("Author: ")
	if err != nil {
		fmt.Println("Could not read input:", err)
		return
	}

	b, err := book.New(code, title, author)
	if err != nil {
		fmt.Println(err)
		return
	}
	b.Display()

	if err := b.Save(); err != nil {
		fmt.Println("Saving failed:", err)
		return
	}
	fmt.Println("Saved to " + b.FileName())
}

// readLine reads a whole line, spaces included.
func readLine(prompt string) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r") // Windows
	return text, nil
}
