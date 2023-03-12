package main

import "github.com/trojsten/ksp-proboj/common"

func cmdKillPlayer(m *Match, args []string, _ string) common.RunnerResponse {
	if len(args) < 1 {
		m.logger.Error("Invalid command syntax: missing arguments")
		return common.RunnerResponse{Status: common.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.logger.Error("Unknown player", "player", player)
		return common.RunnerResponse{Status: common.Error}
	}

	if proc.IsRunning() {
		err := proc.Kill()
		if err != nil {
			m.logger.Error("Failed to kill player", "player", player, "err", err)
			return common.RunnerResponse{Status: common.Error}
		}
	}

	return common.RunnerResponse{Status: common.Ok}
}
