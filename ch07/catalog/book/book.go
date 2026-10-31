// Package book holds one catalog entry and saves it.
package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Book is one title in the catalog.
type Book struct {
	Code    string    `json:"code"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	AddedAt time.Time `json:"added_at"`
}

// New checks that every field has a value.
func New(code, title, author string) (Book, error) {
	if code == "" || title == "" || author == "" {
		return Book{}, errors.New("code, title, author needed")
	}
	return Book{
		Code:    code,
		Title:   title,
		Author:  author,
		AddedAt: time.Now(),
	}, nil
}

// Display prints the book for the person at the desk.
func (b Book) Display() {
	fmt.Printf("[%s] %s, by %s\n", b.Code, b.Title, b.Author)
}

// FileName turns "CS 101" into "cs_101.json".
func (b Book) FileName() string {
	name := strings.ReplaceAll(b.Code, " ", "_")
	return strings.ToLower(name) + ".json"
}

// Save writes the book as JSON next to the program.
func (b Book) Save() error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return os.WriteFile(b.FileName(), data, 0644)
}
