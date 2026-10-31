// Command once loads the holiday list the first time it is
// needed, exactly once, however many goroutines ask.
package main

import (
	"fmt"
	"sync"
)

var (
	holidaysOnce sync.Once
	holidays     map[string]bool
)

func loadHolidays() {
	fmt.Println("loading holidays (runs once)")
	holidays = map[string]bool{
		"2026-01-26": true, "2026-08-15": true,
	}
}

func isHoliday(date string) bool {
	holidaysOnce.Do(loadHolidays)
	return holidays[date]
}

func main() {
	var wg sync.WaitGroup
	for _, d := range []string{"2026-01-26", "2026-01-27",
		"2026-08-15"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(d, isHoliday(d))
		}()
	}
	wg.Wait()
}
