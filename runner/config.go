package main

import (
	"github.com/charmbracelet/log"
	"github.com/trojsten/ksp-proboj/runner/process"
)

type Config struct {
	Server              string            `json:"server"`
	Players             map[string]string `json:"players"`
	Timeout             float64           `json:"timeout"`
	DisableLogs         bool              `json:"disable_logs"`
	DisableGzip         bool              `json:"disable_gzip"`
	ServerWorkDirectory string            `json:"server_workdir"`
}

type Game struct {
	Gamefolder string   `json:"gamefolder"`
	Players    []string `json:"players"`
	Arguments  string   `json:"args"`
}

type Match struct {
	Game    Game
	Config  Config
	Server  process.ProbojProcess
	Players map[string]process.ProbojProcess
	logger  log.Logger
	started bool
	ended   bool
}
