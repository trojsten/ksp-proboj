package client

import (
	"fmt"

	"github.com/trojsten/ksp-proboj/libproboj"
)

// ToPlayer sends the given data to the player.
// Optional comment can be provided, which will get logged by
// the runner
func (r Runner) ToPlayer(player string, comment string, data string) RunnerResponse {
	r.sendCommandWithArgs("TO PLAYER", []string{player, comment}, data)

	response, err := r.readResponse()
	if err != nil {
		r.Log(fmt.Sprintf("error while reading response: %s", err.Error()))
		return Unknown
	}
	if response.Status == libproboj.Ok {
		return Ok
	} else if response.Status == libproboj.Died {
		return Died
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'TO PLAYER' from runner: %s", response.String()))
	return Unknown
}

// ReadPlayer reads all data from the player until end-of-transmittion mark
func (r Runner) ReadPlayer(player string) (RunnerResponse, string) {
	return r.ReadPlayerWithTimeout(player, 1.0)
}

// ReadPlayerWithTimeout reads all data from the player until end-of-transmittion mark
// with an optional timeout multiplier. multiplier=1.0 uses default timeout, multiplier=2.0 doubles it, etc.
func (r Runner) ReadPlayerWithTimeout(player string, timeoutMultiplier float64) (RunnerResponse, string) {
	args := []string{player}
	if timeoutMultiplier != 1.0 {
		args = append(args, fmt.Sprintf("%.1f", timeoutMultiplier))
	}
	r.sendCommandWithArgs("READ PLAYER", args, "")

	response, err := r.readResponse()
	if err != nil {
		r.Log(fmt.Sprintf("error while reading response: %s", err.Error()))
		return Unknown, ""
	}

	if response.Status == libproboj.Ok {
		return Ok, response.Payload
	} else if response.Status == libproboj.Died {
		return Died, ""
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'READ PLAYER' from runner: %s", response.String()))
	return Unknown, ""
}

// KillPlayer instructs the runner to kill the player's process
func (r Runner) KillPlayer(player string) RunnerResponse {
	r.sendCommandWithArgs("KILL PLAYER", []string{player}, "")

	response, err := r.readResponse()
	if err != nil {
		r.Log(fmt.Sprintf("error while reading response: %s", err.Error()))
		return Unknown
	}

	if response.Status == libproboj.Ok {
		return Ok
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'KILL PLAYER' from runner: %s", response.String()))
	return Unknown
}

// PausePlayer instructs the runner to pause player's process
func (r Runner) PausePlayer(player string) RunnerResponse {
	r.sendCommandWithArgs("PAUSE PLAYER", []string{player}, "")

	response, err := r.readResponse()
	if err != nil {
		r.Log(fmt.Sprintf("error while reading response: %s", err.Error()))
		return Unknown
	}

	if response.Status == libproboj.Ok {
		return Ok
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'PAUSE PLAYER' from runner: %s", response.String()))
	return Unknown
}

// ResumePlayer instructs the runner to resume player's process
func (r Runner) ResumePlayer(player string) RunnerResponse {
	r.sendCommandWithArgs("RESUME PLAYER", []string{player}, "")

	response, err := r.readResponse()
	if err != nil {
		r.Log(fmt.Sprintf("error while reading response: %s", err.Error()))
		return Unknown
	}

	if response.Status == libproboj.Ok {
		return Ok
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'RESUME PLAYER' from runner: %s", response.String()))
	return Unknown
}
