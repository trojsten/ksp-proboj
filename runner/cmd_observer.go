package main

import "github.com/trojsten/ksp-proboj/libproboj"

func cmdToObserver(m *Match, _ []string, payload string) libproboj.RunnerResponse {
	err := m.Observer.Observe(payload)
	if err != nil {
		m.Log.Error("Could not write data to observer.", "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	return libproboj.RunnerResponse{Status: libproboj.Ok}
}
