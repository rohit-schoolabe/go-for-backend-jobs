// Package bufpool compares building reminder emails with a
// new buffer each time and with buffers from a sync.Pool.
package bufpool

import (
	"bytes"
	"fmt"
	"sync"
)

// Render allocates a fresh buffer on every call.
func Render(name, title string, days int) string {
	var b bytes.Buffer
	write(&b, name, title, days)
	return b.String()
}

var pool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// RenderPooled reuses buffers. Reset before use: a pooled
// buffer still holds the last caller's text.
func RenderPooled(name, title string, days int) string {
	b := pool.Get().(*bytes.Buffer)
	b.Reset()
	defer pool.Put(b)
	write(b, name, title, days)
	return b.String()
}

func write(b *bytes.Buffer, name, title string, days int) {
	fmt.Fprintf(b, "Dear %s,\n\n", name)
	fmt.Fprintf(b, "%q is %d days overdue. ", title, days)
	b.WriteString("Please return it to the desk.\n\n")
	b.WriteString("College Library\n")
}
