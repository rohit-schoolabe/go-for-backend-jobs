// Command oncevalue answers Chapter 17, exercise 4.
package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// holidays loads the file once. A missing file is an
// error every caller sees, and it is not retried.
var holidays = sync.OnceValues(func() (map[string]bool,
	error) {
	f, err := os.Open("holidays.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m := map[string]bool{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		m[sc.Text()] = true
	}
	return m, sc.Err()
})

func isHoliday(date string) (bool, error) {
	m, err := holidays()
	if err != nil {
		return false, fmt.Errorf("holiday list: %w", err)
	}
	return m[date], nil
}

func main() {
	for _, d := range []string{"2026-01-26", "2026-01-27"} {
		ok, err := isHoliday(d)
		fmt.Println(d, ok, err)
	}
}
