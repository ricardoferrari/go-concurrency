package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func producer(wg *sync.WaitGroup, locker sync.Locker) {
	defer wg.Done()
	locker.Lock()
	defer locker.Unlock()
	// *counter++
}

func consumer(wg *sync.WaitGroup, locker sync.Locker) {
	defer wg.Done()
	locker.Lock()
	defer locker.Unlock()
}

func main() {
	test := func(count int, mutex, rwmutex sync.Locker) time.Duration {
		var wg sync.WaitGroup

		// counter := 1

		start := time.Now()

		for range count {
			wg.Add(2)
			go producer(&wg, mutex)
			go consumer(&wg, rwmutex)
		}
		wg.Wait()
		return time.Since(start)
	}

	var mutex sync.Mutex = sync.Mutex{}
	var rwMutex sync.RWMutex = sync.RWMutex{}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Count\tMutex Duration\tRWMutex Duration")
	tw.Flush()
	for i := 0; i <= 8; i++ {
		count := int(math.Pow(10, float64(i)))
		mutexDuration := test(count, &mutex, &mutex)
		rwMutexDuration := test(count, &mutex, &rwMutex)
		fmt.Fprintf(tw, "%d\t%v\t%v\n", count, mutexDuration, rwMutexDuration)
		tw.Flush()
	}

}
