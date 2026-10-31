package main

import (
	"fmt"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch04/desk/fileops"
)

const copiesFile = "copies.txt"

func main() {
	copies, err := fileops.ReadInt(copiesFile, 0)
	if err != nil {
		fmt.Println("Warning:", err)
	}

	for {
		showMenu()
		var choice int
		fmt.Print("Your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("On the shelf:", copies)
			continue
		case 2:
			n, err := readCount("How many returned? ")
			if err != nil {
				fmt.Println(err)
				continue
			}
			copies += n
		case 3:
			n, err := readCount("How many to issue? ")
			if err != nil {
				fmt.Println(err)
				continue
			}
			if n > copies {
				fmt.Println("Only", copies, "on the shelf.")
				continue
			}
			copies -= n
		case 4:
			fmt.Println("Goodbye.")
			return
		default:
			fmt.Println("Unknown choice.")
			continue
		}

		err := fileops.WriteInt(copiesFile, copies)
		if err != nil {
			fmt.Println("Could not save:", err)
		}
		fmt.Println("On the shelf:", copies)
	}
}
