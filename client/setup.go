package client

import (
	"fmt"
	"strings"
)

// ReadConfig reads the initial configuration from the runner.
// Returns a slice of player names and the configuration data string.
func (r Runner) ReadConfig() ([]string, string) {
	line, err := r.readLine()
	if err != nil {
		panic(fmt.Errorf("error while reading config: %s", err.Error()))
	}
	if line != "CONFIG" {
		panic(fmt.Errorf("expected CONFIG, got %s", line))
	}

	pl, err := r.readLine()
	if err != nil {
		panic(fmt.Errorf("error while reading config: %s", err.Error()))
	}

	players := strings.Split(pl, " ")
	data, err := r.readLines()
	if err != nil {
		panic(fmt.Errorf("error while reading config: %s", err.Error()))
	}

	return players, data
}

// End sends an END command to the runner to signal the end of the game.
// This notifies the runner that the game server has finished execution.
func (r Runner) End() {
	r.sendCommand("END", "")
}
