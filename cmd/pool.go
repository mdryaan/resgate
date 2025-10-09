package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/internal/utils"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/mdryaan/resgate/pkg/pool"
	"github.com/spf13/cobra"
)

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Manage resource pools",
}

var (
	poolName   string
	poolCPU    int
	poolMemory int
	poolGPU    int
)

var poolCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new resource pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		if poolName == "" {
			output.Fatal("--name is required")
		}
		p := &models.Pool{
			Name:      poolName,
			Total:     models.Resources{CPU: poolCPU, Memory: poolMemory, GPU: poolGPU},
			CreatedAt: time.Now(),
		}
		if err := engine.CreatePool(p); err != nil {
			output.Fatal(err.Error())
		}
		output.Successf("Pool '%s' created — CPU: %d  Memory: %dMB  GPU: %d", poolName, poolCPU, poolMemory, poolGPU)
		return nil
	},
}

var poolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all resource pools with utilization",
	RunE: func(cmd *cobra.Command, args []string) error {
		pools := engine.ListPools()
		if len(pools) == 0 {
			output.Warnf("No pools registered")
			return nil
		}
		t := output.NewTable([]string{"Name", "CPU (used/total)", "Memory MB (used/total)", "GPU (used/total)", "Status", "Created"})
		for _, p := range pools {
			status := pool.UtilizationStatus(p)
			var statusStr string
			switch status {
			case pool.StatusHealthy:
				statusStr = output.Success("healthy")
			case pool.StatusWarning:
				statusStr = output.Warn("warning")
			case pool.StatusExhausted:
				statusStr = output.Error("exhausted")
			}
			t.Append([]string{
				output.Info(p.Name),
				fmt.Sprintf("%d / %d (%.0f%%)", p.Reserved.CPU, p.Total.CPU, p.CPUPercent()),
				fmt.Sprintf("%d / %d (%.0f%%)", p.Reserved.Memory, p.Total.Memory, p.MemoryPercent()),
				fmt.Sprintf("%d / %d (%.0f%%)", p.Reserved.GPU, p.Total.GPU, p.GPUPercent()),
				statusStr,
				utils.FormatTime(p.CreatedAt),
			})
		}
		t.Render()
		return nil
	},
}

var poolShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a specific pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		if poolName == "" {
			output.Fatal("--name is required")
		}
		p, err := engine.GetPool(poolName)
		if err != nil {
			output.Fatal(err.Error())
		}
		avail := p.Available()
		output.Header(fmt.Sprintf("Pool: %s", p.Name))
		fmt.Println(strings.Repeat("─", 50))
		fmt.Printf("  %-22s %s\n", "Created:", utils.FormatTime(p.CreatedAt))
		fmt.Println()
		fmt.Printf("  %-22s %s\n", "CPU Total:", output.Info(fmt.Sprintf("%d cores", p.Total.CPU)))
		fmt.Printf("  %-22s %s\n", "CPU Reserved:", output.Warn(fmt.Sprintf("%d cores (%.0f%%)", p.Reserved.CPU, p.CPUPercent())))
		fmt.Printf("  %-22s %s\n", "CPU Available:", output.Success(fmt.Sprintf("%d cores", avail.CPU)))
		fmt.Println()
		fmt.Printf("  %-22s %s\n", "Memory Total:", output.Info(fmt.Sprintf("%d MB", p.Total.Memory)))
		fmt.Printf("  %-22s %s\n", "Memory Reserved:", output.Warn(fmt.Sprintf("%d MB (%.0f%%)", p.Reserved.Memory, p.MemoryPercent())))
		fmt.Printf("  %-22s %s\n", "Memory Available:", output.Success(fmt.Sprintf("%d MB", avail.Memory)))
		fmt.Println()
		fmt.Printf("  %-22s %s\n", "GPU Total:", output.Info(fmt.Sprintf("%d units", p.Total.GPU)))
		fmt.Printf("  %-22s %s\n", "GPU Reserved:", output.Warn(fmt.Sprintf("%d units (%.0f%%)", p.Reserved.GPU, p.GPUPercent())))
		fmt.Printf("  %-22s %s\n", "GPU Available:", output.Success(fmt.Sprintf("%d units", avail.GPU)))
		return nil
	},
}

func init() {
	poolCreateCmd.Flags().StringVar(&poolName, "name", "", "pool name (required)")
	poolCreateCmd.Flags().IntVar(&poolCPU, "cpu", 0, "total CPU cores")
	poolCreateCmd.Flags().IntVar(&poolMemory, "memory", 0, "total memory in MB")
	poolCreateCmd.Flags().IntVar(&poolGPU, "gpu", 0, "total GPU units")

	poolShowCmd.Flags().StringVar(&poolName, "name", "", "pool name (required)")

	poolCmd.AddCommand(poolCreateCmd, poolListCmd, poolShowCmd)
	rootCmd.AddCommand(poolCmd)
}
