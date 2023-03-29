package process

import (
	"bufio"
	"fmt"
	"github.com/trojsten/ksp-proboj/runner/log"
	"strings"
	"sync"
)

type ProbojProcess struct {
	*Process
	stdoutScanner *bufio.Scanner
	stderrScanner *bufio.Scanner
	log           log.Log
	logMutex      *sync.Mutex
}

func NewProbojProcess(command string, dir string, logConfig LogConfig) (pp ProbojProcess, err error) {
	proc, err := NewProcess(Options{
		Command: command,
		Dir:     dir,
		Stdin:   true,
		Stdout:  true,
		Stderr:  logConfig.Enabled,
	})
	if err != nil {
		return
	}
	pp.Process = &proc

	pp.stdoutScanner = bufio.NewScanner(pp.Process.Stdout)

	if logConfig.Enabled {
		pp.stderrScanner = bufio.NewScanner(pp.Process.Stderr)
		pp.logMutex = &sync.Mutex{}
		pp.log = logConfig.Log
		go pp.stderrLoop()
		go pp.closeLogOnExit()
	} else {
		pp.log = log.NewNullLog()
	}
	return
}

func (pp *ProbojProcess) Write(data string) error {
	if !pp.IsRunning() {
		return fmt.Errorf("process is not running")
	}
	_, err := pp.Process.Stdin.Write([]byte(data))
	return err
}

func (pp *ProbojProcess) readLine() (string, error) {
	if !pp.stdoutScanner.Scan() && pp.stdoutScanner.Err() != nil {
		return "", pp.stdoutScanner.Err()
	}
	return pp.stdoutScanner.Text(), nil
}

func (pp *ProbojProcess) Read() (string, error) {
	result := []string{}
	for true {
		input, err := pp.readLine()
		if err != nil {
			return "", err
		}
		if input == "." {
			break
		}
		result = append(result, input)
	}
	return strings.Join(result, "\n"), nil
}

type ReadResult struct {
	Data  string
	Error error
}

func (pp *ProbojProcess) AsyncRead() <-chan ReadResult {
	ch := make(chan ReadResult)
	go func() {
		data, err := pp.Read()
		ch <- ReadResult{
			Data:  data,
			Error: err,
		}
	}()
	return ch
}

func (pp *ProbojProcess) WriteLog(data string) error {
	defer pp.logMutex.Unlock()
	pp.logMutex.Lock()
	_, err := pp.log.Write([]byte(data))
	return err
}

func (pp *ProbojProcess) stderrLoop() {
	for pp.stderrScanner.Scan() {
		// Write error ignored here.
		_ = pp.WriteLog(fmt.Sprintf("%s\n", pp.stderrScanner.Text()))
	}
	if err := pp.stderrScanner.Err(); err != nil {
		// Write error ignored here.
		_ = pp.WriteLog(fmt.Sprintf("[proboj] error while scanning stderr: %s\n", err.Error()))
	}
}

func (pp *ProbojProcess) closeLogOnExit() {
	<-pp.OnExit()

	defer pp.logMutex.Unlock()
	pp.logMutex.Lock()
	_, _ = pp.log.Write([]byte(fmt.Sprintf("[proboj] process terminated\n exit: %d\n err: %v\n", pp.Exit, pp.Error)))
	_ = pp.log.Close()
}
