package models

import "time"

// Weather represents current weather conditions at a specific location
type Weather struct {
	Location    *Location
	Timestamp   time.Time
	UpdatedAt   time.Time
	Temperature float64 // Celsius
	Humidity    float64 // Percentage (0-100)
	Pressure    float64 // hPa (hectopascal)
	CloudCover  float64 // Percentage (0-100)
	WindSpeed   float64 // m/s
	WindDir     float64 // degrees (0-360)
	Precipitation float64 // mm for next hour
	Symbol      string  // Weather symbol code
	Description string  // Human-readable description
}

// WindDirection returns a human-readable wind direction
func (w *Weather) WindDirection() string {
	degrees := int(w.WindDir)

	switch {
	case degrees >= 338 || degrees < 23:
		return "N (North)"
	case degrees >= 23 && degrees < 68:
		return "NE (Northeast)"
	case degrees >= 68 && degrees < 113:
		return "E (East)"
	case degrees >= 113 && degrees < 158:
		return "SE (Southeast)"
	case degrees >= 158 && degrees < 203:
		return "S (South)"
	case degrees >= 203 && degrees < 248:
		return "SW (Southwest)"
	case degrees >= 248 && degrees < 293:
		return "W (West)"
	case degrees >= 293 && degrees < 338:
		return "NW (Northwest)"
	default:
		return "Unknown"
	}
}

// WindDescription returns a description of wind strength
func (w *Weather) WindDescription() string {
	speed := int(w.WindSpeed)

	switch {
	case speed < 2:
		return "Calm"
	case speed < 6:
		return "Light breeze"
	case speed < 12:
		return "Moderate breeze"
	case speed < 20:
		return "Strong breeze"
	case speed < 29:
		return "Gale"
	default:
		return "Storm"
	}
}

// Forecast represents hourly weather forecast data
type Forecast struct {
	Location *Location
	Hours    []HourlyForecast
}

// HourlyForecast represents weather forecast for a specific hour
type HourlyForecast struct {
	Time          time.Time
	Temperature   float64
	Humidity      float64
	WindSpeed     float64
	Precipitation float64
	Symbol        string
	Description   string
}

// DailySummary represents aggregated weather data for a day
type DailySummary struct {
	Location           *Location
	Date               time.Time
	TemperatureMin     float64
	TemperatureMax     float64
	TemperatureAvg     float64
	PrecipitationTotal float64
	Symbol             string  // Most common symbol for the day
	WindSpeedMax       float64 // Maximum wind speed
}

// DailyForecast represents multi-day forecast
type DailyForecast struct {
	Location *Location
	Days     []DailySummary
}
