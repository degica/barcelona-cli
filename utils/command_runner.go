package utils

import (
	"os"
	"os/exec"
)

type CommandRunner struct {
}

func (cr CommandRunner) RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
