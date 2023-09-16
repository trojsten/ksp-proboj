package client

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Runner struct {
	scanner *bufio.Scanner
}

type RunnerResponse int

const (
	Ok RunnerResponse = iota
	Unknown
	Died
)

func NewRunner() Runner {
	r := Runner{}
	r.scanner = bufio.NewScanner(os.Stdin)
	return r
}

// Log prints the message to the stderr
func (r Runner) Log(message string) {
	_, err := fmt.Fprintf(os.Stderr, "[%s] %s\n", time.Now().Format("15:04:05.000"), message)
	if err != nil {
		panic(err)
	}
}
