package main

import (
	"errors"
	"fmt"
)

const maxExtension = 14

func extendDue(dueDay *int, days int) error {
	if dueDay == nil {
		return errors.New("no due day to extend")
	}
	if days > maxExtension {
		return fmt.Errorf("cannot extend by %d days, max is %d",
			days, maxExtension)
	}
	*dueDay += days
	return nil
}

func main() {
	due := 120
	fmt.Println(extendDue(&due, 7), due) // <nil> 127
	fmt.Println(extendDue(&due, 30))     // cannot extend ...
	fmt.Println(extendDue(nil, 7))       // no due day to extend
}
