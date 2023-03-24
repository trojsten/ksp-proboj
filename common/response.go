package common

import "fmt"

type Status int

const (
	Error Status = iota
	Ok
	Died
	Ignore // this response will not be sent to the server
)

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

type RunnerResponse struct {
	Status  Status
	Payload string
}

func (r RunnerResponse) String() string {
	if r.Payload == "" {
		return fmt.Sprintf("%s\n", r.Status.String())
	}
	return fmt.Sprintf("%s\n%s\n.\n", r.Status.String(), r.Payload)
}
