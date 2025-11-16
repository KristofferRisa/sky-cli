package formatter

import (
	"fmt"
	"io"
	"time"

	"github.com/kristofferrisa/sky-cli/internal/models"
	"github.com/kristofferrisa/sky-cli/internal/ui"
)

// FullFormatter provides detailed weather output similar to the bash script
type FullFormatter struct{}

// NewFullFormatter creates a new full formatter
func NewFullFormatter() *FullFormatter {
	return &FullFormatter{}
}

// Name returns the formatter name
func (f *FullFormatter) Name() string {
	return "full"
}

// FormatCurrent formats current weather data with full details
func (f *FullFormatter) FormatCurrent(w io.Writer, weather *models.Weather, opts Options) error {
	if opts.NoColor {
		ui.DisableColors()
	}

	// Header
	fmt.Fprintln(w, ui.Header(fmt.Sprintf("CURRENT WEATHER - %s", weather.Location)))
	fmt.Fprintf(w, "API: MET Norway (Meteorologisk institutt)\n")
	fmt.Fprintf(w, "Coordinates: %.2f°N, %.2f°E\n", weather.Location.Latitude, weather.Location.Longitude)
	fmt.Fprintf(w, "Request time: %s\n", time.Now().Format(opts.TimeFormat))
	fmt.Fprintln(w)

	// Current conditions
	emoji, description := ui.WeatherSymbol(weather.Symbol)
	if opts.NoEmoji {
		emoji = ""
	}

	fmt.Fprintln(w, ui.Bold("Location:    "), weather.Location.String())
	fmt.Fprintln(w, ui.Bold("Time:        "), formatTime(weather.Timestamp))
	fmt.Fprintln(w, ui.Bold("Updated:     "), formatTime(weather.UpdatedAt))
	fmt.Fprintln(w)

	fmt.Fprintln(w, ui.GreenBold("Conditions:  "), emoji, description)
	fmt.Fprintf(w, "%s  %.1f°C\n", ui.Bold("Temperature:"), weather.Temperature)
	fmt.Fprintf(w, "%s     %.0f%%\n", ui.Bold("Humidity:"), weather.Humidity)
	fmt.Fprintf(w, "%s  %.0f%%\n", ui.Bold("Cloud Cover:"), weather.CloudCover)
	fmt.Fprintf(w, "%s %.1f mm (next hour)\n", ui.Bold("Precipitation:"), weather.Precipitation)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "%s         %.1f m/s from %s\n", ui.Bold("Wind:"), weather.WindSpeed, weather.WindDirection())
	fmt.Fprintf(w, "%s  %s\n", ui.Bold("Wind Status:"), weather.WindDescription())
	fmt.Fprintf(w, "%s     %.1f hPa\n", ui.Bold("Pressure:"), weather.Pressure)
	fmt.Fprintln(w)

	return nil
}

// FormatForecast formats hourly forecast data
func (f *FullFormatter) FormatForecast(w io.Writer, forecast *models.Forecast, opts Options) error {
	if opts.NoColor {
		ui.DisableColors()
	}

	fmt.Fprintln(w, ui.Header(fmt.Sprintf("HOURLY FORECAST (Next %d Hours)", len(forecast.Hours))))
	fmt.Fprintf(w, "%s\n", ui.Bold("Time     Temp    Symbol                Precip  Wind    Humidity"))
	fmt.Fprintln(w, "────────────────────────────────────────────────────────────────────")

	for _, hour := range forecast.Hours {
		_, description := ui.WeatherSymbol(hour.Symbol)
		if opts.NoEmoji {
			description = stripEmoji(description)
		}

		fmt.Fprintf(w, "%-8s %-7s %-20s %-7s %-7s %.0f%%\n",
			hour.Time.Format("15:04"),
			fmt.Sprintf("%.1f°C", hour.Temperature),
			description,
			fmt.Sprintf("%.1fmm", hour.Precipitation),
			fmt.Sprintf("%.1fm/s", hour.WindSpeed),
			hour.Humidity,
		)
	}

	fmt.Fprintln(w)
	return nil
}

// FormatDailySummary formats daily summary data
func (f *FullFormatter) FormatDailySummary(w io.Writer, summary *models.DailySummary, opts Options) error {
	if opts.NoColor {
		ui.DisableColors()
	}

	fmt.Fprintln(w, ui.Header("DAILY SUMMARY"))

	fmt.Fprintln(w, ui.Bold("Temperature Range:"))
	fmt.Fprintf(w, "  • Minimum: %.1f°C\n", summary.TemperatureMin)
	fmt.Fprintf(w, "  • Maximum: %.1f°C\n", summary.TemperatureMax)
	fmt.Fprintf(w, "  • Average: %.0f°C\n", summary.TemperatureAvg)
	fmt.Fprintln(w)

	fmt.Fprintln(w, ui.Bold("Precipitation:"))
	fmt.Fprintf(w, "  • Total (24h): %.1f mm\n", summary.PrecipitationTotal)
	fmt.Fprintln(w)

	return nil
}

// FormatComplete formats current weather with forecast and daily summary
func (f *FullFormatter) FormatComplete(w io.Writer, weather *models.Weather, forecast *models.Forecast, summary *models.DailySummary, opts Options) error {
	// Current weather
	if err := f.FormatCurrent(w, weather, opts); err != nil {
		return err
	}

	// Hourly forecast
	if forecast != nil && len(forecast.Hours) > 0 {
		if err := f.FormatForecast(w, forecast, opts); err != nil {
			return err
		}
	}

	// Daily summary
	if summary != nil {
		if err := f.FormatDailySummary(w, summary, opts); err != nil {
			return err
		}
	}

	// Weather warnings
	if weather.WindSpeed > 15 {
		fmt.Fprintln(w, ui.YellowBold("⚠️  Weather Warnings:"))
		fmt.Fprintf(w, "  • Strong winds detected (%.1f m/s)\n", weather.WindSpeed)
		fmt.Fprintln(w)
	}

	if summary != nil && summary.PrecipitationTotal > 10 {
		fmt.Fprintln(w, ui.YellowBold("⚠️  Weather Warnings:"))
		fmt.Fprintf(w, "  • Significant precipitation expected (%.1f mm)\n", summary.PrecipitationTotal)
		fmt.Fprintln(w)
	}

	// Structured data for LLM
	f.formatStructuredData(w, weather, summary, opts)

	// Footer
	fmt.Fprintln(w, ui.Header("Weather Data Retrieved Successfully"))
	fmt.Fprintln(w, "✅ Data is fresh and ready for analysis")
	fmt.Fprintln(w)

	return nil
}

// formatStructuredData outputs LLM-friendly structured data
func (f *FullFormatter) formatStructuredData(w io.Writer, weather *models.Weather, summary *models.DailySummary, opts Options) {
	fmt.Fprintln(w, ui.Header("STRUCTURED DATA (For LLM Processing)"))

	_, description := ui.WeatherSymbol(weather.Symbol)
	description = stripEmoji(description)

	fmt.Fprintf(w, "LOCATION: %s (%.2f°N, %.2f°E)\n",
		weather.Location.Name,
		weather.Location.Latitude,
		weather.Location.Longitude)
	fmt.Fprintf(w, "TIMESTAMP: %s\n", time.Now().Format(opts.TimeFormat))
	fmt.Fprintln(w, "DATA_SOURCE: MET Norway (yr.no/met.no)")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "CURRENT_CONDITIONS:")
	fmt.Fprintf(w, "  temperature: %.1f°C\n", weather.Temperature)
	fmt.Fprintf(w, "  weather: %s\n", description)
	fmt.Fprintf(w, "  humidity: %.0f%%\n", weather.Humidity)
	fmt.Fprintf(w, "  wind_speed: %.1f m/s\n", weather.WindSpeed)
	fmt.Fprintf(w, "  wind_direction: %s\n", weather.WindDirection())
	fmt.Fprintf(w, "  wind_description: %s\n", weather.WindDescription())
	fmt.Fprintf(w, "  pressure: %.1f hPa\n", weather.Pressure)
	fmt.Fprintf(w, "  cloud_cover: %.0f%%\n", weather.CloudCover)
	fmt.Fprintf(w, "  precipitation_next_hour: %.1f mm\n", weather.Precipitation)
	fmt.Fprintln(w)

	if summary != nil {
		fmt.Fprintln(w, "DAILY_SUMMARY:")
		fmt.Fprintf(w, "  temperature_min: %.1f°C\n", summary.TemperatureMin)
		fmt.Fprintf(w, "  temperature_max: %.1f°C\n", summary.TemperatureMax)
		fmt.Fprintf(w, "  temperature_avg: %.0f°C\n", summary.TemperatureAvg)
		fmt.Fprintf(w, "  precipitation_24h: %.1f mm\n", summary.PrecipitationTotal)
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "DATA_UNITS:")
	fmt.Fprintln(w, "  temperature: celsius")
	fmt.Fprintln(w, "  wind_speed: meters_per_second")
	fmt.Fprintln(w, "  pressure: hectopascal (hPa)")
	fmt.Fprintln(w, "  precipitation: millimeters (mm)")
	fmt.Fprintln(w, "  humidity: percent")
	fmt.Fprintln(w)
}

// Helper functions

func formatTime(t time.Time) string {
	return t.Format("Monday, January 02, 2006 at 15:04")
}

func stripEmoji(s string) string {
	// Simple emoji stripping - remove common weather emojis
	emojis := []string{"☀️", "🌤️", "⛅", "☁️", "🌦️", "🌧️", "⛈️", "🌨️", "❄️", "🌫️", "🌡️"}
	result := s
	for _, emoji := range emojis {
		result = stripString(result, emoji)
	}
	return result
}

func stripString(s, remove string) string {
	result := ""
	for _, c := range s {
		if string(c) != remove {
			result += string(c)
		}
	}
	return result
}
