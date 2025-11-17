package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kristofferrisa/sky-cli/internal/formatter"
	"github.com/kristofferrisa/sky-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	// Daily command flags
	dailyLocation string
	dailyLat      float64
	dailyLon      float64
	dailyDays     int
	dailyFormat   string
)

// dailyCmd represents the daily command
var dailyCmd = &cobra.Command{
	Use:   "daily [location]",
	Short: "Get daily weather forecast",
	Long: `Get daily weather forecast for a location.

You can specify a location by:
  - Name (from saved locations): sky daily stavern
  - Coordinates: sky daily --lat 59.0 --lon 10.0
  - Default location (if no arguments): sky daily

Examples:
  sky daily                       # 7-day forecast (default location)
  sky daily stavern               # 7-day forecast for saved location
  sky daily --lat 59.0 --lon 10.0 # Forecast for coordinates
  sky daily --days 3              # 3-day forecast
  sky daily --days 10             # 10-day forecast
  sky daily --format json         # JSON output
  sky daily --format summary      # Brief summary`,
	RunE: runDaily,
}

func init() {
	dailyCmd.Flags().StringVarP(&dailyLocation, "location", "l", "", "Location name from config")
	dailyCmd.Flags().Float64Var(&dailyLat, "lat", 0, "Latitude")
	dailyCmd.Flags().Float64Var(&dailyLon, "lon", 0, "Longitude")
	dailyCmd.Flags().IntVar(&dailyDays, "days", 7, "Number of days for forecast (default: 7)")
	dailyCmd.Flags().StringVarP(&dailyFormat, "format", "f", "", "Output format (full, json, summary, markdown)")

	rootCmd.AddCommand(dailyCmd)
}

func runDaily(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Determine location
	loc, err := getDailyLocation(args)
	if err != nil {
		return err
	}

	// Create weather client (with caching if enabled)
	client := getWeatherClient()

	// Fetch daily forecast
	dailyForecast, err := client.GetDailyForecast(ctx, loc, dailyDays)
	if err != nil {
		return fmt.Errorf("failed to fetch daily forecast: %w", err)
	}

	// Determine format
	format := dailyFormat
	if format == "" {
		format = cfg.DefaultFormat
	}
	if format == "" {
		format = "full"
	}

	// Get formatter
	fmtr, err := formatter.GetFormatter(format)
	if err != nil {
		return err
	}

	// Format options
	opts := formatter.Options{
		NoColor:    cfg.NoColor,
		NoEmoji:    cfg.NoEmoji,
		TimeFormat: "2006-01-02 15:04:05",
	}

	// Format and display
	return fmtr.FormatDailyForecast(os.Stdout, dailyForecast, opts)
}

// getDailyLocation determines the location from command arguments and flags
func getDailyLocation(args []string) (*models.Location, error) {
	// Priority 1: Coordinates from flags
	if dailyLat != 0 || dailyLon != 0 {
		if dailyLat == 0 || dailyLon == 0 {
			return nil, fmt.Errorf("both --lat and --lon must be specified")
		}
		loc := &models.Location{
			Latitude:  dailyLat,
			Longitude: dailyLon,
		}
		if err := loc.Validate(); err != nil {
			return nil, err
		}
		return loc, nil
	}

	// Priority 2: Location name from flag
	if dailyLocation != "" {
		return cfg.GetLocation(dailyLocation)
	}

	// Priority 3: Location name from argument
	if len(args) > 0 {
		return cfg.GetLocation(args[0])
	}

	// Priority 4: Default location from config
	return cfg.GetDefaultLocation()
}
