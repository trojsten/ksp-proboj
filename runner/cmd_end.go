package main

import "github.com/trojsten/ksp-proboj/common"

func cmdEnd(m *Match, _ []string, _ string) common.RunnerResponse {
	m.logger.Info("Server ended the game.")
	err := m.Server.Kill()
	if err != nil {
		m.logger.Error("Failed to kill server", "err", err)
	}
	m.ended = true
	return common.RunnerResponse{Status: common.Ignore}
}
