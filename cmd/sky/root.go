package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kristofferrisa/sky-cli/internal/api"
	"github.com/kristofferrisa/sky-cli/internal/api/met"
	"github.com/kristofferrisa/sky-cli/internal/cache"
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

	// Add version command
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version, commit, and build information for Sky CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sky CLI %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
		if builtBy != "unknown" {
			fmt.Printf("Built by: %s\n", builtBy)
		}
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// getWeatherClient creates a weather client with optional caching
func getWeatherClient() api.WeatherClient {
	// Check if cache is enabled
	if !cfg.Cache.Enabled {
		return met.NewClient()
	}

	// Create cache directory
	cacheDir := cfg.Cache.Directory
	if cacheDir == "" {
		cacheDir = os.Getenv("HOME") + "/.sky/cache"
	}

	// Create file cache
	fileCache, err := cache.NewFileCache(cacheDir)
	if err != nil {
		// Fall back to no-op cache if creation fails
		fmt.Fprintf(os.Stderr, "Warning: Failed to create cache: %v\n", err)
		return met.NewClient()
	}

	// Get TTL from config
	ttl := time.Duration(cfg.Cache.TTLMinutes) * time.Minute
	if ttl == 0 {
		ttl = 10 * time.Minute
	}

	// Return cached client
	return met.NewCachedClient(fileCache, ttl)
}
