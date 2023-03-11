package main

import (
	"fmt"
	"github.com/trojsten/ksp-proboj/common"
	"strings"
)

type handlerFunc func(m *Match, args []string, payload string) common.RunnerResponse

var Handlers = map[string]handlerFunc{
	"TO PLAYER": cmdToPlayer,
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

func cmdToPlayer(m *Match, args []string, payload string) common.RunnerResponse {
	player := args[0]
	var note string
	if len(args) > 1 {
		note = strings.Join(args[1:], " ")
	}

	fmt.Println(note)

	proc, ok := m.Players[player]
	if !ok {
		m.logger.Error("Unknown player", "player", player)
		return common.RunnerResponse{Status: common.Error}
	}

	if !proc.IsRunning() {
		return common.RunnerResponse{Status: common.Died}
	}

	m.logger.Debug("Sending data to player", "player", player)
	err := proc.Write(payload)
	if err != nil {
		m.logger.Error("Failed writing data to player", "player", player, "err", err)
		return common.RunnerResponse{Status: common.Error}
	}
	return common.RunnerResponse{Status: common.Ok}
}
