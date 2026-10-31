// Package conversion turns text into numbers.
package conversion

import (
	"fmt"
	"strconv"
	"strings"
)

// StringsToInts parses every line as a whole number.
func StringsToInts(lines []string) ([]int, error) {
	out := make([]int, 0, len(lines))
	for i, line := range lines {
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		out = append(out, n)
	}
	return out, nil
}
