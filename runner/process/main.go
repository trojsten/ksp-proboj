package process

import (
	"fmt"
	"github.com/google/shlex"
	"io"
	"os"
	"os/exec"
	"path"
)

type Options struct {
	Command string
	Dir     string
	Stdin   bool
	Stdout  bool
	Stderr  bool
}

type Process struct {
	cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser

	started  bool
	ended    bool
	exitChan chan struct{} // closed on process exit
	Exit     int
	Error    error
}

func NewProcess(options Options) (p Process, err error) {
	parts, err := shlex.Split(options.Command)
	if err != nil {
		return
	}

	if !path.IsAbs(parts[0]) {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return
		}
		parts[0] = path.Join(wd, parts[0])
	}

	p.cmd = exec.Command(parts[0], parts[1:]...)
	p.cmd.Dir = options.Dir

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

	err = p.cmd.Wait()
	if exiterr, ok := err.(*exec.ExitError); ok {
		p.Exit = exiterr.ExitCode()
	} else {
		return err
	}
	return nil
}

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

func (p *Process) OnExit() chan struct{} {
	return p.exitChan
}

func (p *Process) IsRunning() bool {
	return p.started && !p.ended
}

func (p *Process) Kill() error {
	if !p.IsRunning() {
		return fmt.Errorf("process is not running")
	}

	return p.cmd.Process.Kill()
}
