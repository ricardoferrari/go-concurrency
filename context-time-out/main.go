package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(channel chan string, wg *sync.WaitGroup, context context.Context) {
	defer wg.Done()

	for {
		select {
		case msg, ok := <-channel:
			if !ok {
				fmt.Println("Message channel closed, listener exiting")
				return
			}
			fmt.Printf("Received: %s at %v\n", msg, time.Now().Format(time.StampMilli))
			time.Sleep(500 * time.Millisecond) // Simulate processing delay
		case <-context.Done():
			fmt.Println("Listener stopped due to context cancellation")
			return
		}

	}
}

func produce(channel chan string, wg *sync.WaitGroup, message string, context context.Context) {
	defer wg.Done()
	select {
	case <-context.Done():
		fmt.Println("Producer stopped due to context cancellation")
	case channel <- message:
	}
}

func main() {

	var wgListener sync.WaitGroup

	request := make(chan string)

	rootContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer close(request)

	wgListener.Add(2)
	go worker(request, &wgListener, rootContext)

	go produce(request, &wgListener, "Slow request", rootContext)

	time.Sleep(2 * time.Second) // Wait to ensure context timeout occurs

	go produce(request, &wgListener, "Another slow request", rootContext)

	wgListener.Wait()

	// Calculate total execution time

	fmt.Println("All listeners have finished processing")
}
