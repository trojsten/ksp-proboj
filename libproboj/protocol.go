package libproboj

import (
	"fmt"
	"github.com/trojsten/ksp-proboj/common"
	"strings"
)

// sendCommand sends the command with payload to the runner
func (r Runner) sendCommand(command string, payload string) {
	if payload == "" {
		fmt.Printf("%s\n.\n", command)
	} else {
		fmt.Printf("%s\n%s\n.\n", command, payload)
	}
}

func (r Runner) sendCommandWithArgs(command string, args []string, payload string) {
	r.sendCommand(fmt.Sprintf("%s %s", command, strings.Join(args, " ")), payload)
}

// readLine reads one line from the runner
func (r Runner) readLine() string {
	r.scanner.Scan()
	return r.scanner.Text()
}

// readLines reads multiple lines from the runner until the end-of-transmittion mark
func (r Runner) readLines() string {
	result := []string{}
	for true {
		input := r.readLine()
		if input == "." {
			break
		}
		result = append(result, input)
	}
	return strings.Join(result, "\n")
}

// readLines reads multiple lines from the runner until the end-of-transmittion mark
func (r Runner) readResponse() common.RunnerResponse {
	return common.RunnerResponse{
		Status:  common.GetStatus(r.readLine()),
		Payload: r.readLines(),
	}
}
