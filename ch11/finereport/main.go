package main

import (
	"fmt"
	"os"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch11/finereport/cmdmanager"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch11/finereport/filemanager"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch11/finereport/fines"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch11/finereport/iomanager"
)

func main() {
	rates := []float64{1, 2, 5}
	interactive := len(os.Args) > 1 && os.Args[1] == "-i"

	for _, rate := range rates {
		var io iomanager.IOManager
		if interactive {
			io = cmdmanager.New()
		} else {
			out := fmt.Sprintf("fines_%.0f.json", rate)
			io = filemanager.New("overdue.txt", out)
		}

		job := fines.NewJob(io, rate)
		if err := job.Process(); err != nil {
			fmt.Println("Could not process job:", err)
			return
		}
	}
	fmt.Println("Done.")
}
