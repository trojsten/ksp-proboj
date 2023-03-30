package libproboj

import (
	"fmt"
	"github.com/trojsten/ksp-proboj/common"
	"io"
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

// readLines reads multiple lines from the runner until the end-of-transmittion mark
func (r Runner) readLines() (string, error) {
	result := []string{}
	for true {
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

// readLines reads multiple lines from the runner until the end-of-transmittion mark
func (r Runner) readResponse() (common.RunnerResponse, error) {
	line, err := r.readLine()
	if err != nil {
		return common.RunnerResponse{}, err
	}

	lines, err := r.readLines()
	if err != nil {
		return common.RunnerResponse{}, err
	}

	return common.RunnerResponse{
		Status:  common.GetStatus(line),
		Payload: lines,
	}, nil
}
