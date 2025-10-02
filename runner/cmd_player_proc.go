package main

import "github.com/trojsten/ksp-proboj/libproboj"

// cmdKillPlayer handles the "KILL PLAYER" command from the server.
// It terminates the specified player's process. Returns OK on success,
// ERROR if the player doesn't exist or kill fails, or OK if the
// player is already not running.
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

// cmdPausePlayer handles the "PAUSE PLAYER" command from the server.
// It suspends the specified player's process execution (SIGSTOP on Unix).
// Returns OK on success, ERROR if the player doesn't exist or pause fails,
// or OK if the player is already not running. Platform-dependent behavior.
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

// cmdResumePlayer handles the "RESUME PLAYER" command from the server.
// It resumes the specified player's process execution (SIGCONT on Unix).
// Returns OK on success, ERROR if the player doesn't exist or resume fails,
// or OK if the player is already not running. Platform-dependent behavior.
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
