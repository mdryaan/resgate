package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mdryaan/resgate/internal/utils"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/mdryaan/resgate/pkg/pool"
	"github.com/spf13/cobra"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Live watch of reservation state with periodic refresh",
	RunE: func(cmd *cobra.Command, args []string) error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		printWatchView()

		ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				engine.SweepExpired()
				fmt.Print("\033[H\033[2J")
				printWatchView()
			case <-sig:
				fmt.Println()
				output.Println("Watch stopped.")
				return nil
			}
		}
	},
}

func printWatchView() {
	pools := engine.ListPools()
	reservations := engine.ActiveReservations()

	output.Header(fmt.Sprintf("Resgate Watch  [%s]  (Ctrl+C to exit)", utils.FormatTime(time.Now())))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("  Active Reservations: %s    Pools: %s    Tenants: %s\n",
		output.Info(fmt.Sprintf("%d", len(reservations))),
		output.Info(fmt.Sprintf("%d", len(pools))),
		output.Info(fmt.Sprintf("%d", len(engine.ListTenants()))),
	)
	fmt.Println()

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

	if len(reservations) > 0 {
		t := output.NewTable([]string{"ID", "Tenant", "Pool", "CPU", "Mem MB", "GPU", "Expires"})
		for _, r := range reservations {
			t.Append([]string{
				output.Dim(r.ID),
				output.Info(r.TenantName),
				r.PoolName,
				fmt.Sprintf("%d", r.Resources.CPU),
				fmt.Sprintf("%d", r.Resources.Memory),
				fmt.Sprintf("%d", r.Resources.GPU),
				utils.TimeUntil(r.ExpiresAt),
			})
		}
		t.Render()
	} else {
		output.Warnf("No active reservations")
	}
}

func init() {
	watchCmd.Flags().IntVar(&watchInterval, "interval", 3, "refresh interval in seconds")
	rootCmd.AddCommand(watchCmd)
}
