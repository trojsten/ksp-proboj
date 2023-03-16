package main

import (
	"encoding/json"
	"fmt"
	"github.com/trojsten/ksp-proboj/common"
	"os"
	"path"
	"strings"
)

func cmdScores(m *Match, _ []string, payload string) common.RunnerResponse {
	lines := strings.Split(payload, "\n")
	scores := map[string]int{}
	for _, line := range lines {
		var player string
		var score int
		_, err := fmt.Sscanf(line, "%s %d", &player, &score)
		if err != nil {
			m.logger.Warn("Could not parse score data", "err", err)
			return common.RunnerResponse{Status: common.Error}
		}

		scores[player] = score
	}

	fileName := path.Join(m.Game.Gamefolder, "score")
	file, err := os.Create(fileName)
	if err != nil {
		m.logger.Error("Could not open score file", "err", err)
		return common.RunnerResponse{Status: common.Error}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			m.logger.Error("Error while closing score file", "err", err)
		}
	}(file)

	data, err := json.Marshal(scores)
	if err != nil {
		m.logger.Error("Error while marshalling score data", "err", err)
		return common.RunnerResponse{Status: common.Error}
	}
	_, err = file.Write(data)
	if err != nil {
		m.logger.Error("Error while saving score data", "err", err)
		return common.RunnerResponse{Status: common.Error}
	}

	return common.RunnerResponse{Status: common.Ok}
}
