package formatter

import (
	"io"

	"github.com/kristofferrisa/sky-cli/internal/models"
)

// Options contains formatting options
type Options struct {
	NoColor    bool
	NoEmoji    bool
	TimeFormat string
}

// Formatter is the interface for weather data formatters
type Formatter interface {
	// FormatCurrent formats current weather data
	FormatCurrent(w io.Writer, weather *models.Weather, opts Options) error

	// FormatForecast formats hourly forecast data
	FormatForecast(w io.Writer, forecast *models.Forecast, opts Options) error

	// FormatDailySummary formats daily summary data
	FormatDailySummary(w io.Writer, summary *models.DailySummary, opts Options) error

	// Name returns the formatter name
	Name() string
}

// DefaultOptions returns default formatting options
func DefaultOptions() Options {
	return Options{
		NoColor:    false,
		NoEmoji:    false,
		TimeFormat: "2006-01-02 15:04:05",
	}
}
