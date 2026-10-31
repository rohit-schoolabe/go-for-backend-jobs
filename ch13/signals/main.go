// Command signals keeps working until Ctrl+C, then
// finishes the current item and exits cleanly.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("stopping after", i-1, "items")
			return
		case <-time.After(500 * time.Millisecond):
			fmt.Println("processed item", i)
		}
	}
}
