package cmd

import (
	"fmt"
	"os"

	"github.com/mdryaan/resgate/pkg/exporter"
	"github.com/mdryaan/resgate/pkg/output"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	exportOutput string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export reservation history as JSON or CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		exp, err := exporter.New(exportFormat)
		if err != nil {
			output.Fatal(err.Error())
		}

		report := &exporter.Report{
			Pools:        engine.ListPools(),
			Tenants:      engine.ListTenants(),
			Reservations: engine.ActiveReservations(),
			History:      engine.History(),
		}

		data, err := exp.Export(report)
		if err != nil {
			output.Fatal(err.Error())
		}

		if exportOutput != "" {
			if err := os.WriteFile(exportOutput, data, 0644); err != nil {
				output.Fatal(err.Error())
			}
			output.Successf("Report exported to %s (%d bytes)", exportOutput, len(data))
			return nil
		}

		fmt.Print(string(data))
		return nil
	},
}

func init() {
	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "export format: json, csv")
	exportCmd.Flags().StringVar(&exportOutput, "output", "", "write to file instead of stdout")
	rootCmd.AddCommand(exportCmd)
}
