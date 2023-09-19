package main

import (
	"encoding/json"
	"fmt"
	"github.com/trojsten/ksp-proboj/libproboj"
	"os"
	"path"
	"strings"
)

func cmdScores(m *Match, _ []string, payload string) libproboj.RunnerResponse {
	lines := strings.Split(payload, "\n")
	scores := map[string]int{}
	for _, line := range lines {
		var player string
		var score int

		_, err := fmt.Sscanf(line, "%s %d", &player, &score)
		if err != nil {
			m.Log.Warn("Could not parse score data", "err", err)
			return libproboj.RunnerResponse{Status: libproboj.Error}
		}

		scores[player] = score
	}

	fileName := path.Join(m.Directory(), "score.json")
	file, err := os.Create(fileName)
	if err != nil {
		m.Log.Error("Could not open score file", "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			m.Log.Error("Error while closing score file", "err", err)
		}
	}(file)

	data, err := json.Marshal(scores)
	if err != nil {
		m.Log.Error("Error while marshalling score data", "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	_, err = file.Write(data)
	if err != nil {
		m.Log.Error("Error while saving score data", "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	return libproboj.RunnerResponse{Status: libproboj.Ok}
}
