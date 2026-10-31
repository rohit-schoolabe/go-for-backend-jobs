// Package words holds the Chapter 13 exercise answers.
package words

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch13/catalogio"
)

// countWords counts words separated by any white space.
func countWords(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}

// ReadLines reads one JSON book per line, the format
// that WriteLines produces.
func ReadLines(r io.Reader,
	fn func(catalogio.Book) error) error {
	dec := json.NewDecoder(r)
	for {
		var b catalogio.Book
		err := dec.Decode(&b)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := fn(b); err != nil {
			return err
		}
	}
}
