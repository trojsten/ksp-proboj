package main

import (
	"sync"
)

// parellelWorker is a goroutine that processes matches from a channel.
// Each worker runs matches sequentially as they arrive from the channel.
func parellelWorker(ch <-chan *Match, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		match, more := <-ch
		if !more {
			return
		}
		if receivedKillSignal {
			continue
		}
		match.Run()
	}
}

// runParallel executes multiple games concurrently using worker pool pattern.
// Creates specified number of worker goroutines that process matches from a shared channel.
func runParallel(config Config, games []Game, concurrency int) {
	ch := make(chan *Match)
	wg := sync.WaitGroup{}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go parellelWorker(ch, &wg)
	}

	for _, game := range games {
		if receivedKillSignal {
			break
		}
		match := NewMatch(config, game)
		ch <- match
	}

	close(ch)
	wg.Wait()
}
