package main

import (
	"fmt"
	"strings"
)

type handlerFunc func(args []string, payload string) error

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
		err := handler(args, payload)
		if err != nil {
			m.logger.Error("Command handler failed", "err", err)
		}
		return
	}

	m.logger.Warn("Server sent unknown command", "cmd", cmd)
}

func cmdToPlayer(args []string, payload string) error {
	fmt.Println(payload)
	return nil
}
