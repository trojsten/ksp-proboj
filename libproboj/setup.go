package libproboj

import (
	"fmt"
	"strings"
)

func (r Runner) ReadConfig() ([]string, string) {
	line, err := r.readLine()
	if err != nil {
		panic(fmt.Errorf("error while reading config: %s", err.Error()))
	}
	if line != "CONFIG" {
		panic(fmt.Errorf("expected CONFIG, got %s", line))
	}

	pl, err := r.readLine()
	if err != nil {
		panic(fmt.Errorf("error while reading config: %s", err.Error()))
	}

	players := strings.Split(pl, " ")
	data, err := r.readLines()
	if err != nil {
		panic(fmt.Errorf("error while reading config: %s", err.Error()))
	}
	
	return players, data
}

func (r Runner) End() {
	r.sendCommand("END", "")
}
