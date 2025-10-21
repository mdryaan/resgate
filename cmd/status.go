package cmd

import (
	"fmt"
	"strings"

	"github.com/mdryaan/resgate/internal/utils"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/mdryaan/resgate/pkg/pool"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Overall system status and utilization",
	RunE: func(cmd *cobra.Command, args []string) error {
		pools := engine.ListPools()
		tenants := engine.ListTenants()
		reservations := engine.ActiveReservations()
		expired := engine.SweepExpired()

		output.Header("Resgate System Status")
		fmt.Println(strings.Repeat("─", 60))
		fmt.Println()

		fmt.Printf("  %-24s %s\n", "Pools:", output.Info(fmt.Sprintf("%d", len(pools))))
		fmt.Printf("  %-24s %s\n", "Tenants:", output.Info(fmt.Sprintf("%d", len(tenants))))
		fmt.Printf("  %-24s %s\n", "Active Reservations:", output.Info(fmt.Sprintf("%d", len(reservations))))
		fmt.Printf("  %-24s %s\n", "Expired (swept):", output.Warn(fmt.Sprintf("%d", len(expired))))
		fmt.Println()

		if len(pools) > 0 {
			output.Header("Pool Utilization")
			for _, p := range pools {
				status := pool.UtilizationStatus(p)
				var colorFn func(string) string
				switch status {
				case pool.StatusHealthy:
					colorFn = output.Success
				case pool.StatusWarning:
					colorFn = output.Warn
				case pool.StatusExhausted:
					colorFn = output.Error
				}
				fmt.Printf("  %s\n", output.Info(p.Name))
				fmt.Printf("    CPU     %s %s\n", utils.PercentBar(p.CPUPercent(), 20), colorFn(fmt.Sprintf("%.0f%%", p.CPUPercent())))
				fmt.Printf("    Memory  %s %s\n", utils.PercentBar(p.MemoryPercent(), 20), colorFn(fmt.Sprintf("%.0f%%", p.MemoryPercent())))
				fmt.Printf("    GPU     %s %s\n", utils.PercentBar(p.GPUPercent(), 20), colorFn(fmt.Sprintf("%.0f%%", p.GPUPercent())))
				fmt.Println()
			}
		}

		if len(tenants) > 0 {
			output.Header("Tenant Priority Ranking")
			for _, t := range tenants {
				count := 0
				for _, r := range reservations {
					if r.TenantName == t.Name {
						count++
					}
				}
				fmt.Printf("  [%d] %-20s %s active reservation(s)\n",
					t.Priority, t.Name, output.Info(fmt.Sprintf("%d", count)))
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
