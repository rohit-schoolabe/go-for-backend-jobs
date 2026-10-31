// Command ticker spaces out work with time.Ticker:
// at most one send every 200 ms, whatever the loop does.
package main

import (
	"fmt"
	"time"
)

func main() {
	members := []string{"Asha", "Ravi", "Meena", "Joseph"}
	tick := time.NewTicker(200 * time.Millisecond)
	defer tick.Stop() // a ticker runs until stopped

	start := time.Now()
	for _, m := range members {
		<-tick.C
		fmt.Printf("%4dms reminder to %s\n",
			time.Since(start).Milliseconds(), m)
	}
}
