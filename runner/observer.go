package main

import (
	"os"

	"github.com/trojsten/ksp-proboj/runner/log"
)

// Observer handles game recording by writing game events to a compressed log file.
type Observer struct {
	log *log.GzipLog
}

// NewObserver creates a new Observer that writes to the specified file path.
func NewObserver(path string) (Observer, error) {
	file, err := os.Create(path)
	if err != nil {
		return Observer{}, err
	}

	gzipLog, err := log.NewGzipLog(file)
	if err != nil {
		return Observer{}, err
	}

	return Observer{
		log: gzipLog,
	}, nil
}

// Observe writes game data to the observer log file.
func (o Observer) Observe(data string) error {
	if o.log == nil {
		return nil
	}

	_, err := o.log.Write([]byte(data))
	return err
}

// Close finalizes the observer log and releases resources.
// Flushes any pending data and closes the underlying file.
func (o Observer) Close() error {
	if o.log == nil {
		return nil
	}

	err := o.log.Close()
	o.log = nil
	return err
}
