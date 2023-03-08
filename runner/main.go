package main

import (
	"github.com/charmbracelet/log"
	"time"
)

func main() {
	log.SetTimeFormat(time.StampMilli)
	log.SetLevel(log.DebugLevel)

	m := Match{
		Game: Game{
			Gamefolder: "test",
			Players:    []string{"a", "b"},
			Arguments:  "",
		},
		Config: Config{
			//Server: "/usr/bin/false",
			Server: "/home/ano95/Dev/Trojsten/ksp-proboj-2023-jar/srv",
			Players: map[string]string{
				"a": "", "b": "",
			},
			Timeout:             1,
			DisableLogs:         false,
			DisableGzip:         false,
			ServerWorkDirectory: ".",
		},
	}

	m.Run()
}
