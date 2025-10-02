// Package main implements the Proboj runner that manages game servers,
// player processes, and their communication according to the Proboj protocol.
// The runner handles process lifecycle management, command parsing, logging,
// and coordination between the game server and player bots.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Version: %v\n", VERSION)
	flag.PrintDefaults()
	fmt.Println()
}

func main() {
	flag.Usage = printUsage
	debug := false
	flag.BoolVar(&debug, "verbose", debug, "print verbose logs")
	concurrency := 1
	flag.IntVar(&concurrency, "concurrency", concurrency, "number of games to run concurrently")
	configFilename := "config.json"
	gamesFilename := "games.json"
	flag.StringVar(&configFilename, "config", configFilename, "configuration filename")
	flag.StringVar(&gamesFilename, "games", gamesFilename, "game definition filename")

	flag.Parse()
	flag.Usage()

	log.SetTimeFormat(time.StampMilli)
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	registerSignals()

	var config Config
	var games []Game

	configFile, err := os.ReadFile(configFilename)
	if err != nil {
		log.Error("Could not open config file", "file", configFilename, "err", err)
		os.Exit(1)
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Error("Could not parse config file", "err", err)
		os.Exit(1)
	}

	gamesFile, err := os.ReadFile(gamesFilename)
	if err != nil {
		log.Error("Could not open games file", "file", gamesFilename, "err", err)
		os.Exit(1)
	}

	err = json.Unmarshal(gamesFile, &games)
	if err != nil {
		log.Error("Could not parse games file", "err", err)
		os.Exit(1)
	}

	if concurrency > 1 {
		runParallel(config, games, concurrency)
	} else {
		runSequentially(config, games)
	}
}
