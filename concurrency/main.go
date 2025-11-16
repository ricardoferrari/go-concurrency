package main

import (
	"fmt"
	"sync"
)

func safeIncrement(counter *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	*counter++
	mu.Unlock()
}

func main() {

	var wg sync.WaitGroup

	var counter int = 0

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("Final Counter:", counter)

	var mu sync.Mutex
	counter = 0

	for range 1000 {
		wg.Add(1)
		go safeIncrement(&counter, &mu, &wg)
	}

	wg.Wait()
	fmt.Println("Final Counter with Mutex:", counter)
}
