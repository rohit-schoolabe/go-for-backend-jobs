// This program deadlocks on purpose. Run it to see the error.
package main

import "fmt"

func main() {
	done := make(chan bool)
	done <- true // nobody is receiving
	fmt.Println("never printed")
}
