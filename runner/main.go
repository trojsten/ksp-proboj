package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] config games\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage
	debug := false
	flag.BoolVar(&debug, "v", debug, "print verbose logs")
	concurrency := 1
	flag.IntVar(&concurrency, "c", concurrency, "number of games to run concurrently")

	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	log.SetTimeFormat(time.StampMilli)
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	registerSignals()

	configFilename := flag.Arg(0)
	gamesFilename := flag.Arg(1)

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

	<-finish
}
