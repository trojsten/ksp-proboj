package libproboj

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

type RunnerResponse struct {
	Status  Status
	Payload string
}

func (r RunnerResponse) String() string {
	return fmt.Sprintf("%s\n%s\n.\n", r.Status.String(), r.Payload)
}
