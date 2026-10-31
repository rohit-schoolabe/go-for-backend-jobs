package main

import (
	"fmt"
	"time"
)

func notify(name string, done chan bool) {
	fmt.Println("Reminder sent to", name)
	done <- true
}

func slowNotify(name string, done chan bool) {
	time.Sleep(2 * time.Second) // a slow SMS gateway
	fmt.Println("Reminder sent to", name, "(slow)")
	done <- true
}

func countdown(n int, out chan int) {
	for i := n; i > 0; i-- {
		out <- i
	}
	close(out)
}

func main() {
	done := make(chan bool)
	go notify("Asha", done)
	go slowNotify("Rahul", done)
	go notify("Meera", done)
	for range 3 {
		<-done
	}

	dones := make([]chan bool, 2)
	for i, name := range []string{"Kiran", "Dev"} {
		dones[i] = make(chan bool)
		go notify(name, dones[i])
	}
	for _, d := range dones {
		<-d
	}

	numbers := make(chan int)
	go countdown(3, numbers)
	for n := range numbers {
		fmt.Println(n)
	}
}
