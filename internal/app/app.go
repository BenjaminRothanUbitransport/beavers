package app

import (
	"github.com/ubitransports/beavers/internal/config"
)

type App struct {
	Config   *config.Config
	Projects []config.Project
	Git      GitClient
	Exec     CommandExecutor
}
