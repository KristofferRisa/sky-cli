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
	// Current command flags
	locationName  string
	latitude      float64
	longitude     float64
	showForecast  bool
	showSummary   bool
	forecastHours int
	formatType    string
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current [location]",
	Short: "Get current weather conditions",
	Long: `Get current weather conditions for a location.

You can specify a location by:
  - Name (from saved locations): sky current stavern
  - Coordinates: sky current --lat 59.0 --lon 10.0
  - Default location (if no arguments): sky current

Examples:
  sky current                          # Use default location
  sky current stavern                  # Use saved location 'stavern'
  sky current --lat 59.0 --lon 10.0   # Use coordinates
  sky current --forecast               # Include 12-hour forecast
  sky current --summary                # Include daily summary
  sky current --format json            # JSON output
  sky current --format summary         # Brief summary
  sky current --format markdown        # Markdown format`,
	RunE: runCurrent,
}

func init() {
	currentCmd.Flags().StringVarP(&locationName, "location", "l", "", "Location name from config")
	currentCmd.Flags().Float64Var(&latitude, "lat", 0, "Latitude")
	currentCmd.Flags().Float64Var(&longitude, "lon", 0, "Longitude")
	currentCmd.Flags().BoolVar(&showForecast, "forecast", false, "Include hourly forecast")
	currentCmd.Flags().BoolVar(&showSummary, "summary", false, "Include daily summary")
	currentCmd.Flags().IntVar(&forecastHours, "hours", 12, "Number of hours for forecast")
	currentCmd.Flags().StringVarP(&formatType, "format", "f", "", "Output format (full, json, summary, markdown)")

	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Determine location
	loc, err := getLocation(args)
	if err != nil {
		return err
	}

	// Create weather client (with caching if enabled)
	client := getWeatherClient()

	// Fetch current weather
	weather, err := client.GetCurrentWeather(ctx, loc)
	if err != nil {
		return fmt.Errorf("failed to fetch current weather: %w", err)
	}

	// Fetch forecast if requested
	var forecast *models.Forecast
	if showForecast {
		forecast, err = client.GetHourlyForecast(ctx, loc, forecastHours)
		if err != nil {
			return fmt.Errorf("failed to fetch forecast: %w", err)
		}
	}

	// Fetch daily summary if requested
	var summary *models.DailySummary
	if showSummary {
		summary, err = client.GetDailySummary(ctx, loc)
		if err != nil {
			return fmt.Errorf("failed to fetch daily summary: %w", err)
		}
	}

	// Determine format
	format := formatType
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

	// Special handling for full formatter with complete output
	if fullFmt, ok := fmtr.(*formatter.FullFormatter); ok && showForecast && showSummary {
		return fullFmt.FormatComplete(os.Stdout, weather, forecast, summary, opts)
	}

	// Otherwise, format individually
	if err := fmtr.FormatCurrent(os.Stdout, weather, opts); err != nil {
		return err
	}

	if showForecast && forecast != nil {
		if err := fmtr.FormatForecast(os.Stdout, forecast, opts); err != nil {
			return err
		}
	}

	if showSummary && summary != nil {
		if err := fmtr.FormatDailySummary(os.Stdout, summary, opts); err != nil {
			return err
		}
	}

	return nil
}

// getLocation determines the location from command arguments and flags
func getLocation(args []string) (*models.Location, error) {
	// Priority 1: Coordinates from flags
	if latitude != 0 || longitude != 0 {
		if latitude == 0 || longitude == 0 {
			return nil, fmt.Errorf("both --lat and --lon must be specified")
		}
		loc := &models.Location{
			Latitude:  latitude,
			Longitude: longitude,
		}
		if err := loc.Validate(); err != nil {
			return nil, err
		}
		return loc, nil
	}

	// Priority 2: Location name from flag
	if locationName != "" {
		return cfg.GetLocation(locationName)
	}

	// Priority 3: Location name from argument
	if len(args) > 0 {
		return cfg.GetLocation(args[0])
	}

	// Priority 4: Default location from config
	return cfg.GetDefaultLocation()
}
