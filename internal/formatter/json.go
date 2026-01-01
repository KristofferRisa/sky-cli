package formatter

import (
	"encoding/json"
	"io"

	"github.com/kristofferrisa/sky-cli/internal/models"
)

// JSONFormatter provides JSON output
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Name returns the formatter name
func (f *JSONFormatter) Name() string {
	return "json"
}

// JSONWeather is the JSON representation of weather data
type JSONWeather struct {
	Location      *models.Location `json:"location"`
	Timestamp     string           `json:"timestamp"`
	UpdatedAt     string           `json:"updated_at"`
	Temperature   float64          `json:"temperature"`
	FeelsLike     float64          `json:"feels_like"`
	Humidity      float64          `json:"humidity"`
	Pressure      float64          `json:"pressure"`
	CloudCover    float64          `json:"cloud_cover"`
	WindSpeed     float64          `json:"wind_speed"`
	WindDirection string           `json:"wind_direction"`
	WindDegrees   float64          `json:"wind_degrees"`
	Precipitation float64          `json:"precipitation"`
	Symbol        string           `json:"symbol"`
	Description   string           `json:"description"`
	Units         JSONUnits        `json:"units"`
}

// JSONUnits describes the units used
type JSONUnits struct {
	Temperature   string `json:"temperature"`
	WindSpeed     string `json:"wind_speed"`
	Pressure      string `json:"pressure"`
	Precipitation string `json:"precipitation"`
	Humidity      string `json:"humidity"`
}

// JSONForecast is the JSON representation of forecast data
type JSONForecast struct {
	Location *models.Location     `json:"location"`
	Hours    []JSONHourlyForecast `json:"hours"`
}

// JSONHourlyForecast is a single hour in the forecast
type JSONHourlyForecast struct {
	Time          string  `json:"time"`
	Temperature   float64 `json:"temperature"`
	FeelsLike     float64 `json:"feels_like"`
	Humidity      float64 `json:"humidity"`
	WindSpeed     float64 `json:"wind_speed"`
	Precipitation float64 `json:"precipitation"`
	Symbol        string  `json:"symbol"`
	Description   string  `json:"description"`
}

// JSONDailySummary is the JSON representation of daily summary
type JSONDailySummary struct {
	Location           *models.Location `json:"location"`
	Date               string           `json:"date"`
	TemperatureMin     float64          `json:"temperature_min"`
	TemperatureMax     float64          `json:"temperature_max"`
	TemperatureAvg     float64          `json:"temperature_avg"`
	PrecipitationTotal float64          `json:"precipitation_total"`
	Units              JSONUnits        `json:"units"`
}

// FormatCurrent formats current weather as JSON
func (f *JSONFormatter) FormatCurrent(w io.Writer, weather *models.Weather, opts Options) error {
	jw := JSONWeather{
		Location:      weather.Location,
		Timestamp:     weather.Timestamp.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     weather.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Temperature:   weather.Temperature,
		FeelsLike:     weather.FeelsLike(),
		Humidity:      weather.Humidity,
		Pressure:      weather.Pressure,
		CloudCover:    weather.CloudCover,
		WindSpeed:     weather.WindSpeed,
		WindDirection: weather.WindDirection(),
		WindDegrees:   weather.WindDir,
		Precipitation: weather.Precipitation,
		Symbol:        weather.Symbol,
		Description:   weather.Description,
		Units: JSONUnits{
			Temperature:   "celsius",
			WindSpeed:     "meters_per_second",
			Pressure:      "hectopascal",
			Precipitation: "millimeters",
			Humidity:      "percent",
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jw)
}

// FormatForecast formats forecast as JSON
func (f *JSONFormatter) FormatForecast(w io.Writer, forecast *models.Forecast, opts Options) error {
	jf := JSONForecast{
		Location: forecast.Location,
		Hours:    make([]JSONHourlyForecast, len(forecast.Hours)),
	}

	for i, hour := range forecast.Hours {
		jf.Hours[i] = JSONHourlyForecast{
			Time:          hour.Time.Format("2006-01-02T15:04:05Z"),
			Temperature:   hour.Temperature,
			FeelsLike:     hour.FeelsLike(),
			Humidity:      hour.Humidity,
			WindSpeed:     hour.WindSpeed,
			Precipitation: hour.Precipitation,
			Symbol:        hour.Symbol,
			Description:   hour.Description,
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jf)
}

// FormatDailySummary formats daily summary as JSON
func (f *JSONFormatter) FormatDailySummary(w io.Writer, summary *models.DailySummary, opts Options) error {
	js := JSONDailySummary{
		Location:           summary.Location,
		Date:               summary.Date.Format("2006-01-02"),
		TemperatureMin:     summary.TemperatureMin,
		TemperatureMax:     summary.TemperatureMax,
		TemperatureAvg:     summary.TemperatureAvg,
		PrecipitationTotal: summary.PrecipitationTotal,
		Units: JSONUnits{
			Temperature:   "celsius",
			WindSpeed:     "meters_per_second",
			Pressure:      "hectopascal",
			Precipitation: "millimeters",
			Humidity:      "percent",
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(js)
}

// FormatDailyForecast formats daily forecast as JSON
func (f *JSONFormatter) FormatDailyForecast(w io.Writer, dailyForecast *models.DailyForecast, opts Options) error {
	type JSONDailyForecastDay struct {
		Date               string  `json:"date"`
		TemperatureMin     float64 `json:"temperature_min"`
		TemperatureMax     float64 `json:"temperature_max"`
		TemperatureAvg     float64 `json:"temperature_avg"`
		PrecipitationTotal float64 `json:"precipitation_total"`
		Symbol             string  `json:"symbol"`
		WindSpeedMax       float64 `json:"wind_speed_max"`
	}

	type JSONDailyForecastOutput struct {
		Location *models.Location       `json:"location"`
		Days     []JSONDailyForecastDay `json:"days"`
		Units    JSONUnits              `json:"units"`
	}

	output := JSONDailyForecastOutput{
		Location: dailyForecast.Location,
		Days:     make([]JSONDailyForecastDay, len(dailyForecast.Days)),
		Units: JSONUnits{
			Temperature:   "celsius",
			WindSpeed:     "meters_per_second",
			Pressure:      "hectopascal",
			Precipitation: "millimeters",
			Humidity:      "percent",
		},
	}

	for i, day := range dailyForecast.Days {
		output.Days[i] = JSONDailyForecastDay{
			Date:               day.Date.Format("2006-01-02"),
			TemperatureMin:     day.TemperatureMin,
			TemperatureMax:     day.TemperatureMax,
			TemperatureAvg:     day.TemperatureAvg,
			PrecipitationTotal: day.PrecipitationTotal,
			Symbol:             day.Symbol,
			WindSpeedMax:       day.WindSpeedMax,
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
