// Package libproboj provides shared protocol definitions and utilities for the Proboj system.
package libproboj

import "fmt"

// Status represents the result status of a protocol command execution.
// These status codes are used in RunnerResponse to indicate the outcome
// of operations performed by the runner.
type Status int

const (
	// Error indicates that the command failed due to invalid syntax,
	// process error, or other execution problems.
	Error Status = iota
	// Ok indicates that the command completed successfully.
	Ok
	// Died indicates that the player process terminated during
	// the command execution (e.g., crashed or was killed).
	Died
	// Ignore indicates that this response should not be sent to the server.
	// Used internally for filtering responses.
	Ignore
)

// String returns the string representation of the Status code.
func (s Status) String() string {
	switch s {
	case Error:
		return "ERROR"
	case Ok:
		return "OK"
	case Died:
		return "DIED"
	}
	return "?"
}

// GetStatus converts a string representation back to a Status enum value.
func GetStatus(s string) Status {
	switch s {
	case "ERROR":
		return Error
	case "OK":
		return Ok
	case "DIED":
		return Died
	}
	return Error
}

// RunnerResponse represents a response from the runner to the game server.
// It contains a status code indicating the result of the operation and
// a payload string containing any additional data or error messages.
type RunnerResponse struct {
	Status  Status // The result status of the operation (Error, Ok, Died)
	Payload string // Additional data or error message content
}

// String formats the RunnerResponse as a protocol message.
// The format is: "STATUS\nPAYLOAD\n.\n" where STATUS is the string
// representation of the status code, PAYLOAD is the response content,
// and the response ends with a line containing only ".".
func (r RunnerResponse) String() string {
	return fmt.Sprintf("%s\n%s\n.\n", r.Status.String(), r.Payload)
}
