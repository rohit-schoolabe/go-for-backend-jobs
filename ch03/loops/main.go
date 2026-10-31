package main

import "fmt"

func main() {
	for i := 1; i <= 3; i++ {
		fmt.Println("Reminder", i)
	}

	shelf := 0
	for shelf < 3 {
		shelf++
	}
	fmt.Println("shelves checked:", shelf)

	for day := 1; day <= 7; day++ {
		if day == 7 {
			fmt.Println("Sunday: closed")
			continue
		}
		if day == 3 {
			fmt.Println("Wednesday: stock check, stop here")
			break
		}
		fmt.Println("Day", day, "open")
	}

	today := 6
	isWeekend := today == 6 || today == 7
	if isWeekend {
		fmt.Println("Weekend hours")
	}

	count := 0
desk:
	for {
		count++
		switch {
		case count < 3:
			fmt.Println("serving student", count)
		default:
			break desk // plain break leaves only the switch
		}
	}
	fmt.Println("apple" < "banana", "Zebra" < "apple")
}
