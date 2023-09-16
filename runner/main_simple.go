package main

func runSequentially(config Config, games []Game) {
	for _, game := range games {
		if receivedKillSignal {
			return
		}
		match := NewMatch(config, game)
		match.Run()
	}
}
