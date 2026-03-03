package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunMakeTarget executes a make target in the given project path.
func RunMakeTarget(projectPath string, target string) error {
	cmd := exec.Command("make", target)
	cmd.Dir = projectPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Optional: we could check if Makefile exists before running
	if _, err := os.Stat(projectPath + "/Makefile"); os.IsNotExist(err) {
		return fmt.Errorf("no Makefile found in %s", projectPath)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make %s failed: %w", target, err)
	}

	return nil
}
