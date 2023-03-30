package main

func runSequentially(config Config, games []Game) {
	for _, game := range games {
		if receivedKillSignal {
			return
		}
		match := &Match{
			Game:   game,
			Config: config,
		}
		signalMatchStart(match)
		match.Run()
		signalMatchEnd(match)
	}
}
