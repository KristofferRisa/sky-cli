package main

import (
	"fmt"
	"os"

	"github.com/kristofferrisa/sky-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config

	// Global flags
	noColor bool
	noEmoji bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "sky",
	Short: "A beautiful weather CLI tool",
	Long: `Sky is a command-line weather tool that provides current conditions,
forecasts, and weather data in a human-readable and LLM-friendly format.

Powered by MET Norway (Meteorologisk institutt).`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override config with flags if set
		if noColor {
			cfg.NoColor = true
		}
		if noEmoji {
			cfg.NoEmoji = true
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&noEmoji, "no-emoji", false, "Disable emoji output")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
