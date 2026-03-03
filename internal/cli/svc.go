package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ubitransports/beavers/internal/app"
	"github.com/ubitransports/beavers/internal/discovery"
)

func NewSvcCmd(appCtx *app.App) *cobra.Command {
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
				p, err := discovery.ResolveProjectByAlias(appCtx.Projects, args[0])
				if err != nil {
					return err
				}
				
				return appCtx.Exec.RunInteractive(p.Path, "make", cmd.Name())
			},
		}
		svcCmd.AddCommand(actionCmd)
	}

	svcCmd.AddCommand(NewAuditCmd(appCtx))

	return svcCmd
}
