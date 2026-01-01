package formatter

import (
	"fmt"
	"io"

	"github.com/kristofferrisa/sky-cli/internal/models"
	"github.com/kristofferrisa/sky-cli/internal/ui"
)

// SummaryFormatter provides brief one-line weather summaries
type SummaryFormatter struct{}

// NewSummaryFormatter creates a new summary formatter
func NewSummaryFormatter() *SummaryFormatter {
	return &SummaryFormatter{}
}

// Name returns the formatter name
func (f *SummaryFormatter) Name() string {
	return "summary"
}

// FormatCurrent formats current weather as a brief summary
func (f *SummaryFormatter) FormatCurrent(w io.Writer, weather *models.Weather, opts Options) error {
	if opts.NoColor {
		ui.DisableColors()
	}

	emoji, description := ui.WeatherSymbol(weather.Symbol)
	if opts.NoEmoji {
		emoji = ""
		description = stripEmoji(description)
	}

	fmt.Fprintf(w, "%s %s: %s %.1f°C (feels like %.1f°C), Wind: %.1f m/s %s, Humidity: %.0f%%",
		ui.Bold(weather.Location.String()),
		ui.Cyan(weather.Timestamp.Format("15:04")),
		emoji+" "+description,
		weather.Temperature,
		weather.FeelsLike(),
		weather.WindSpeed,
		weather.WindDirection(),
		weather.Humidity,
	)

	if weather.Precipitation > 0 {
		fmt.Fprintf(w, ", Rain: %.1fmm", weather.Precipitation)
	}

	fmt.Fprintln(w)
	return nil
}

// FormatForecast formats forecast as brief summaries
func (f *SummaryFormatter) FormatForecast(w io.Writer, forecast *models.Forecast, opts Options) error {
	if opts.NoColor {
		ui.DisableColors()
	}

	fmt.Fprintln(w, ui.Bold(forecast.Location.String()))
	fmt.Fprintln(w, ui.Bold("Hourly Forecast:"))

	for _, hour := range forecast.Hours {
		emoji, description := ui.WeatherSymbol(hour.Symbol)
		if opts.NoEmoji {
			emoji = ""
			description = stripEmoji(description)
		}

		precip := ""
		if hour.Precipitation > 0 {
			precip = fmt.Sprintf(", %.1fmm", hour.Precipitation)
		}

		fmt.Fprintf(w, "  %s: %s %.1f°C (feels like %.1f°C)%s\n",
			ui.Cyan(hour.Time.Format("15:04")),
			emoji+" "+description,
			hour.Temperature,
			hour.FeelsLike(),
			precip,
		)
	}

	return nil
}

// FormatDailySummary formats daily summary as brief text
func (f *SummaryFormatter) FormatDailySummary(w io.Writer, summary *models.DailySummary, opts Options) error {
	if opts.NoColor {
		ui.DisableColors()
	}

	fmt.Fprintf(w, "%s %s: %.1f-%.1f°C (avg %.0f°C)",
		ui.Bold(summary.Location.String()),
		ui.Cyan(summary.Date.Format("Mon Jan 2")),
		summary.TemperatureMin,
		summary.TemperatureMax,
		summary.TemperatureAvg,
	)

	if summary.PrecipitationTotal > 0 {
		fmt.Fprintf(w, ", Rain: %.1fmm", summary.PrecipitationTotal)
	}

	fmt.Fprintln(w)
	return nil
}

// FormatDailyForecast formats daily forecast as brief summaries
func (f *SummaryFormatter) FormatDailyForecast(w io.Writer, dailyForecast *models.DailyForecast, opts Options) error {
	if opts.NoColor {
		ui.DisableColors()
	}

	fmt.Fprintln(w, ui.Bold(dailyForecast.Location.String()))
	fmt.Fprintln(w, ui.Bold("Daily Forecast:"))

	for _, day := range dailyForecast.Days {
		emoji, description := ui.WeatherSymbol(day.Symbol)
		if opts.NoEmoji {
			emoji = ""
			description = stripEmoji(description)
		}

		precipStr := ""
		if day.PrecipitationTotal > 0 {
			precipStr = fmt.Sprintf(", %.1fmm rain", day.PrecipitationTotal)
		}

		fmt.Fprintf(w, "  %s: %s %.1f-%.1f°C%s\n",
			ui.Cyan(day.Date.Format("Mon Jan 2")),
			emoji+" "+description,
			day.TemperatureMin,
			day.TemperatureMax,
			precipStr,
		)
	}

	return nil
}
