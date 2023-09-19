package main

import "github.com/trojsten/ksp-proboj/libproboj"

func cmdEnd(m *Match, _ []string, _ string) libproboj.RunnerResponse {
	m.Log.Info("Server ended the game.")
	err := m.Server.Kill()
	if err != nil {
		m.Log.Error("Failed to kill server", "err", err)
	}
	m.Ended = true
	return libproboj.RunnerResponse{Status: libproboj.Ignore}
}
