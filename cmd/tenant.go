package cmd

import (
	"time"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/spf13/cobra"
)

var tenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Manage tenants",
}

var (
	tenantName     string
	tenantPriority int
)

var tenantAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Register a new tenant",
	RunE: func(cmd *cobra.Command, args []string) error {
		if tenantName == "" {
			output.Fatal("--name is required")
		}
		t := &models.Tenant{
			Name:      tenantName,
			Priority:  tenantPriority,
			CreatedAt: time.Now(),
		}
		if err := engine.AddTenant(t); err != nil {
			output.Fatal(err.Error())
		}
		output.Successf("Tenant '%s' registered with priority %d", tenantName, tenantPriority)
		return nil
	},
}

func init() {
	tenantAddCmd.Flags().StringVar(&tenantName, "name", "", "tenant name (required)")
	tenantAddCmd.Flags().IntVar(&tenantPriority, "priority", 5, "priority level 1-10 (higher = more priority)")

	tenantCmd.AddCommand(tenantAddCmd)
	rootCmd.AddCommand(tenantCmd)
}
