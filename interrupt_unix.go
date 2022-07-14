//go:build !windows

package main

import (
	"os"
	"os/exec"
)

// interruptCmd here only implemented for completeness. BypaSSH does not makes
// any sense to be ran outside the VSCode+WSL environment.
func interruptCmd(cmd *exec.Cmd) error {
	return cmd.Process.Signal(os.Interrupt)
}
