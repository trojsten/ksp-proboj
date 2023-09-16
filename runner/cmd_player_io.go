package main

import (
	"fmt"
	"github.com/trojsten/ksp-proboj/libproboj"
	"strings"
	"time"
)

func cmdToPlayer(m *Match, args []string, payload string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.logger.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]
	var note string
	if len(args) > 1 {
		note = strings.Join(args[1:], " ")
	}

	proc, ok := m.Players[player]
	if !ok {
		m.logger.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if !proc.IsRunning() {
		m.logger.Warn("Player is not running", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	}

	err := proc.WriteLog(fmt.Sprintf("[proboj] %s\n", note))
	if err != nil {
		m.logger.Error("Failed writing data to players' log", "player", player, "err", err)
	}

	m.logger.Debug("Sending data to player", "player", player)
	select {
	case <-time.After(time.Second):
		m.logger.Error("Write timeouted", "player", player, "err", err)
		err := proc.Kill()
		if err != nil {
			m.logger.Error("Failed to kill player", "player", player, "err", err)
		}
		return libproboj.RunnerResponse{Status: libproboj.Error}
	case <-proc.OnExit():
		m.logger.Warn("Player died", "player", player, "exit", proc.Exit, "err", proc.Error)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	case err := <-proc.AsyncWrite(fmt.Sprintf("%s\n.\n", payload)):
		if err != nil {
			m.logger.Error("Failed writing data to player", "player", player, "err", err)
			return libproboj.RunnerResponse{Status: libproboj.Error}
		}
		return libproboj.RunnerResponse{Status: libproboj.Ok}
	}
}

func cmdReadPlayer(m *Match, args []string, _ string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.logger.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.logger.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if !proc.IsRunning() {
		m.logger.Debug("Ignoring read from dead player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	}

	m.logger.Debug("Reading data from player", "player", player)
	select {
	case <-time.After(time.Millisecond * time.Duration(m.Config.Timeout*1000)):
		m.logger.Warn("Player timeouted", "player", player)
		err := proc.Kill()
		if err != nil {
			m.logger.Error("Failed to kill player", "player", player, "err", err)
		}
		return libproboj.RunnerResponse{Status: libproboj.Died}
	case <-proc.OnExit():
		m.logger.Warn("Player died", "player", player, "exit", proc.Exit, "err", proc.Error)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	case result := <-proc.AsyncRead():
		if result.Error != nil {
			m.logger.Error("Error while reading from player", "player", player, "err", result.Error)
			return libproboj.RunnerResponse{Status: libproboj.Error}
		}
		return libproboj.RunnerResponse{
			Status:  libproboj.Ok,
			Payload: result.Data,
		}
	}
}
