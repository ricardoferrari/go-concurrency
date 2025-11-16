package main

import (
	"fmt"
	"sync"
)

func counterManager(sumRequest <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	var counter int = 0
	for increment := range sumRequest {
		counter += increment
	}
	fmt.Println("Final Counter from Manager:", counter)
}

func main() {

	var wg sync.WaitGroup
	var workersWg sync.WaitGroup

	var sumRequest = make(chan int)
	wg.Add(1)
	go counterManager(sumRequest, &wg)

	for range 10000 {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			sumRequest <- 1
		}()
	}

	workersWg.Wait()

	close(sumRequest)
	wg.Wait()
}
