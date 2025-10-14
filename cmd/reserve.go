package cmd

import (
	"fmt"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/internal/utils"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/spf13/cobra"
)

var (
	reserveTenant string
	reservePool   string
	reserveCPU    int
	reserveMemory int
	reserveGPU    int
	reserveTTL    int
)

var reserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve resources for a tenant",
	RunE: func(cmd *cobra.Command, args []string) error {
		if reserveTenant == "" || reservePool == "" {
			output.Fatal("--tenant and --pool are required")
		}
		res := models.Resources{CPU: reserveCPU, Memory: reserveMemory, GPU: reserveGPU}
		if res.IsZero() {
			output.Fatal("at least one of --cpu, --memory, --gpu must be non-zero")
		}
		r, err := engine.Reserve(reserveTenant, reservePool, res, reserveTTL)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.Successf("Reservation created: %s", r.ID)
		fmt.Printf("  Tenant:  %s\n", output.Info(r.TenantName))
		fmt.Printf("  Pool:    %s\n", output.Info(r.PoolName))
		fmt.Printf("  CPU:     %d cores\n", r.Resources.CPU)
		fmt.Printf("  Memory:  %d MB\n", r.Resources.Memory)
		fmt.Printf("  GPU:     %d units\n", r.Resources.GPU)
		fmt.Printf("  Expires: %s\n", utils.TimeUntil(r.ExpiresAt))
		return nil
	},
}

func init() {
	reserveCmd.Flags().StringVar(&reserveTenant, "tenant", "", "tenant name (required)")
	reserveCmd.Flags().StringVar(&reservePool, "pool", "", "pool name (required)")
	reserveCmd.Flags().IntVar(&reserveCPU, "cpu", 0, "CPU cores to reserve")
	reserveCmd.Flags().IntVar(&reserveMemory, "memory", 0, "memory in MB to reserve")
	reserveCmd.Flags().IntVar(&reserveGPU, "gpu", 0, "GPU units to reserve")
	reserveCmd.Flags().IntVar(&reserveTTL, "ttl", 0, "TTL in seconds (0 = no expiry)")
	rootCmd.AddCommand(reserveCmd)
}
