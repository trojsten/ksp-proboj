package main

// runSequentially executes games one at a time in a single thread.
// Each game completes before the next one starts.
func runSequentially(config Config, games []Game) {
	for _, game := range games {
		if receivedKillSignal {
			return
		}
		match := NewMatch(config, game)
		match.Run()
	}
}
