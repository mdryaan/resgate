package cmd

import (
	"fmt"

	"github.com/mdryaan/resgate/internal/config"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s  %s\n", output.Bold(config.AppName), output.Info("v"+config.Version))
		fmt.Println("Resource Reservation Engine for multi-tenant environments")
		fmt.Println("Built with cobra · viper · sync · encoding/json")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
