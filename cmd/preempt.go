package cmd

import (
	"fmt"

	"github.com/mdryaan/resgate/pkg/output"
	"github.com/spf13/cobra"
)

var (
	preemptTenant string
	preemptTarget string
)

var preemptCmd = &cobra.Command{
	Use:   "preempt",
	Short: "Preempt a lower-priority tenant's reservations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if preemptTenant == "" || preemptTarget == "" {
			output.Fatal("--tenant and --target are required")
		}
		preempted, err := engine.Preempt(preemptTenant, preemptTarget)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.Successf("Preempted %d reservation(s) from '%s'", len(preempted), preemptTarget)
		for _, r := range preempted {
			fmt.Printf("  %s  cpu=%d mem=%d gpu=%d  pool=%s\n",
				output.Warn(r.ID),
				r.Resources.CPU, r.Resources.Memory, r.Resources.GPU,
				r.PoolName,
			)
		}
		return nil
	},
}

func init() {
	preemptCmd.Flags().StringVar(&preemptTenant, "tenant", "", "preemptor tenant name (higher priority, required)")
	preemptCmd.Flags().StringVar(&preemptTarget, "target", "", "target tenant to preempt (lower priority, required)")
	rootCmd.AddCommand(preemptCmd)
}
