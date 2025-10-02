//go:build windows

package process

import (
	"os"
	"os/exec"
	"syscall"
)

func terminateProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// pauseProcess is a no-op on Windows.
// Windows doesn't have a direct equivalent to Unix SIGSTOP.
func pauseProcess(pid int) error {
	// there is no way afaik to pause process on windows :(
	return nil
}

// resumeProcess is a no-op on Windows.
// Windows doesn't have a direct equivalent to Unix SIGCONT.
func resumeProcess(pid int) error {
	// there is no way afaik to resume process on windows :(
	return nil
}

// setProcessGroupID is a minimal implementation on Windows.
// Windows doesn't have process groups like Unix systems.
func setProcessGroupID(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{}
}
