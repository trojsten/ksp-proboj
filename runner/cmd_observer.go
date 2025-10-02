package main

import "github.com/trojsten/ksp-proboj/libproboj"

// cmdToObserver handles the "TO OBSERVER" command from the server.
// It writes the payload data to the observer log for game recording.
// Returns OK on successful write, ERROR if observer write fails.
func cmdToObserver(m *Match, _ []string, payload string) libproboj.RunnerResponse {
	err := m.Observer.Observe(payload)
	if err != nil {
		m.Log.Error("Could not write data to observer.", "err", err)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	return libproboj.RunnerResponse{Status: libproboj.Ok}
}
