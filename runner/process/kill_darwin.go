//go:build darwin

package process

import (
	"os/exec"
	"syscall"
)

func terminateProcess(pid int) error {
	return syscall.Kill(-pid, syscall.SIGTERM)
}

func pauseProcess(pid int) error {
	return syscall.Kill(-pid, syscall.SIGSTOP)
}

func resumeProcess(pid int) error {
	return syscall.Kill(-pid, syscall.SIGCONT)
}

func setProcessGroupID(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
