package cmd

import (
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/spf13/cobra"
)

var unreserveID string

var unreserveCmd = &cobra.Command{
	Use:   "unreserve",
	Short: "Release a reservation by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if unreserveID == "" {
			output.Fatal("--id is required")
		}
		if err := engine.Unreserve(unreserveID); err != nil {
			output.Fatal(err.Error())
		}
		output.Successf("Reservation '%s' released", unreserveID)
		return nil
	},
}

func init() {
	unreserveCmd.Flags().StringVar(&unreserveID, "id", "", "reservation ID to release (required)")
	rootCmd.AddCommand(unreserveCmd)
}
