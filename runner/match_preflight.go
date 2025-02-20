package main

import (
	"fmt"
	"maps"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	log2 "github.com/trojsten/ksp-proboj/runner/log"
	"github.com/trojsten/ksp-proboj/runner/process"
	"github.com/trojsten/ksp-proboj/runner/websockets"
)

func (m *Match) preflight() error {
	if m.Started {
		return fmt.Errorf("the match was already started")
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

	var humanPlayers = make(map[string]PlayerConf)
	for _, player := range m.Game.Players {
		if m.Config.Players[player].Language == "human" {
			humanPlayers[player] = m.Config.Players[player]
		}
	}

	if len(humanPlayers) > 0 {
		if m.Config.http.port == 0 {
			m.Config.http.port = 8080
		}
		go func() {
			err := websockets.StartWebSocketServer(m.Config.http.port, m.Config.http.sourceRoot)
			if err != nil {
				//return fmt.Errorf("start websocket server: %w", err)
				return
			}
		}()

		m.Log.Info("human players are ready to connect at http://localhost:" + strconv.Itoa(m.Config.http.port))
		// wait for all connections to be established
		websockets.WaitForPlayers(maps.Keys(humanPlayers))

		m.Log.Info("All players connected, starting game")
	}

	err = m.preparePlayers()
	if err != nil {
		return fmt.Errorf("prepare players: %w", err)
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

func (m *Match) preparePlayers() error {
	if m.Config.ProcessesPerPlayer <= 1 {
		return nil
	}

	var newPlayers []string
	newPlayerConfig := make(map[string]PlayerConf)

	for _, player := range m.Game.Players {
		for i := 0; i < m.Config.ProcessesPerPlayer; i++ {
			playerName := fmt.Sprintf("%s_%d", player, i)
			newPlayers = append(newPlayers, playerName)

			if _, exists := m.Config.Players[player]; !exists {
				return fmt.Errorf("player %s not found in config", player)
			}
			newPlayerConfig[playerName] = m.Config.Players[player]
		}
	}

	m.Config.Players = newPlayerConfig
	m.Game.Players = newPlayers

	return nil
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
	player, exists := m.Config.Players[name]

	if !exists {
		return fmt.Errorf("player %s not found in config", name)
	}

	if player.Language == "human" {
		return nil
	}

	logConfig, err := m.logConfig(name)
	if err != nil {
		return err
	}

	_, exists = m.Config.Timeout[player.Language]
	if !exists {
		return fmt.Errorf("language %s not found in config", player.Language)
	}

	m.Log.Debug("Creating player process", "player", name, "command", player.Command)
	proc, err := process.NewProbojProcess(player.Command, m.Directory(), logConfig)
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
