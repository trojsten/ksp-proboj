package main

import (
	"github.com/trojsten/ksp-proboj/libproboj"
	"strings"
)

type handlerFunc func(m *Match, args []string, payload string) libproboj.RunnerResponse

var Handlers = map[string]handlerFunc{
	"TO PLAYER":     cmdToPlayer,
	"READ PLAYER":   cmdReadPlayer,
	"TO OBSERVER":   cmdToObserver,
	"KILL PLAYER":   cmdKillPlayer,
	"PAUSE PLAYER":  cmdPausePlayer,
	"RESUME PLAYER": cmdResumePlayer,
	"SCORES":        cmdScores,
	"END":           cmdEnd,
}

func (m *Match) parseCommand(data string) {
	cmd, payload, _ := strings.Cut(data, "\n")
	m.Log.Debug("Parsing command", "cmd", cmd)

	for prefix, handler := range Handlers {
		if !strings.HasPrefix(cmd, prefix) {
			continue
		}

		args := strings.Split(strings.TrimSpace(strings.TrimPrefix(cmd, prefix)), " ")
		m.Log.Debug("Using command handler", "handler", prefix, "args", args)
		response := handler(m, args, payload)
		if response.Status == libproboj.Ignore {
			return
		}

		err := m.Server.Write(response.String())
		if err != nil {
			m.Log.Error("Failed writing response back to the server", "err", err)
		}
		return
	}

	m.Log.Warn("Server sent unknown command", "cmd", cmd)
}
