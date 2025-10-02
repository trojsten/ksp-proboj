package client

import (
	"fmt"
	"io"
	"strings"

	"github.com/trojsten/ksp-proboj/libproboj"
)

// sendCommand sends the command with payload to the runner.
// Formats the command as "COMMAND\nPAYLOAD\n.\n" where PAYLOAD is optional.
func (r Runner) sendCommand(command string, payload string) {
	if payload == "" {
		fmt.Printf("%s\n.\n", command)
	} else {
		fmt.Printf("%s\n%s\n.\n", command, payload)
	}
}

// sendCommandWithArgs sends a command with arguments and optional payload to the runner.
// Joins the arguments with spaces and appends them to the command, then sends using sendCommand.
// Format: "COMMAND ARG1 ARG2...\nPAYLOAD\n.\n"
func (r Runner) sendCommandWithArgs(command string, args []string, payload string) {
	r.sendCommand(fmt.Sprintf("%s %s", command, strings.Join(args, " ")), payload)
}

// readLine reads one line from the runner's stdin.
// Returns the line content and any error encountered.
// Returns io.EOF error if the input stream is closed.
func (r Runner) readLine() (string, error) {
	if !r.scanner.Scan() {
		if r.scanner.Err() != nil {
			return "", r.scanner.Err()
		} else {
			return "", io.EOF
		}
	}
	return r.scanner.Text(), nil
}

// readLines reads multiple lines from the runner until the end-of-transmission mark.
// Reads lines until encountering a line containing only "." which marks the end.
// Returns the concatenated lines with newline separators and any error encountered.
func (r Runner) readLines() (string, error) {
	result := []string{}
	for {
		input, err := r.readLine()
		if err != nil {
			return "", err
		}
		if input == "." {
			break
		}
		result = append(result, input)
	}
	return strings.Join(result, "\n"), nil
}

// readResponse reads a complete protocol response from the runner.
// Reads the status line first, then reads the payload lines until the end marker.
// Returns a libproboj.RunnerResponse containing the status and payload, or any error encountered.
func (r Runner) readResponse() (libproboj.RunnerResponse, error) {
	line, err := r.readLine()
	if err != nil {
		return libproboj.RunnerResponse{}, err
	}

	lines, err := r.readLines()
	if err != nil {
		return libproboj.RunnerResponse{}, err
	}

	return libproboj.RunnerResponse{
		Status:  libproboj.GetStatus(line),
		Payload: lines,
	}, nil
}
