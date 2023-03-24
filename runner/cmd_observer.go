package main

import "github.com/trojsten/ksp-proboj/common"

func cmdToObserver(m *Match, _ []string, payload string) common.RunnerResponse {
	_, err := m.observer.Write([]byte(payload + "\n"))
	if err != nil {
		return common.RunnerResponse{Status: common.Error}
	}
	return common.RunnerResponse{Status: common.Ok}
}
