package cmd

import (
	"fmt"
	"time"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/internal/utils"
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

var tenantListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered tenants",
	RunE: func(cmd *cobra.Command, args []string) error {
		tenants := engine.ListTenants()
		if len(tenants) == 0 {
			output.Warnf("No tenants registered")
			return nil
		}
		t := output.NewTable([]string{"Name", "Priority", "Registered"})
		for _, tenant := range tenants {
			pStr := fmt.Sprintf("%d", tenant.Priority)
			var colored string
			switch {
			case tenant.Priority >= 8:
				colored = output.Error(pStr)
			case tenant.Priority >= 5:
				colored = output.Warn(pStr)
			default:
				colored = output.Success(pStr)
			}
			t.Append([]string{output.Info(tenant.Name), colored, utils.FormatTime(tenant.CreatedAt)})
		}
		t.Render()
		return nil
	},
}

func init() {
	tenantAddCmd.Flags().StringVar(&tenantName, "name", "", "tenant name (required)")
	tenantAddCmd.Flags().IntVar(&tenantPriority, "priority", 5, "priority level 1-10 (higher = more priority)")

	tenantCmd.AddCommand(tenantAddCmd, tenantListCmd)
	rootCmd.AddCommand(tenantCmd)
}
