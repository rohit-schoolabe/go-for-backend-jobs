// Command count reports lines and bytes from any reader.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// countLines works on a file, a network body, a string,
// or anything else with a Read method.
func countLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}

func main() {
	n, _ := countLines(strings.NewReader("a\nb\nc\n"))
	fmt.Println("from a string:", n)

	f, err := os.Open("count.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Count the file and copy it to the screen in one pass.
	var bytes countingWriter
	tee := io.TeeReader(f, io.MultiWriter(os.Stdout, &bytes))
	n, err = countLines(tee)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("from the file: %d lines, %d bytes\n",
		n, bytes.n)
}

// countingWriter is an io.Writer that only counts.
type countingWriter struct{ n int }

func (c *countingWriter) Write(p []byte) (int, error) {
	c.n += len(p)
	return len(p), nil
}
