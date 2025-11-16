package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	fmt.Println("worker started")

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Worker completed after timeout")
	case <-ctx.Done():
		fmt.Println("worker canceled due to context cancellation: ", ctx.Err())
	}
}

func main() {

	var rootContext context.Context = context.Background()
	ctx, cancel := context.WithCancel(rootContext)
	defer cancel()

	var wgListener sync.WaitGroup

	channel := make(chan string)
	defer close(channel)

	// First worker that will complete after timeout
	wgListener.Add(1)
	go worker(&wgListener, ctx)
	wgListener.Wait()

	// Second worker that will be canceled before completion
	wgListener.Add(1)
	go worker(&wgListener, ctx)
	time.Sleep(2 * time.Second)
	fmt.Println("main canceling context")
	cancel()
	wgListener.Wait()

	fmt.Println("all dispatched")
}
