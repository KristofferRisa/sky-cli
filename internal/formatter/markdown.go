package formatter

import (
	"fmt"
	"io"

	"github.com/kristofferrisa/sky-cli/internal/models"
	"github.com/kristofferrisa/sky-cli/internal/ui"
)

// MarkdownFormatter provides markdown-formatted output
type MarkdownFormatter struct{}

// NewMarkdownFormatter creates a new markdown formatter
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

// Name returns the formatter name
func (f *MarkdownFormatter) Name() string {
	return "markdown"
}

// FormatCurrent formats current weather as markdown
func (f *MarkdownFormatter) FormatCurrent(w io.Writer, weather *models.Weather, opts Options) error {
	emoji, description := ui.WeatherSymbol(weather.Symbol)
	if opts.NoEmoji {
		emoji = ""
	}

	fmt.Fprintf(w, "# Weather for %s\n\n", weather.Location)
	fmt.Fprintf(w, "**Updated:** %s\n\n", weather.UpdatedAt.Format("2006-01-02 15:04:05"))

	fmt.Fprintln(w, "## Current Conditions")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "- **Conditions:** %s %s\n", emoji, description)
	fmt.Fprintf(w, "- **Temperature:** %.1f°C\n", weather.Temperature)
	fmt.Fprintf(w, "- **Humidity:** %.0f%%\n", weather.Humidity)
	fmt.Fprintf(w, "- **Cloud Cover:** %.0f%%\n", weather.CloudCover)
	fmt.Fprintf(w, "- **Wind:** %.1f m/s from %s (%s)\n",
		weather.WindSpeed,
		weather.WindDirection(),
		weather.WindDescription())
	fmt.Fprintf(w, "- **Pressure:** %.1f hPa\n", weather.Pressure)
	fmt.Fprintf(w, "- **Precipitation (next hour):** %.1f mm\n", weather.Precipitation)
	fmt.Fprintln(w)

	return nil
}

// FormatForecast formats forecast as markdown
func (f *MarkdownFormatter) FormatForecast(w io.Writer, forecast *models.Forecast, opts Options) error {
	fmt.Fprintf(w, "## Hourly Forecast (%d hours)\n\n", len(forecast.Hours))

	fmt.Fprintln(w, "| Time | Conditions | Temp | Precip | Wind | Humidity |")
	fmt.Fprintln(w, "|------|-----------|------|--------|------|----------|")

	for _, hour := range forecast.Hours {
		emoji, description := ui.WeatherSymbol(hour.Symbol)
		if opts.NoEmoji {
			emoji = ""
		}

		fmt.Fprintf(w, "| %s | %s %s | %.1f°C | %.1fmm | %.1fm/s | %.0f%% |\n",
			hour.Time.Format("15:04"),
			emoji,
			description,
			hour.Temperature,
			hour.Precipitation,
			hour.WindSpeed,
			hour.Humidity,
		)
	}

	fmt.Fprintln(w)
	return nil
}

// FormatDailySummary formats daily summary as markdown
func (f *MarkdownFormatter) FormatDailySummary(w io.Writer, summary *models.DailySummary, opts Options) error {
	fmt.Fprintf(w, "## Daily Summary (%s)\n\n", summary.Date.Format("Monday, January 2"))

	fmt.Fprintln(w, "### Temperature")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "- **Minimum:** %.1f°C\n", summary.TemperatureMin)
	fmt.Fprintf(w, "- **Maximum:** %.1f°C\n", summary.TemperatureMax)
	fmt.Fprintf(w, "- **Average:** %.0f°C\n", summary.TemperatureAvg)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "### Precipitation")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "- **Total (24h):** %.1f mm\n", summary.PrecipitationTotal)
	fmt.Fprintln(w)

	return nil
}

// FormatDailyForecast formats daily forecast as markdown
func (f *MarkdownFormatter) FormatDailyForecast(w io.Writer, dailyForecast *models.DailyForecast, opts Options) error {
	fmt.Fprintf(w, "## Daily Forecast (%d days)\n\n", len(dailyForecast.Days))

	fmt.Fprintln(w, "| Date | Conditions | Temp (Min/Max) | Precip | Wind Max |")
	fmt.Fprintln(w, "|------|-----------|----------------|--------|----------|")

	for _, day := range dailyForecast.Days {
		emoji, description := ui.WeatherSymbol(day.Symbol)
		if opts.NoEmoji {
			emoji = ""
		}

		fmt.Fprintf(w, "| %s | %s %s | %.1f-%.1f°C | %.1fmm | %.1fm/s |\n",
			day.Date.Format("Mon Jan 2"),
			emoji,
			description,
			day.TemperatureMin,
			day.TemperatureMax,
			day.PrecipitationTotal,
			day.WindSpeedMax,
		)
	}

	fmt.Fprintln(w)
	return nil
}
