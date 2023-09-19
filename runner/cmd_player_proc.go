package main

import "github.com/trojsten/ksp-proboj/libproboj"

func cmdKillPlayer(m *Match, args []string, _ string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.Log.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.Log.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if !proc.IsRunning() {
		return libproboj.RunnerResponse{Status: libproboj.Ok}
	}

	err := proc.Kill()
	if err != nil {
		m.Log.Error("Failed to kill player", "player", player, "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	return libproboj.RunnerResponse{Status: libproboj.Ok}
}

func cmdPausePlayer(m *Match, args []string, _ string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.Log.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.Log.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if !proc.IsRunning() {
		return libproboj.RunnerResponse{Status: libproboj.Ok}
	}

	err := proc.Pause()
	if err != nil {
		m.Log.Error("Failed to pause player", "player", player, "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	return libproboj.RunnerResponse{Status: libproboj.Ok}
}

func cmdResumePlayer(m *Match, args []string, _ string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.Log.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.Log.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if !proc.IsRunning() {
		return libproboj.RunnerResponse{Status: libproboj.Ok}
	}

	err := proc.Resume()
	if err != nil {
		m.Log.Error("Failed to resume player", "player", player, "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	return libproboj.RunnerResponse{Status: libproboj.Ok}
}
