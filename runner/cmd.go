package main

import (
	"github.com/trojsten/ksp-proboj/common"
	"strings"
)

type handlerFunc func(m *Match, args []string, payload string) common.RunnerResponse

var Handlers = map[string]handlerFunc{
	"TO PLAYER":   cmdToPlayer,
	"READ PLAYER": cmdReadPlayer,
	"TO OBSERVER": cmdToObserver,
	"KILL PLAYER": cmdKillPlayer,
}

func (m *Match) parseCommand(data string) {
	cmd, payload, _ := strings.Cut(data, "\n")
	m.logger.Debug("Parsing command", "cmd", cmd)

	for prefix, handler := range Handlers {
		if !strings.HasPrefix(cmd, prefix) {
			continue
		}

		args := strings.Split(strings.TrimSpace(strings.TrimPrefix(cmd, prefix)), " ")
		m.logger.Debug("Using command handler", "handler", prefix, "args", args)
		response := handler(m, args, payload)
		err := m.Server.Write(response.String())
		if err != nil {
			m.logger.Error("Failed writing response back to the server", "err", err)
		}
		return
	}

	m.logger.Warn("Server sent unknown command", "cmd", cmd)
}
