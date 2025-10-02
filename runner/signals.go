package main

import (
	"os"
	"os/signal"
	"syscall"
)

// runningMatches tracks all currently active matches for graceful shutdown.
var runningMatches []*Match

// receivedKillSignal indicates whether a termination signal has been received.
var receivedKillSignal = false

// signalMatchStart registers a new match for signal tracking.
// Called when a match starts to ensure it can be gracefully terminated.
func signalMatchStart(m *Match) {
	runningMatches = append(runningMatches, m)
}

// signalMatchEnd removes a completed match from signal tracking.
// Called when a match finishes to prevent unnecessary cleanup.
func signalMatchEnd(m *Match) {
	for i, match := range runningMatches {
		if match == m {
			runningMatches = append(runningMatches[:i], runningMatches[i+1:]...)
			break
		}
	}
}

// registerSignals sets up signal handlers for graceful shutdown.
// Listens for SIGINT and SIGTERM signals and initiates cleanup of all
// running matches by killing their server processes. Runs in a separate
// goroutine to avoid blocking the main thread.
func registerSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		receivedKillSignal = true
		for _, match := range runningMatches {
			err := match.Server.Kill()
			if err != nil {
				match.Log.Error("Could not kill server", "err", err)
			}
		}
	}()
}
