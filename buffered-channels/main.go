package main

import (
	"fmt"
	"sync"
	"time"
)

func dispatch(channel chan string, text string) {
	time.Sleep(10 * time.Second)
	channel <- text
}

func listen(channel chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case msg := <-channel:
			fmt.Println("received:", msg)
		case <-time.After(5 * time.Second):
			fmt.Println("listener timed out")
			return
		}
	}
}

func main() {

	var wgListener sync.WaitGroup

	channel := make(chan string, 3)
	defer close(channel)

	wgListener.Add(1)
	go listen(channel, &wgListener)

	go dispatch(channel, "Slow dispatcher")

	wgListener.Wait()

	fmt.Println("all dispatched")
}
