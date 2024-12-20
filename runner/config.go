package main

import (
	"github.com/charmbracelet/log"
	"github.com/trojsten/ksp-proboj/runner/process"
	"path"
)

type PlayerConf struct {
	Command  string `json:"command"`
	Language string `json:"language"`
}

type Config struct {
	Server             string                `json:"server"`
	Players            map[string]PlayerConf `json:"players"`
	ProcessesPerPlayer int                   `json:"processes_per_player"`
	Timeout            map[string]float64    `json:"timeout"`
	DisableLogs        bool                  `json:"disable_logs"`
	DisableGzip        bool                  `json:"disable_gzip"`
	GameRoot           string                `json:"game_root"`
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
	Players map[string]*process.ProbojProcess
	Log     log.Logger
	Started bool
	Ended   bool

	Observer Observer
}

func (m *Match) Directory() string {
	return path.Join(m.Config.GameRoot, m.Game.Gamefolder)
}

func NewMatch(config Config, game Game) *Match {
	return &Match{Game: game, Config: config, Players: map[string]*process.ProbojProcess{}}
}
