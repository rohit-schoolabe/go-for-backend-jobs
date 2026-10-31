// Package catalogio reads and writes the catalog as a
// stream, one book at a time, whatever its size.
package catalogio

import (
	"encoding/json"
	"fmt"
	"io"
)

// Book is one entry in the old system's export.
type Book struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// ReadArray decodes a JSON array of books from r and calls
// fn for each one, without loading the whole array.
func ReadArray(r io.Reader, fn func(Book) error) error {
	dec := json.NewDecoder(r)
	if _, err := dec.Token(); err != nil { // the [
		return fmt.Errorf("reading start of array: %w", err)
	}
	for i := 0; dec.More(); i++ {
		var b Book
		if err := dec.Decode(&b); err != nil {
			return fmt.Errorf("book %d: %w", i, err)
		}
		if err := fn(b); err != nil {
			return err
		}
	}
	_, err := dec.Token() // the ]
	return err
}

// WriteLines writes one JSON object per line, a format
// that other tools can read line by line.
func WriteLines(w io.Writer, books []Book) error {
	enc := json.NewEncoder(w)
	for _, b := range books {
		if err := enc.Encode(b); err != nil {
			return err
		}
	}
	return nil
}
