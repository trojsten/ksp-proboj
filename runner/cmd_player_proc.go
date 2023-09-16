package main

import "github.com/trojsten/ksp-proboj/libproboj"

func cmdKillPlayer(m *Match, args []string, _ string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.logger.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.logger.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if proc.IsRunning() {
		err := proc.Kill()
		if err != nil {
			m.logger.Error("Failed to kill player", "player", player, "err", err)
			return libproboj.RunnerResponse{Status: libproboj.Error}
		}
	}

	return libproboj.RunnerResponse{Status: libproboj.Ok}
}
