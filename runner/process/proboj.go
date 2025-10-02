package process

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/trojsten/ksp-proboj/runner/log"
)

// ProbojProcess extends Process with Proboj protocol-specific functionality.
// It handles protocol-compliant I/O (messages terminated by "."), logging,
// and asynchronous operations for communication with game servers and players.
type ProbojProcess struct {
	*Process
	// stdoutReader provides buffered reading from process stdout
	stdoutReader *bufio.Reader
	// stderrReader provides buffered reading from process stderr
	stderrReader *bufio.Reader
	// log handles logging for this process
	log log.Log
	// logMutex synchronizes access to the log
	logMutex *sync.Mutex
	// wait tracks background goroutines for proper cleanup
	wait *sync.WaitGroup
}

// NewProbojProcess creates a new ProbojProcess with protocol-specific setup.
// It creates the underlying Process, sets up buffered readers for I/O,
// configures logging, and starts stderr monitoring goroutine if enabled.
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
	pp.wait = &sync.WaitGroup{}

	if logConfig.Enabled {
		pp.stderrReader = bufio.NewReader(pp.Process.Stderr)
		pp.log = logConfig.Log
		go pp.stderrLoop()
	} else {
		pp.log = log.NewNullLog()
	}
	return
}

// Write sends data to the process stdin.
func (pp *ProbojProcess) Write(data string) error {
	if !pp.IsRunning() {
		return fmt.Errorf("process is not running")
	}

	_, err := pp.Process.Stdin.Write([]byte(data))
	return err
}

// AsyncWrite sends data to process stdin asynchronously.
func (pp *ProbojProcess) AsyncWrite(data string) <-chan error {
	ch := make(chan error)
	go func() {
		ch <- pp.Write(data)
	}()
	return ch
}

// readLine reads a single line from process stdout.
func (pp *ProbojProcess) readLine() (string, error) {
	return readln(pp.stdoutReader)
}

// Read reads protocol-compliant data from process stdout.
// Reads lines until encountering a "." line, then returns
// the concatenated content.
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

// ReadResult represents the result of an asynchronous read operation.
// It contains either the successfully read data or an error.
type ReadResult struct {
	Data  string
	Error error
}

// AsyncRead reads protocol-compliant data from process stdout asynchronously.
// Returns a channel that will receive a ReadResult.
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

// WriteLog writes data to the process log with thread-safe access.
func (pp *ProbojProcess) WriteLog(data string) error {
	defer pp.logMutex.Unlock()
	pp.logMutex.Lock()
	_, err := pp.log.Write([]byte(data))
	return err
}

// stderrLoop continuously reads from process stderr and logs the output.
// Runs in a background goroutine until EOF or error is encountered.
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

// closeLogOnExit waits for process exit and writes termination information to the log.
// Records the exit code and error, then closes the log file.
func (pp *ProbojProcess) closeLogOnExit() {
	<-pp.OnExit()

	defer pp.logMutex.Unlock()
	pp.logMutex.Lock()
	_, _ = pp.log.Write([]byte(fmt.Sprintf("[proboj] process terminated\n exit: %d\n err: %v\n", pp.Exit, pp.Error)))
	_ = pp.log.Close()
}

// WaitForEnd blocks until all background goroutines for this process have completed.
// Ensures proper cleanup before process destruction.
func (pp *ProbojProcess) WaitForEnd() {
	pp.wait.Wait()
}
