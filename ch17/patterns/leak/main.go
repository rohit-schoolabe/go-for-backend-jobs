// Command leak shows a goroutine leak and its fix. Each
// lookup races a slow source against a timeout.
package main

import (
	"fmt"
	"runtime"
	"time"
)

func slowSearch() string {
	time.Sleep(50 * time.Millisecond)
	return "Optics, shelf B4"
}

// leaky gives up after 10ms, but the goroutine it started
// blocks forever on an unbuffered send nobody receives.
func leaky() string {
	ch := make(chan string)
	go func() { ch <- slowSearch() }()
	select {
	case r := <-ch:
		return r
	case <-time.After(10 * time.Millisecond):
		return "timed out"
	}
}

// fixed uses a buffer of one, so the send always succeeds
// and the goroutine exits even when nobody reads.
func fixed() string {
	ch := make(chan string, 1)
	go func() { ch <- slowSearch() }()
	select {
	case r := <-ch:
		return r
	case <-time.After(10 * time.Millisecond):
		return "timed out"
	}
}

func run(name string, lookup func() string) {
	for range 100 {
		lookup()
	}
	time.Sleep(100 * time.Millisecond) // let stragglers end
	fmt.Printf("%-6s goroutines left: %d\n",
		name, runtime.NumGoroutine())
}

func main() {
	run("fixed", fixed)
	run("leaky", leaky)
}
