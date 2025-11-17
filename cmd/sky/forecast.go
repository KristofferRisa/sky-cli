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
	// Forecast command flags
	forecastLocation string
	forecastLat      float64
	forecastLon      float64
	forecastHoursCmd int
	forecastFormat   string
)

// forecastCmd represents the forecast command
var forecastCmd = &cobra.Command{
	Use:   "forecast [location]",
	Short: "Get weather forecast",
	Long: `Get hourly weather forecast for a location.

You can specify a location by:
  - Name (from saved locations): sky forecast stavern
  - Coordinates: sky forecast --lat 59.0 --lon 10.0
  - Default location (if no arguments): sky forecast

Examples:
  sky forecast                         # Use default location (12 hours)
  sky forecast stavern                 # Use saved location
  sky forecast --lat 59.0 --lon 10.0  # Use coordinates
  sky forecast --hours 24              # 24-hour forecast
  sky forecast --format json           # JSON output
  sky forecast --format summary        # Brief summary`,
	RunE: runForecast,
}

func init() {
	forecastCmd.Flags().StringVarP(&forecastLocation, "location", "l", "", "Location name from config")
	forecastCmd.Flags().Float64Var(&forecastLat, "lat", 0, "Latitude")
	forecastCmd.Flags().Float64Var(&forecastLon, "lon", 0, "Longitude")
	forecastCmd.Flags().IntVar(&forecastHoursCmd, "hours", 12, "Number of hours for forecast")
	forecastCmd.Flags().StringVarP(&forecastFormat, "format", "f", "", "Output format (full, json, summary, markdown)")

	rootCmd.AddCommand(forecastCmd)
}

func runForecast(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Determine location
	loc, err := getForecastLocation(args)
	if err != nil {
		return err
	}

	// Create weather client (with caching if enabled)
	client := getWeatherClient()

	// Fetch forecast
	forecast, err := client.GetHourlyForecast(ctx, loc, forecastHoursCmd)
	if err != nil {
		return fmt.Errorf("failed to fetch forecast: %w", err)
	}

	// Determine format
	format := forecastFormat
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
	return fmtr.FormatForecast(os.Stdout, forecast, opts)
}

// getForecastLocation determines the location from command arguments and flags
func getForecastLocation(args []string) (*models.Location, error) {
	// Priority 1: Coordinates from flags
	if forecastLat != 0 || forecastLon != 0 {
		if forecastLat == 0 || forecastLon == 0 {
			return nil, fmt.Errorf("both --lat and --lon must be specified")
		}
		loc := &models.Location{
			Latitude:  forecastLat,
			Longitude: forecastLon,
		}
		if err := loc.Validate(); err != nil {
			return nil, err
		}
		return loc, nil
	}

	// Priority 2: Location name from flag
	if forecastLocation != "" {
		return cfg.GetLocation(forecastLocation)
	}

	// Priority 3: Location name from argument
	if len(args) > 0 {
		return cfg.GetLocation(args[0])
	}

	// Priority 4: Default location from config
	return cfg.GetDefaultLocation()
}
