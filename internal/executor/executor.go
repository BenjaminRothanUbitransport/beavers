package executor

import (
	"os"
	"os/exec"
)

type RealCommandExecutor struct{}

func NewExecutor() *RealCommandExecutor {
	return &RealCommandExecutor{}
}

func (e *RealCommandExecutor) RunInteractive(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (e *RealCommandExecutor) Run(dir, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	return cmd.Output()
}
