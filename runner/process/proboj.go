package process

import (
	"bufio"
	"strings"
)

type ProbojProcess struct {
	Process
	stdoutScanner *bufio.Scanner
	stderrScanner *bufio.Scanner
}

func probojSplitFunc(buffer []byte, eof bool) (int, []byte, error) {
	before, _, found := strings.Cut(string(buffer), "\n.\n")
	if !found {
		return 0, nil, nil
	}
	token := []byte(before)
	return len(token) + 3, token, nil // + 3 so we also advance over "\n.\n"
}

func NewProbojProcess(command string, dir string) (pp ProbojProcess, err error) {
	pp.Process, err = NewProcess(Options{
		Command: command,
		Dir:     dir,
		Stdin:   true,
		Stdout:  true,
		Stderr:  false, // todo: logs
	})
	if err != nil {
		return
	}

	pp.stdoutScanner = bufio.NewScanner(pp.Process.Stdout)
	pp.stdoutScanner.Split(probojSplitFunc)
	return
}

func (pp *ProbojProcess) Write(data string) error {
	_, err := pp.Process.Stdin.Write([]byte(data))
	return err
}

func (pp *ProbojProcess) Read() (string, error) {
	if !pp.stdoutScanner.Scan() && pp.stdoutScanner.Err() != nil {
		return "", pp.stdoutScanner.Err()
	}
	return pp.stdoutScanner.Text(), nil
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
