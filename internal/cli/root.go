package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ubitransports/beavers/internal/app"
	"github.com/ubitransports/beavers/internal/config"
	"github.com/ubitransports/beavers/internal/discovery"
)

var configPath string

func NewRootCmd(appCtx *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "beavers",
		Short: "Beavers CLI - The Unified Developer Experience Plane",
		Long:  `Beavers is an industrious CLI tool designed to streamline the management of complex microservice ecosystems.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			appCtx.Config = cfg

			cacheValid, cachedProjects, _ := discovery.ReadCache()
			if cacheValid {
				appCtx.Projects = cachedProjects
				// Kick off background refresh
				go func() {
					freshProjects, _ := discovery.DiscoverProjects(appCtx)
					_ = discovery.WriteCache(freshProjects)
				}()
			} else {
				projects, err := discovery.DiscoverProjects(appCtx)
				if err != nil {
					return fmt.Errorf("failed to discover projects: %w", err)
				}
				appCtx.Projects = projects
				_ = discovery.WriteCache(projects)
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file (default is ./beavers.yaml or ~/.beavers/config.yaml)")

	cmd.AddCommand(NewProjectCmd(appCtx))
	cmd.AddCommand(NewPathCmd(appCtx))
	cmd.AddCommand(NewSvcCmd(appCtx))

	return cmd
}
