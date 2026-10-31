package catalogio

import (
	"bytes"
	"strings"
	"testing"
)

const export = `[
 {"code":"CS-101","title":"Go Basics",
  "author":"Rao","year":2023},
 {"code":"PH-201","title":"Optics","author":"Iyer","year":2019}
]`

func TestReadArray(t *testing.T) {
	var titles []string
	err := ReadArray(strings.NewReader(export),
		func(b Book) error {
			titles = append(titles, b.Title)
			return nil
		})
	if err != nil {
		t.Fatal(err)
	}
	if len(titles) != 2 || titles[1] != "Optics" {
		t.Fatalf("got %v", titles)
	}
}

func TestReadArrayBadBook(t *testing.T) {
	bad := `[{"code":"X","year":"old"}]`
	err := ReadArray(strings.NewReader(bad),
		func(Book) error { return nil })
	if err == nil {
		t.Fatal("want an error for a string year")
	}
}

func TestWriteLines(t *testing.T) {
	var buf bytes.Buffer
	books := []Book{{Code: "A"}, {Code: "B"}}
	if err := WriteLines(&buf, books); err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(buf.String(), "\n"); got != 2 {
		t.Fatalf("want 2 lines, got %d", got)
	}
}
