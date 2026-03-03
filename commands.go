package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	configPath string
	cfg        *Config
	projects   []Project
	rootCmd    = &cobra.Command{
		Use:   "beavers",
		Short: "Beavers CLI - The Unified Developer Experience Plane",
		Long:  `Beavers is an industrious CLI tool designed to streamline the management of complex microservice ecosystems.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cfg, err = loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			
			cacheValid, cachedProjects, _ := readCache()
			if cacheValid {
				projects = cachedProjects
				// Kick off background refresh
				go func() {
					freshProjects, _ := discoverProjects(cfg)
					_ = writeCache(freshProjects)
				}()
			} else {
				// Blocking full discovery
				projects, err = discoverProjects(cfg)
				if err != nil {
					return fmt.Errorf("failed to discover projects: %w", err)
				}
				_ = writeCache(projects)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file (default is ./beavers.yaml or ~/.beavers/config.yaml)")

	// Define Project List Command
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manage and list projects",
	}

	projectListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all discovered projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProjects(os.Stdout, projects)
		},
	}

	projectCmd.AddCommand(projectListCmd)

	// Define Path Command
	pathCmd := &cobra.Command{
		Use:   "path <alias>",
		Short: "Resolve and print the absolute path of a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return printProjectPath(os.Stdout, projects, args[0])
		},
	}

	// Define Svc Command
	svcCmd := &cobra.Command{
		Use:   "svc",
		Short: "Manage and execute commands on services",
	}

	for _, action := range []string{"install", "build", "pull"} {
		actionCmd := &cobra.Command{
			Use:          action + " <alias>",
			Short:        fmt.Sprintf("Execute make %s for a project", action),
			Args:         cobra.ExactArgs(1),
			SilenceUsage: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				p, err := resolveProjectByAlias(projects, args[0])
				if err != nil {
					return err
				}
				// Use cmd.Name() instead of cmd.Use because Use has args.
				return RunMakeTarget(p.Path, cmd.Name())
			},
		}
		svcCmd.AddCommand(actionCmd)
	}

	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(pathCmd)
	rootCmd.AddCommand(svcCmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// listProjects prints a table of projects to the provided writer.
func listProjects(w io.Writer, projects []Project) error {
	data := pterm.TableData{
		{"Name", "Alias", "Type", "Workspace", "Git Branch", "Sync Status", "Path"},
	}

	for _, p := range projects {
		data = append(data, []string{
			p.Name,
			p.Alias,
			p.Type,
			p.Workspace,
			p.GitBranch,
			p.SyncStatus,
			p.Path,
		})
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data).WithWriter(w).Render()
}

// printProjectPath prints the absolute path of a project to the provided writer.
func printProjectPath(w io.Writer, projects []Project, identifier string) error {
	p, err := resolveProjectByAlias(projects, identifier)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, p.Path)
	return nil
}
