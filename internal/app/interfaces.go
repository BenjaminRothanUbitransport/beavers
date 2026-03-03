package app

type GitClient interface {
	DetectStatus(path string) (branch, status string)
}

type CommandExecutor interface {
	// RunInteractive executes a command and streams output to os.Stdout/Stderr
	RunInteractive(dir, command string, args ...string) error
	// Run executes a command and returns its stdout
	Run(dir, command string, args ...string) ([]byte, error)
}
