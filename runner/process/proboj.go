package process

import (
	"bufio"
	"fmt"
	"github.com/trojsten/ksp-proboj/runner/log"
	"io"
	"strings"
	"sync"
)

type ProbojProcess struct {
	*Process
	stdoutReader *bufio.Reader
	stderrReader *bufio.Reader
	log          log.Log
	logMutex     *sync.Mutex
	wait         *sync.WaitGroup
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

	pp.stdoutReader = bufio.NewReader(pp.Process.Stdout)
	pp.logMutex = &sync.Mutex{}

	if logConfig.Enabled {
		pp.stderrReader = bufio.NewReader(pp.Process.Stderr)
		pp.log = logConfig.Log
		go pp.stderrLoop()
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

func (pp *ProbojProcess) AsyncWrite(data string) <-chan error {
	ch := make(chan error)
	go func() {
		ch <- pp.Write(data)
	}()
	return ch
}

func (pp *ProbojProcess) readLine() (string, error) {
	return readln(pp.stdoutReader)
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
	pp.wait.Add(1)
	defer pp.wait.Done()

	for {
		data, err := readln(pp.stderrReader)
		if err != nil {
			if err != io.EOF {
				_ = pp.WriteLog(fmt.Sprintf("[proboj] error while reading stderr: %s\n", err.Error()))
			}
			break
		}
		_ = pp.WriteLog(fmt.Sprintf("%s\n", data))
	}

	pp.closeLogOnExit()
}

func (pp *ProbojProcess) closeLogOnExit() {
	<-pp.OnExit()

	defer pp.logMutex.Unlock()
	pp.logMutex.Lock()
	_, _ = pp.log.Write([]byte(fmt.Sprintf("[proboj] process terminated\n exit: %d\n err: %v\n", pp.Exit, pp.Error)))
	_ = pp.log.Close()
}

func (pp *ProbojProcess) WaitForEnd() {
	pp.wait.Wait()
}
