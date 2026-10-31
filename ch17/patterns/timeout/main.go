// Command timeout passes one deadline down through every
// call, so slow work stops when the caller stops waiting.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// lookupMember pretends to be a slow call to the college's
// student database. It honors ctx.
func lookupMember(ctx context.Context, id string) (string,
	error) {
	select {
	case <-time.After(300 * time.Millisecond):
		return "Ananya, B.Sc. 2nd year", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func issue(ctx context.Context, code, member string) error {
	name, err := lookupMember(ctx, member)
	if err != nil {
		return fmt.Errorf("issuing %s: %w", code, err)
	}
	fmt.Println("issued", code, "to", name)
	return nil
}

func main() {
	for _, limit := range []time.Duration{
		time.Second, 100 * time.Millisecond} {
		ctx, cancel := context.WithTimeout(
			context.Background(), limit)
		err := issue(ctx, "CS-101", "S1042")
		cancel()
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("gave up after", limit, "-", err)
		}
	}
}
