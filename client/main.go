// Package client provides a game server client library for communicating with the Proboj runner.
// This package enables game servers to interact with the runner through a simple protocol.
package client

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Runner represents the interface for communicating with the Proboj runner.
// It handles protocol communication, logging, and provides methods for
// interacting with players and observers. The Runner reads from stdin
// and writes to stdout for communication with the runner process.
type Runner struct {
	scanner *bufio.Scanner // Scanner for reading responses from the runner
}

// RunnerResponse represents the response status from runner operations.
type RunnerResponse int

const (
	// Ok indicates that the command completed successfully.
	Ok RunnerResponse = iota
	// Unknown indicates an unknown or unexpected response status.
	Unknown
	// Died indicates that the player process terminated during command execution.
	Died
)

// NewRunner creates and initializes a new Runner instance.
func NewRunner() Runner {
	r := Runner{}
	r.scanner = bufio.NewScanner(os.Stdin)
	return r
}

// Log prints a timestamped message to stderr for debugging and logging purposes.
func (r Runner) Log(message string) {
	_, err := fmt.Fprintf(os.Stderr, "[%s] %s\n", time.Now().Format("15:04:05.000"), message)
	if err != nil {
		panic(err)
	}
}
