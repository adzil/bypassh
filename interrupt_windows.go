package main

import (
	"os/exec"

	"golang.org/x/sys/windows"
)

func interruptCmd(cmd *exec.Cmd) error {
	return windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(cmd.Process.Pid))
}
