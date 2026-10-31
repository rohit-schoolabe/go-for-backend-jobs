package main

import (
	"fmt"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch12/finereport/filemanager"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch12/finereport/fines"
)

func main() {
	start := time.Now()
	rates := []float64{1, 2, 5}
	doneChans := make([]chan bool, len(rates))
	errorChans := make([]chan error, len(rates))

	for i, rate := range rates {
		doneChans[i] = make(chan bool)
		errorChans[i] = make(chan error)

		out := fmt.Sprintf("fines_%.0f.json", rate)
		fm := filemanager.New("overdue.txt", out)
		job := fines.NewJob(fm, rate)
		go job.Process(doneChans[i], errorChans[i])
	}

	for i := range rates {
		select {
		case err := <-errorChans[i]:
			fmt.Println("Failed:", err)
		case <-doneChans[i]:
			fmt.Printf("Rate %.0f done\n", rates[i])
		}
	}
	fmt.Println("Took", time.Since(start).Round(time.Second))
}
