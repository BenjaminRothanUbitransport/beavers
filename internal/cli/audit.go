package cli

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/ubitransports/beavers/internal/app"
	"github.com/ubitransports/beavers/internal/audit"
	"github.com/ubitransports/beavers/internal/discovery"
)

func NewAuditCmd(appCtx *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "audit <alias>",
		Short: "Run compliance audit against a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := discovery.ResolveProjectByAlias(appCtx.Projects, args[0])
			if err != nil {
				return err
			}

			if len(appCtx.Config.AuditRules) == 0 {
				pterm.Warning.Println("No audit rules defined in configuration.")
				return nil
			}

			results := audit.RunAudit(appCtx, *p, appCtx.Config.AuditRules)

			data := pterm.TableData{
				{"Rule", "Status", "Message"},
			}

			for _, res := range results {
				statusStr := pterm.Green("PASS")
				if res.Status == "FAIL" {
					statusStr = pterm.Red("FAIL")
				}
				data = append(data, []string{res.RuleName, statusStr, res.Message})
			}

			pterm.DefaultHeader.WithFullWidth().Printf("Audit Report: %s", p.Name)
			return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
		},
	}
}
