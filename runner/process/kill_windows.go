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

func pauseProcess(pid int) error {
	// there is no way afaik to pause process on windows :(
	return nil
}

func resumeProcess(pid int) error {
	// there is no way afaik to resume process on windows :(
	return nil
}

func setProcessGroupID(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{}
}
