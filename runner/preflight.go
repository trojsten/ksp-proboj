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

	// TODO: start players

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
