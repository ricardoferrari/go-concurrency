package main

import (
	"fmt"
	"sync"
)

func listen(channel chan string, wg *sync.WaitGroup, id int, listenersActivities *int) {
	defer wg.Done()

	for msg := range channel {
		*listenersActivities++
		fmt.Printf("Listener %d received: %s\n", id, msg)
	}
}

func main() {

	// Create a map to track each listener's activities
	// using a single counter for all listeners avoiding map concurrency issues
	listeners1Activities := 0
	listeners2Activities := 0

	var wgListener sync.WaitGroup

	channel := make(chan string, 3)

	wgListener.Add(2)
	go listen(channel, &wgListener, 1, &listeners1Activities)
	go listen(channel, &wgListener, 2, &listeners2Activities)
	for i := 1; i <= 1000; i++ {
		channel <- fmt.Sprintf("Slow dispatcher %d", i)
	}

	// Close the channel after all messages are sent
	// certifies listeners to finish by ending their range loops
	close(channel)

	wgListener.Wait()

	fmt.Println("All listeners have finished processing")
	fmt.Println("Listeners 1 activities:", listeners1Activities, "Listeners 2 activities:", listeners2Activities)
	fmt.Println("Total activities:", listeners1Activities+listeners2Activities)
}
