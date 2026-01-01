package models

import (
	"math"
	"time"
)

// Weather represents current weather conditions at a specific location
type Weather struct {
	Location      *Location
	Timestamp     time.Time
	UpdatedAt     time.Time
	Temperature   float64 // Celsius
	Humidity      float64 // Percentage (0-100)
	Pressure      float64 // hPa (hectopascal)
	CloudCover    float64 // Percentage (0-100)
	WindSpeed     float64 // m/s
	WindDir       float64 // degrees (0-360)
	Precipitation float64 // mm for next hour
	Symbol        string  // Weather symbol code
	Description   string  // Human-readable description
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

// FeelsLike calculates the apparent temperature (feels like)
// using temperature, humidity, and wind speed
func (w *Weather) FeelsLike() float64 {
	return calculateApparentTemperature(w.Temperature, w.Humidity, w.WindSpeed)
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

// FeelsLike calculates the apparent temperature (feels like)
// using temperature, humidity, and wind speed
func (h *HourlyForecast) FeelsLike() float64 {
	return calculateApparentTemperature(h.Temperature, h.Humidity, h.WindSpeed)
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

// calculateApparentTemperature calculates the "feels like" temperature
// using industry-standard formulas: Wind Chill for cold conditions,
// Heat Index for hot/humid conditions, or actual temperature otherwise.
//
// Parameters:
//   - temperature: Air temperature in °C
//   - humidity: Relative humidity as percentage (0-100)
//   - windSpeed: Wind speed in m/s
//
// Returns: Apparent temperature in °C
func calculateApparentTemperature(temperature, humidity, windSpeed float64) float64 {
	// Convert wind speed from m/s to km/h for Wind Chill formula
	windKmh := windSpeed * 3.6

	// Cold conditions: Use Wind Chill formula
	// Applied when temp ≤ 10°C and wind > 4.8 km/h (> 1.33 m/s)
	if temperature <= 10.0 && windKmh > 4.8 {
		return calculateWindChill(temperature, windKmh)
	}

	// Hot/humid conditions: Use Heat Index formula
	// Applied when temp ≥ 27°C and humidity > 40%
	if temperature >= 27.0 && humidity > 40.0 {
		return calculateHeatIndex(temperature, humidity)
	}

	// Moderate conditions: Return actual temperature
	// No significant wind chill or heat index effect
	return temperature
}

// calculateWindChill calculates wind chill temperature using the
// North American and UK standard formula (adopted 2001).
//
// Formula: WC = 13.12 + 0.6215×T - 11.37×V^0.16 + 0.3965×T×V^0.16
//
// Parameters:
//   - temperature: Air temperature in °C
//   - windKmh: Wind speed in km/h
//
// Returns: Wind chill temperature in °C
func calculateWindChill(temperature, windKmh float64) float64 {
	// Standard Wind Chill formula
	// WC = 13.12 + 0.6215×T - 11.37×V^0.16 + 0.3965×T×V^0.16
	wc := 13.12 +
		0.6215*temperature -
		11.37*math.Pow(windKmh, 0.16) +
		0.3965*temperature*math.Pow(windKmh, 0.16)

	return wc
}

// calculateHeatIndex calculates heat index using the Rothfusz regression
// formula used by the US National Weather Service.
//
// This is a simplified version that works well for typical conditions.
// For extreme conditions, the full formula with adjustments should be used.
//
// Parameters:
//   - temperature: Air temperature in °C
//   - humidity: Relative humidity as percentage (0-100)
//
// Returns: Heat index in °C
func calculateHeatIndex(temperature, humidity float64) float64 {
	// Convert to Fahrenheit for the standard formula
	tempF := temperature*9.0/5.0 + 32.0

	// Rothfusz regression (simplified Steadman formula)
	// Used by US National Weather Service
	hi := -42.379 +
		2.04901523*tempF +
		10.14333127*humidity -
		0.22475541*tempF*humidity -
		6.83783e-3*tempF*tempF -
		5.481717e-2*humidity*humidity +
		1.22874e-3*tempF*tempF*humidity +
		8.5282e-4*tempF*humidity*humidity -
		1.99e-6*tempF*tempF*humidity*humidity

	// Convert back to Celsius
	return (hi - 32.0) * 5.0 / 9.0
}
