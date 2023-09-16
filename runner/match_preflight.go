package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	log2 "github.com/trojsten/ksp-proboj/runner/log"
	"github.com/trojsten/ksp-proboj/runner/process"
	"os"
	"path"
	"strings"
)

func (m *Match) preflight() error {
	if m.Started {
		return fmt.Errorf("the match was already Started")
	}
	signalMatchStart(m)
	m.Started = true
	m.Log = log.With()
	m.Log.SetPrefix(m.Game.Gamefolder)

	// Create folders
	err := os.MkdirAll(m.Directory(), 0o755)
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", m.Directory(), err)
	}
	err = os.MkdirAll(path.Join(m.Directory(), "logs"), 0o755)
	if err != nil {
		return fmt.Errorf("mkdir %s/logs: %w", m.Directory(), err)
	}

	err = m.startServer()
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	err = m.sendConfigToServer()
	if err != nil {
		return fmt.Errorf("send config to server: %w", err)
	}

	m.startPlayers()

	err = m.openObserver()
	if err != nil {
		return fmt.Errorf("open observer: %w", err)
	}

	return nil
}

func (m *Match) sendConfigToServer() error {
	return m.Server.Write(fmt.Sprintf("CONFIG\n%s\n%s\n.\n", strings.Join(m.Game.Players, " "), m.Game.Arguments))
}

func (m *Match) startServer() (err error) {
	m.Log.Debug("Creating server process", "server", m.Config.Server)

	logConfig, err := m.logConfig("__server")
	if err != nil {
		return
	}

	m.Server, err = process.NewProbojProcess(m.Config.Server, m.Directory(), logConfig)
	if err != nil {
		return
	}

	m.Log.Info("Starting server process")
	m.Server.Start()
	return
}

func (m *Match) startPlayers() {
	for _, player := range m.Game.Players {
		err := m.startPlayer(player)
		if err != nil {
			m.Log.Error("Failed to start player", "player", player, "err", err)
		}
	}
}

func (m *Match) logConfig(name string) (process.LogConfig, error) {
	var logConfig process.LogConfig
	if m.Config.DisableLogs {
		logConfig.Enabled = false
	} else {
		logConfig.Enabled = true
		suffix := "gz"
		if m.Config.DisableGzip {
			suffix = "txt"
		}

		fileName := path.Join(m.Directory(), "logs", fmt.Sprintf("%s.%s", name, suffix))
		file, err := os.Create(fileName)
		if err != nil {
			return process.LogConfig{}, err
		}

		if m.Config.DisableGzip {
			logConfig.Log = log2.NewPlainLog(file)
		} else {
			logConfig.Log, err = log2.NewGzipLog(file)
			if err != nil {
				return process.LogConfig{}, err
			}
		}
	}
	return logConfig, nil
}

func (m *Match) startPlayer(name string) error {
	program, exists := m.Config.Players[name]
	if !exists {
		return fmt.Errorf("player %s not found in config", name)
	}

	logConfig, err := m.logConfig(name)
	if err != nil {
		return err
	}

	m.Log.Debug("Creating player process", "player", name, "program", program)
	proc, err := process.NewProbojProcess(program, m.Directory(), logConfig)
	if err != nil {
		return err
	}

	m.Log.Info("Starting player process", "player", name)
	m.Players[name] = &proc
	proc.Start()
	return nil
}

func (m *Match) openObserver() (err error) {
	fileName := path.Join(m.Directory(), "observer.gz")
	m.Log.Debug("Opening observer file", "file", fileName)
	m.Observer, err = NewObserver(fileName)
	return
}
