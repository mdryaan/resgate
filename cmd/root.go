package cmd

import (
	"os"

	"github.com/mdryaan/resgate/internal/config"
	"github.com/mdryaan/resgate/internal/store"
	"github.com/mdryaan/resgate/pkg/reservation"
	"github.com/spf13/cobra"
)

var (
	engine    *reservation.Engine
	stateFile string
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:   "resgate",
	Short: "Resource reservation engine for multi-tenant environments",
	Long: `Resgate manages resource reservations across multiple tenants sharing
a common resource pool with priority-based preemption, TTL expiry,
conflict detection, and concurrent-safe access.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" {
			return nil
		}
		s := store.New(stateFile)
		if err := s.Load(); err != nil {
			return err
		}
		engine = reservation.NewEngine(s)
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if engine == nil || cmd.Name() == "watch" {
			return nil
		}
		return engine.Save()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&stateFile, "state", config.DefaultStateFile(), "path to state file")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")
}
