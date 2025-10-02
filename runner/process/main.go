// Package process provides subprocess management capabilities for the Proboj runner.
// It handles process creation, lifecycle management (start/kill/pause/resume),
// I/O redirection, and cross-platform process control. The package includes
// both basic Process management and ProbojProcess for protocol-specific communication.
package process

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/shlex"
)

// Options contains configuration for creating a new Process.
type Options struct {
	// Command is the executable command to run
	Command string
	// Dir is the working directory for the process
	Dir string
	// Stdin enables stdin pipe when true
	Stdin bool
	// Stdout enables stdout pipe when true
	Stdout bool
	// Stderr enables stderr pipe when true
	Stderr bool
}

// Process represents a managed subprocess with stdin/stdout/stderr pipes,
// lifecycle management (start/kill/pause/resume), and exit status tracking.
type Process struct {
	// cmd is the underlying exec.Cmd instance
	cmd *exec.Cmd
	// Stdin is the writeable pipe for process input
	Stdin io.WriteCloser
	// Stdout is the readable pipe for process output
	Stdout io.ReadCloser
	// Stderr is the readable pipe for process error output
	Stderr io.ReadCloser

	// pid is the process ID
	pid int
	// started indicates whether the process has been started
	started bool
	// ended indicates whether the process has terminated
	ended bool
	// paused indicates whether the process is currently paused
	paused bool
	// exitChan is closed when the process exits
	exitChan chan struct{}
	// Exit contains the process exit code
	Exit int
	// Error contains any error that occurred during process execution
	Error error
}

// NewProcess creates a new Process instance with the given options.
// It parses the command, resolves absolute paths, sets up the working
// directory, and creates pipes for enabled I/O streams. Returns the
// Process instance or an error if setup fails.
func NewProcess(options Options) (p Process, err error) {
	parts, err := shlex.Split(options.Command)
	if err != nil {
		return
	}

	if !filepath.IsAbs(parts[0]) {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return
		}
		parts[0] = filepath.Join(wd, parts[0])
	}

	p.cmd = exec.Command(parts[0], parts[1:]...)
	p.cmd.Dir = options.Dir
	setProcessGroupID(p.cmd)

	if options.Stdin {
		p.Stdin, err = p.cmd.StdinPipe()
		if err != nil {
			return
		}
	}
	if options.Stdout {
		p.Stdout, err = p.cmd.StdoutPipe()
		if err != nil {
			return
		}
	}
	if options.Stderr {
		p.Stderr, err = p.cmd.StderrPipe()
		if err != nil {
			return
		}
	}

	p.exitChan = make(chan struct{})
	return
}

func (p *Process) run() error {
	err := p.cmd.Start()
	if err != nil {
		return err
	}

	p.pid = p.cmd.Process.Pid

	err = p.cmd.Wait()
	if exiterr, ok := err.(*exec.ExitError); ok {
		p.Exit = exiterr.ExitCode()
	} else {
		return err
	}
	return nil
}

// Start begins process execution in a goroutine.
// Returns a channel that will be closed when the process exits.
func (p *Process) Start() chan struct{} {
	if p.started {
		p.Error = fmt.Errorf("process was already started")
		return p.exitChan
	}
	p.started = true

	go func() {
		err := p.run()
		if err != nil {
			p.Error = err
		}
		p.ended = true
		close(p.exitChan)
	}()
	return p.exitChan
}

// OnExit returns a channel that is closed when the process exits.
func (p *Process) OnExit() chan struct{} {
	return p.exitChan
}

// IsRunning returns true if the process has been started and has not yet ended.
func (p *Process) IsRunning() bool {
	return p.started && !p.ended
}

// Kill terminates the process forcefully.
func (p *Process) Kill() error {
	if !p.IsRunning() || p.pid == 0 {
		return fmt.Errorf("process is not running")
	}

	// Killing paused process tends to have unknown consequences.
	if p.IsPaused() {
		err := p.Resume()
		if err != nil {
			return err
		}
	}

	return terminateProcess(p.pid)
}

// Pause suspends process execution (platform-dependent).
// Uses SIGSTOP on Unix systems.
func (p *Process) Pause() error {
	if !p.IsRunning() || p.pid == 0 {
		return fmt.Errorf("process is not running")
	}

	if p.IsPaused() {
		return fmt.Errorf("process is already paused")
	}

	p.paused = true
	return pauseProcess(p.pid)
}

// Resume resumes process execution (platform-dependent).
// Uses SIGCONT on Unix systems.
func (p *Process) Resume() error {
	if !p.IsRunning() || p.pid == 0 {
		return fmt.Errorf("process is not running")
	}

	if !p.IsPaused() {
		return fmt.Errorf("process is not paused")
	}

	p.paused = false
	return resumeProcess(p.pid)
}

// IsPaused returns true if the process is currently paused.
func (p *Process) IsPaused() bool {
	return p.paused
}

// readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix       = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
