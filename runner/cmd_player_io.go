package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/trojsten/ksp-proboj/libproboj"
)

// cmdToPlayer handles the "TO PLAYER" command from the server.
// It sends the payload to the specified player's stdin with optional
// logging comment. Returns OK on success, ERROR on failure, or DIED
// if the player process is not running. Implements write timeout handling
// and process cleanup on timeout.
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

// cmdReadPlayer handles the "READ PLAYER" command from the server.
// It reads data from the specified player's stdout until a "." line
// is encountered. Supports optional timeout multiplier argument.
// Returns OK with player data on success, ERROR on failure, or DIED
// if the player process is not running. Implements timeout handling
// and automatic process cleanup on timeout.
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

	// Parse optional timeout multiplier
	var multiplier float64 = 1.0 // Default multiplier (no change)
	if len(args) > 1 {
		multiplier, err := strconv.ParseFloat(args[1], 64)
		if err != nil || multiplier <= 0 {
			m.Log.Error("Invalid timeout multiplier", "player", player, "multiplier", args[1])
			return libproboj.RunnerResponse{Status: libproboj.Error}
		}
	}

	actualTimeout := timeout * multiplier
	m.Log.Debug("Reading data from player", "player", player, "timeout", actualTimeout, "multiplier", multiplier)
	select {
	case <-time.After(time.Millisecond * time.Duration(actualTimeout*1000)):
		m.Log.Warn("Player timeouted", "player", player, "timeout", actualTimeout)
		_ = proc.WriteLog(fmt.Sprintf("[proboj] killing process due to read timeout (adjusted by %.1fx)\n", multiplier))
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
