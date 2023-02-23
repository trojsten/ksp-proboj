package libproboj

import "fmt"

// ToPlayer sends the given data to the player.
// Optional comment can be provided, which will get logged by
// the runner
func (r Runner) ToPlayer(player string, comment string, data string) RunnerResponse {
	r.sendCommandWithArgs("TO PLAYER", []string{player, comment}, data)

	response := r.readLine()
	if response == "OK" {
		return Ok
	} else if response == "DIED" {
		return Died
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'TO PLAYER' from runner: %s", response))
	return Unknown
}

// ReadPlayer reads all data from the player until end-of-transmittion mark
func (r Runner) ReadPlayer(player string) (RunnerResponse, string) {
	r.sendCommandWithArgs("READ PLAYER", []string{player}, "")

	response := r.readLine()
	if response == "OK" {
		return Ok, r.readLines()
	} else if response == "DIED" {
		return Died, ""
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'READ PLAYER' from runner: %s", response))
	return Unknown, ""
}

// KillPlayer instructs the runner to kill the player's process
func (r Runner) KillPlayer(player string) RunnerResponse {
	r.sendCommandWithArgs("KILL PLAYER", []string{player}, "")

	response := r.readLine()
	if response == "OK" {
		return Ok
	}
	r.Log(fmt.Sprintf("unknown response to cmd 'KILL PLAYER' from runner: %s", response))
	return Unknown
}
