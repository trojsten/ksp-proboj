package main

import (
	"sync"
)

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
