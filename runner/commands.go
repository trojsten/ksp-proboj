package main

import (
	"fmt"
	"github.com/trojsten/ksp-proboj/common"
	"strings"
	"time"
)

type handlerFunc func(m *Match, args []string, payload string) common.RunnerResponse

var Handlers = map[string]handlerFunc{
	"TO PLAYER":   cmdToPlayer,
	"READ PLAYER": cmdReadPlayer,
	"TO OBSERVER": cmdToObserver,
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
	if len(args) < 1 {
		m.logger.Error("Invalid command syntax: missing arguments")
		return common.RunnerResponse{Status: common.Error}
	}
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

func cmdReadPlayer(m *Match, args []string, _ string) common.RunnerResponse {
	if len(args) < 1 {
		m.logger.Error("Invalid command syntax: missing arguments")
		return common.RunnerResponse{Status: common.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.logger.Error("Unknown player", "player", player)
		return common.RunnerResponse{Status: common.Error}
	}

	if !proc.IsRunning() {
		m.logger.Debug("Ignoring read from dead player", "player", player)
		return common.RunnerResponse{Status: common.Died}
	}

	m.logger.Debug("Reading data from player", "player", player)
	select {
	case <-time.After(time.Second * time.Duration(m.Config.Timeout)):
		m.logger.Warn("Player timeouted", "player", player)
		err := proc.Kill()
		if err != nil {
			m.logger.Error("Failed to kill player", "player", player, "err", err)
		}
		return common.RunnerResponse{Status: common.Died}
	case <-proc.OnExit():
		m.logger.Warn("Player died", "player", player, "exit", proc.Exit, "err", proc.Error)
		return common.RunnerResponse{Status: common.Died}
	case result := <-proc.AsyncRead():
		if result.Error != nil {
			m.logger.Error("Error while reading from player", "player", player, "err", result.Error)
			return common.RunnerResponse{Status: common.Error}
		}
		return common.RunnerResponse{
			Status:  common.Ok,
			Payload: result.Data,
		}
	}
}

func cmdToObserver(m *Match, _ []string, payload string) common.RunnerResponse {
	_, err := m.observer.Write([]byte(payload))
	if err != nil {
		return common.RunnerResponse{Status: common.Error}
	}
	return common.RunnerResponse{Status: common.Ok}
}
