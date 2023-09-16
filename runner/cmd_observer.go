package main

import "github.com/trojsten/ksp-proboj/libproboj"

func cmdToObserver(m *Match, _ []string, payload string) libproboj.RunnerResponse {
	_, err := m.observer.Write([]byte(payload + "\n"))
	if err != nil {
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	return libproboj.RunnerResponse{Status: libproboj.Ok}
}
