package main

import (
	"os"
	"os/signal"
	"syscall"
)

var runningMatches []*Match
var receivedKillSignal = false

func signalMatchStart(m *Match) {
	runningMatches = append(runningMatches, m)
}

func signalMatchEnd(m *Match) {
	for i, match := range runningMatches {
		if match == m {
			runningMatches = append(runningMatches[:i], runningMatches[i+1:]...)
			break
		}
	}
}

func registerSignals() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		receivedKillSignal = true
		for _, match := range runningMatches {
			err := match.Server.Kill()
			if err != nil {
				match.logger.Error("Could not kill server", "err", err)
			}
		}
	}()
}
