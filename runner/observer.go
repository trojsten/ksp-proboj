package main

import (
	"github.com/trojsten/ksp-proboj/runner/log"
	"os"
)

type Observer struct {
	log *log.GzipLog
}

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

func (o Observer) Observe(data string) error {
	if o.log == nil {
		return nil
	}

	_, err := o.log.Write([]byte(data))
	return err
}

func (o Observer) Close() error {
	if o.log == nil {
		return nil
	}

	err := o.log.Close()
	o.log = nil
	return err
}
