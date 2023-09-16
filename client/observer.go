package client

import (
	"fmt"
	"github.com/trojsten/ksp-proboj/libproboj"
	"strings"
)

type Scores map[string]int

// ToObserver sends the given data to the observer
func (r Runner) ToObserver(data string) RunnerResponse {
	r.sendCommand("TO OBSERVER", data)

	response, err := r.readResponse()
	if err != nil {
		r.Log(fmt.Sprintf("error while reading response: %s", err.Error()))
		return Unknown
	}
	if response.Status == libproboj.Ok {
		return Ok
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'TO OBSERVER' from runner: %s", response.String()))
	return Unknown
}

// Scores sends game scores to the observer
func (r Runner) Scores(scores Scores) {
	payload := []string{}
	for player, score := range scores {
		payload = append(payload, fmt.Sprintf("%s %d", player, score))
	}

	r.sendCommand("SCORES", strings.Join(payload, "\n"))
	response, err := r.readResponse()
	if err != nil {
		r.Log(fmt.Sprintf("error while reading response: %s", err.Error()))
		return
	}
	if response.Status != libproboj.Ok {
		r.Log(fmt.Sprintf("unknown response to cmd 'SCORES' from runner: %s", response.String()))
	}
}
