package cli

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/ubitransports/beavers/internal/app"
	"github.com/ubitransports/beavers/internal/discovery"
)

func NewPathCmd(appCtx *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "path <alias>",
		Short: "Resolve and print the absolute path of a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Disable pterm output for the path command to avoid corrupting shell cd $(beavers path ...)
			pterm.DisableOutput()

			p, err := discovery.ResolveProjectByAlias(appCtx.Projects, args[0])
			if err != nil {
				// Print error to stderr so it doesn't break cd, then exit
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}
			fmt.Fprintln(os.Stdout, p.Path)
			return nil
		},
	}
}
