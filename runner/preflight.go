package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/trojsten/ksp-proboj/runner/process"
	"strings"
)

func (m *Match) preflight() (err error) {
	if m.started {
		return fmt.Errorf("the match was already started")
	}
	m.started = true
	m.logger = log.With()
	m.logger.SetPrefix(m.Game.Gamefolder)

	err = m.startServer()
	if err != nil {
		return
	}

	m.Players = map[string]process.ProbojProcess{}
	m.startPlayers()

	err = m.sendConfigToServer()
	if err != nil {
		return
	}

	return
}

func (m *Match) sendConfigToServer() error {
	return m.Server.Write(fmt.Sprintf("CONFIG\n%s\n%s\n.\n", strings.Join(m.Game.Players, " "), m.Game.Arguments))
}

func (m *Match) startServer() (err error) {
	m.logger.Debug("Creating server process", "server", m.Config.Server)
	m.Server, err = process.NewProbojProcess(m.Config.Server, m.Config.ServerWorkDirectory)
	if err != nil {
		return
	}

	m.logger.Info("Starting server process")
	m.Server.Start()
	return
}

func (m *Match) startPlayers() {
	for _, player := range m.Game.Players {
		err := m.startPlayer(player)
		if err != nil {
			m.logger.Error("Failed to start player", "player", player, "err", err)
		}
	}
}

func (m *Match) startPlayer(name string) error {
	program, exists := m.Config.Players[name]
	if !exists {
		return fmt.Errorf("player %s not found in config", name)
	}
	m.logger.Debug("Creating player process", "player", name, "program", program)
	proc, err := process.NewProbojProcess(program, m.Config.ServerWorkDirectory)
	if err != nil {
		return err
	}

	m.logger.Info("Starting player process", "player", name)
	m.Players[name] = proc
	proc.Start()
	return nil
}
