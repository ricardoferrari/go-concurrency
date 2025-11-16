package main

import (
	"fmt"
	"sync"
	"time"
)

func listen(channel chan string, wg *sync.WaitGroup, id int, listenersActivities *int, rateLimiter *time.Ticker) {
	defer wg.Done()
	fmt.Println("Listener", id, "started with ticker", rateLimiter.C)

	for msg := range channel {
		// Rate limiting: wait for the ticker before processing the next message
		<-rateLimiter.C
		time.Sleep(500 * time.Millisecond) // Simulate rate limiting delay
		*listenersActivities++
		fmt.Printf("Listener %d received: %s at %v\n", id, msg, time.Now().Format(time.StampMilli))
	}
}

func main() {
	// Start timer to measure total execution time
	startTime := time.Now()

	rateLimiter := time.NewTicker(time.Second)
	defer rateLimiter.Stop()

	// Create a map to track each listener's activities
	// using a single counter for all listeners avoiding map concurrency issues
	listeners1Activities := 0
	listeners2Activities := 0
	listeners3Activities := 0

	var wgListener sync.WaitGroup

	request := make(chan string, 2)

	wgListener.Add(3)
	go listen(request, &wgListener, 1, &listeners1Activities, rateLimiter)
	go listen(request, &wgListener, 2, &listeners2Activities, rateLimiter)
	go listen(request, &wgListener, 3, &listeners3Activities, rateLimiter)

	// Although we spawned 3 listeners, we will send 100 requests.
	// With a rate limit of 1 request per second, this should take about 100 seconds.

	for i := 1; i <= 8; i++ {
		request <- fmt.Sprintf("Slow request %d", i)
	}

	// Close the channel after all messages are sent
	// certifies listeners to finish by ending their range loops
	close(request)

	wgListener.Wait()

	// Calculate total execution time
	elapsedTime := time.Since(startTime)

	fmt.Println("All listeners have finished processing")
	fmt.Println("Listeners 1 activities:", listeners1Activities, "Listeners 2 activities:", listeners2Activities, "Listeners 3 activities:", listeners3Activities)
	fmt.Println("Total activities:", listeners1Activities+listeners2Activities+listeners3Activities)
	fmt.Printf("Total execution time: %v\n", elapsedTime)
}
