package main

func runSequentially(config Config, games []Game) {
	for _, game := range games {
		match := Match{
			Game:   game,
			Config: config,
		}
		match.Run()
	}
}
