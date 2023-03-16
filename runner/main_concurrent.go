package main

import (
	"sync"
)

func parellelWorker(ch <-chan Match, wg *sync.WaitGroup) {
	defer wg.Done()
	for true {
		match, more := <-ch
		if !more {
			return
		}
		match.Run()
	}
}

func runParallel(config Config, games []Game, concurrency int) {
	ch := make(chan Match)
	wg := sync.WaitGroup{}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go parellelWorker(ch, &wg)
	}

	for _, game := range games {
		match := Match{
			Game:   game,
			Config: config,
		}
		ch <- match
	}

	close(ch)
	wg.Wait()
}
