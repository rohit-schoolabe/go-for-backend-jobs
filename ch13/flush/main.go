// Command flush shows the bufio.Writer mistake: data sits
// in the buffer until Flush, and is lost without it.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func size(name string) int64 {
	info, err := os.Stat(name)
	if err != nil {
		return -1
	}
	return info.Size()
}

func main() {
	f, err := os.Create("report.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "CS-101  Ananya  12 days  Rs 24")
	fmt.Println("before Flush:", size("report.txt"), "bytes")

	if err := w.Flush(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("after Flush: ", size("report.txt"), "bytes")
}
