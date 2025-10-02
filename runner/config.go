package main

import (
	"github.com/charmbracelet/log"
	"github.com/trojsten/ksp-proboj/runner/process"
	"path"
)

// PlayerConf contains configuration for a single player/bot.
type PlayerConf struct {
	// Command is the executable command to run for this player
	Command string `json:"command"`
	// Language identifies the programming language for timeout configuration
	Language string `json:"language"`
}

// Config contains global runner configuration for all games.
type Config struct {
	// Server is the path to the game server executable
	Server string `json:"server"`
	// Players maps player names to their configurations
	Players map[string]PlayerConf `json:"players"`
	// ProcessesPerPlayer specifies how many processes to spawn per player (default 1)
	// If > 1, players are named player_0, player_1, etc.
	ProcessesPerPlayer int `json:"processes_per_player"`
	// Timeout maps programming languages to their maximum response time in seconds
	Timeout map[string]float64 `json:"timeout"`
	// DisableLogs when true prevents any log file creation
	DisableLogs bool `json:"disable_logs"`
	// DisableGzip when true uses plain text logs instead of gzip compression
	DisableGzip bool `json:"disable_gzip"`
	// GameRoot is the root directory where game data will be stored
	GameRoot string `json:"game_root"`
}

// Game defines a single game instance to be run by the runner.
type Game struct {
	// Gamefolder is the name of the directory created under GameRoot for this game
	Gamefolder string `json:"gamefolder"`
	// Players is the list of player names participating in this game
	Players []string `json:"players"`
	// Arguments contains additional configuration data passed to the server
	Arguments string `json:"args"`
}

// Match represents a single game instance with its configuration,
// server process, player processes, and logging infrastructure.
type Match struct {
	// Game contains the specific game configuration
	Game Game
	// Config contains the global runner configuration
	Config Config
	// Server is the game server process that controls the game
	Server process.ProbojProcess
	// Players maps player names to their process instances
	Players map[string]*process.ProbojProcess
	// Log is the logger for this match instance
	Log log.Logger
	// Started indicates whether the match has begun
	Started bool
	// Ended indicates whether the match has finished
	Ended bool

	// Observer handles recording game events for replay
	Observer Observer
}

func (m *Match) Directory() string {
	return path.Join(m.Config.GameRoot, m.Game.Gamefolder)
}

func NewMatch(config Config, game Game) *Match {
	return &Match{Game: game, Config: config, Players: map[string]*process.ProbojProcess{}}
}
