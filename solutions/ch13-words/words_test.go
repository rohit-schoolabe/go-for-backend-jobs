package words

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch13/catalogio"
)

func TestCountWords(t *testing.T) {
	n, err := countWords(strings.NewReader(
		"Wings of Fire\n  Ignited\tMinds "))
	if err != nil || n != 5 {
		t.Fatalf("got %d, %v", n, err)
	}
}

func TestRoundTrip(t *testing.T) {
	in := []catalogio.Book{{Code: "A", Title: "Optics"},
		{Code: "B", Title: "Go Basics"}}
	var buf bytes.Buffer
	if err := catalogio.WriteLines(&buf, in); err != nil {
		t.Fatal(err)
	}
	var out []catalogio.Book
	err := ReadLines(&buf, func(b catalogio.Book) error {
		out = append(out, b)
		return nil
	})
	if err != nil || len(out) != 2 || out[1] != in[1] {
		t.Fatalf("got %v, %v", out, err)
	}
}
