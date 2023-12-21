package main

import (
	"fmt"
	"github.com/trojsten/ksp-proboj/libproboj"
	"strings"
	"time"
)

func cmdToPlayer(m *Match, args []string, payload string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.Log.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]
	var note string
	if len(args) > 1 {
		note = strings.Join(args[1:], " ")
	}

	proc, ok := m.Players[player]
	if !ok {
		m.Log.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if !proc.IsRunning() {
		m.Log.Warn("Player is not running", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	}

	err := proc.WriteLog(fmt.Sprintf("[proboj] %s\n", note))
	if err != nil {
		m.Log.Error("Failed writing data to players' log", "player", player, "err", err)
	}

	m.Log.Debug("Sending data to player", "player", player)
	select {
	case <-time.After(5 * time.Second):
		m.Log.Error("Write timeouted", "player", player, "err", err)
		_ = proc.WriteLog(fmt.Sprintf("[proboj] killing process due to write timeout\n"))
		err := proc.Kill()
		if err != nil {
			m.Log.Error("Failed to kill player", "player", player, "err", err)
		}
		return libproboj.RunnerResponse{Status: libproboj.Error}
	case <-proc.OnExit():
		m.Log.Warn("Player died", "player", player, "exit", proc.Exit, "err", proc.Error)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	case err := <-proc.AsyncWrite(fmt.Sprintf("%s\n.\n", payload)):
		if err != nil {
			m.Log.Error("Failed writing data to player", "player", player, "err", err)
			return libproboj.RunnerResponse{Status: libproboj.Error}
		}
		return libproboj.RunnerResponse{Status: libproboj.Ok}
	}
}

func cmdReadPlayer(m *Match, args []string, _ string) libproboj.RunnerResponse {
	if len(args) < 1 {
		m.Log.Error("Invalid command syntax: missing arguments")
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}
	player := args[0]

	proc, ok := m.Players[player]
	if !ok {
		m.Log.Error("Unknown player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Error}
	}

	if !proc.IsRunning() {
		m.Log.Debug("Ignoring read from dead player", "player", player)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	}

	playerConf := m.Config.Players[player]
	timeout := m.Config.Timeout[playerConf.Language]

	m.Log.Debug("Reading data from player", "player", player)
	select {
	case <-time.After(time.Millisecond * time.Duration(timeout*1000)):
		m.Log.Warn("Player timeouted", "player", player)
		_ = proc.WriteLog(fmt.Sprintf("[proboj] killing process due to read timeout\n"))
		err := proc.Kill()
		if err != nil {
			m.Log.Error("Failed to kill player", "player", player, "err", err)
		}
		return libproboj.RunnerResponse{Status: libproboj.Died}
	case <-proc.OnExit():
		m.Log.Warn("Player died", "player", player, "exit", proc.Exit, "err", proc.Error)
		return libproboj.RunnerResponse{Status: libproboj.Died}
	case result := <-proc.AsyncRead():
		if result.Error != nil {
			m.Log.Error("Error while reading from player", "player", player, "err", result.Error)
			return libproboj.RunnerResponse{Status: libproboj.Error}
		}
		return libproboj.RunnerResponse{
			Status:  libproboj.Ok,
			Payload: result.Data,
		}
	}
}
