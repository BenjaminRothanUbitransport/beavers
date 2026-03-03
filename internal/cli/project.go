package cli

import (
	"io"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/ubitransports/beavers/internal/app"
	"github.com/ubitransports/beavers/internal/config"
)

func NewProjectCmd(appCtx *app.App) *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manage and list projects",
	}

	projectListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all discovered projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProjects(os.Stdout, appCtx.Projects)
		},
	}

	projectCmd.AddCommand(projectListCmd)
	return projectCmd
}

func listProjects(w io.Writer, projects []config.Project) error {
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
