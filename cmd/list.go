package cmd

import (
	"fmt"

	"github.com/mdryaan/resgate/internal/utils"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active reservations",
	RunE: func(cmd *cobra.Command, args []string) error {
		reservations := engine.ActiveReservations()
		if len(reservations) == 0 {
			output.Warnf("No active reservations")
			return nil
		}
		t := output.NewTable([]string{"ID", "Tenant", "Pool", "CPU", "Memory MB", "GPU", "Status", "Expires"})
		for _, r := range reservations {
			var statusStr string
			switch r.Status {
			case "active":
				statusStr = output.Success(string(r.Status))
			case "expired":
				statusStr = output.Error(string(r.Status))
			case "preempted":
				statusStr = output.Warn(string(r.Status))
			default:
				statusStr = string(r.Status)
			}
			t.Append([]string{
				output.Dim(r.ID),
				output.Info(r.TenantName),
				r.PoolName,
				fmt.Sprintf("%d", r.Resources.CPU),
				fmt.Sprintf("%d", r.Resources.Memory),
				fmt.Sprintf("%d", r.Resources.GPU),
				statusStr,
				utils.TimeUntil(r.ExpiresAt),
			})
		}
		t.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
