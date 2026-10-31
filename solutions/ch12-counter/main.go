package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		total int
	)
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				mu.Lock()
				total++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println(total) // 10000, every time
}
