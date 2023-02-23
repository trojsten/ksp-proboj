package libproboj

import (
	"fmt"
	"strings"
)

type Scores map[string]int

// ToObserver sends the given data to the observer
func (r Runner) ToObserver(data string) RunnerResponse {
	r.sendCommand("TO OBSERVER", data)

	response := r.readLine()
	if response == "OK" {
		return Ok
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'TO OBSERVER' from runner: %s", response))
	return Unknown
}

// Scores sends game scores to the observer
func (r Runner) Scores(scores Scores) {
	payload := []string{}
	for player, score := range scores {
		payload = append(payload, fmt.Sprintf("%s %d", player, score))
	}

	r.sendCommand("SCORES", strings.Join(payload, "\n"))
}
