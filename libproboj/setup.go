package libproboj

import (
	"fmt"
	"strings"
)

func (r Runner) ReadConfig() ([]string, string) {
	line := r.readLine()
	if line != "CONFIG" {
		panic(fmt.Errorf("expected CONFIG, got %s", line))
	}

	players := strings.Split(r.readLine(), " ")
	data := r.readLines()
	return players, data
}

func (r Runner) End() {
	r.sendCommand("END", "")
}
