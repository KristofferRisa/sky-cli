package models

import (
	"testing"
)

func TestWeatherWindDirection(t *testing.T) {
	tests := []struct {
		name     string
		degrees  float64
		expected string
	}{
		{"North", 0, "N (North)"},
		{"North upper bound", 22, "N (North)"},
		{"Northeast", 45, "NE (Northeast)"},
		{"East", 90, "E (East)"},
		{"Southeast", 135, "SE (Southeast)"},
		{"South", 180, "S (South)"},
		{"Southwest", 225, "SW (Southwest)"},
		{"West", 270, "W (West)"},
		{"Northwest", 315, "NW (Northwest)"},
		{"North wraparound", 350, "N (North)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weather := &Weather{WindDir: tt.degrees}
			result := weather.WindDirection()
			if result != tt.expected {
				t.Errorf("WindDirection() for %f degrees = %s; want %s", tt.degrees, result, tt.expected)
			}
		})
	}
}

func TestWeatherWindDescription(t *testing.T) {
	tests := []struct {
		name     string
		speed    float64
		expected string
	}{
		{"Calm", 1, "Calm"},
		{"Light breeze", 4, "Light breeze"},
		{"Moderate breeze", 10, "Moderate breeze"},
		{"Strong breeze", 15, "Strong breeze"},
		{"Gale", 25, "Gale"},
		{"Storm", 30, "Storm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weather := &Weather{WindSpeed: tt.speed}
			result := weather.WindDescription()
			if result != tt.expected {
				t.Errorf("WindDescription() for %f m/s = %s; want %s", tt.speed, result, tt.expected)
			}
		})
	}
}

func TestCalculateApparentTemperature(t *testing.T) {
	tests := []struct {
		name        string
		temperature float64
		humidity    float64
		windSpeed   float64
		expected    float64
		tolerance   float64
		description string
	}{
		// WIND CHILL ZONE (temp ≤ 10°C AND wind > 4.8 km/h)
		{
			name:        "Wind Chill: Cold with moderate wind",
			temperature: -2.0,
			humidity:    70.0,
			windSpeed:   4.0, // 14.4 km/h
			expected:    -6.76,
			tolerance:   0.1,
			description: "Typical Norwegian winter day",
		},
		{
			name:        "Wind Chill: Extreme cold with high wind",
			temperature: -10.0,
			humidity:    60.0,
			windSpeed:   10.0, // 36 km/h
			expected:    -20.30,
			tolerance:   0.1,
			description: "Dangerous wind chill conditions",
		},
		{
			name:        "Wind Chill: Freezing with light wind",
			temperature: 0.0,
			humidity:    50.0,
			windSpeed:   5.0, // 18 km/h
			expected:    -4.94,
			tolerance:   0.1,
			description: "Just below freezing with wind",
		},
		{
			name:        "Wind Chill: Boundary case at 10°C",
			temperature: 10.0,
			humidity:    70.0,
			windSpeed:   15.0, // 54 km/h - high wind
			expected:    5.32,
			tolerance:   0.1,
			description: "At upper temperature threshold",
		},
		{
			name:        "Wind Chill: Cold with strong wind",
			temperature: 5.0,
			humidity:    60.0,
			windSpeed:   8.0, // 28.8 km/h
			expected:    0.16,
			tolerance:   0.1,
			description: "Moderate cold with noticeable wind",
		},

		// HEAT INDEX ZONE (temp ≥ 27°C AND humidity > 40%)
		{
			name:        "Heat Index: Hot with high humidity",
			temperature: 30.0,
			humidity:    80.0,
			windSpeed:   2.0,
			expected:    37.67,
			tolerance:   0.1,
			description: "Oppressive summer conditions",
		},
		{
			name:        "Heat Index: Extreme heat with high humidity",
			temperature: 35.0,
			humidity:    90.0,
			windSpeed:   1.0,
			expected:    63.67,
			tolerance:   0.1,
			description: "Dangerous heat conditions",
		},
		{
			name:        "Heat Index: Boundary case at 27°C",
			temperature: 27.0,
			humidity:    50.0,
			windSpeed:   3.0,
			expected:    27.42,
			tolerance:   0.1,
			description: "At lower temperature threshold",
		},

		// MODERATE ZONE (no wind chill or heat index applied)
		{
			name:        "Moderate: Comfortable spring day",
			temperature: 20.0,
			humidity:    50.0,
			windSpeed:   5.0,
			expected:    20.0, // Returns actual temp
			tolerance:   0.01,
			description: "No adjustment needed",
		},
		{
			name:        "Moderate: Mild with no wind",
			temperature: 15.0,
			humidity:    60.0,
			windSpeed:   0.0,
			expected:    15.0, // Returns actual temp
			tolerance:   0.01,
			description: "Calm conditions",
		},
		{
			name:        "Moderate: Warm but low humidity",
			temperature: 25.0,
			humidity:    20.0, // Too low for heat index
			windSpeed:   3.0,
			expected:    25.0, // Returns actual temp
			tolerance:   0.01,
			description: "Dry heat - no heat index",
		},
		{
			name:        "Moderate: Perfect comfortable weather",
			temperature: 22.0,
			humidity:    45.0,
			windSpeed:   3.0,
			expected:    22.0, // Returns actual temp
			tolerance:   0.01,
			description: "Ideal conditions",
		},
		{
			name:        "Moderate: Cool but insufficient wind",
			temperature: 8.0,
			humidity:    50.0,
			windSpeed:   1.0, // 3.6 km/h - below threshold
			expected:    8.0, // Returns actual temp (wind too low)
			tolerance:   0.01,
			description: "Wind below 4.8 km/h threshold",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateApparentTemperature(tt.temperature, tt.humidity, tt.windSpeed)
			diff := result - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("calculateApparentTemperature(%f, %f, %f) = %f; want %f (±%f)",
					tt.temperature, tt.humidity, tt.windSpeed, result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestWeatherFeelsLike(t *testing.T) {
	tests := []struct {
		name        string
		weather     *Weather
		expectLower bool // Expect feels like to be lower than actual temp
	}{
		{
			name: "Cold with wind - feels colder (Wind Chill)",
			weather: &Weather{
				Temperature: -2.0,
				Humidity:    70.0,
				WindSpeed:   5.0, // 18 km/h
			},
			expectLower: true,
		},
		{
			name: "Hot with high humidity - feels hotter (Heat Index)",
			weather: &Weather{
				Temperature: 30.0,
				Humidity:    85.0,
				WindSpeed:   2.0,
			},
			expectLower: false,
		},
		{
			name: "Moderate temperature - no change",
			weather: &Weather{
				Temperature: 15.0,
				Humidity:    60.0,
				WindSpeed:   2.0, // Not enough to trigger wind chill
			},
			expectLower: false, // Will return actual temp (no change)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feelsLike := tt.weather.FeelsLike()

			// Verify the method returns a reasonable value
			if feelsLike < -50 || feelsLike > 70 {
				t.Errorf("FeelsLike() = %f is outside reasonable range", feelsLike)
			}

			// Verify wind chill / heat index direction
			if tt.expectLower && feelsLike >= tt.weather.Temperature {
				t.Errorf("Expected feels like (%f) to be lower than temperature (%f)",
					feelsLike, tt.weather.Temperature)
			}
			if !tt.expectLower && tt.weather.Temperature >= 27.0 && feelsLike <= tt.weather.Temperature {
				// For heat index cases, should feel hotter
				t.Errorf("Expected feels like (%f) to be higher than temperature (%f)",
					feelsLike, tt.weather.Temperature)
			}
		})
	}
}

func TestHourlyForecastFeelsLike(t *testing.T) {
	tests := []struct {
		name     string
		forecast *HourlyForecast
		minDiff  float64 // Minimum expected difference from actual temp
	}{
		{
			name: "Wind Chill: Strong wind creates significant difference",
			forecast: &HourlyForecast{
				Temperature: 0.0, // Cold enough for wind chill
				Humidity:    60.0,
				WindSpeed:   12.0, // 43.2 km/h - strong wind
			},
			minDiff: 7.5, // Should feel at least 7.5°C colder
		},
		{
			name: "Moderate zone: No change expected",
			forecast: &HourlyForecast{
				Temperature: 20.0,
				Humidity:    50.0,
				WindSpeed:   1.0,
			},
			minDiff: 0.0, // No difference in moderate zone
		},
		{
			name: "Heat Index: High humidity creates difference",
			forecast: &HourlyForecast{
				Temperature: 32.0,
				Humidity:    80.0,
				WindSpeed:   2.0,
			},
			minDiff: 5.0, // Should feel at least 5°C hotter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feelsLike := tt.forecast.FeelsLike()

			// Verify reasonable range
			if feelsLike < -50 || feelsLike > 60 {
				t.Errorf("FeelsLike() = %f is outside reasonable range", feelsLike)
			}

			// Verify minimum difference
			diff := tt.forecast.Temperature - feelsLike
			if diff < 0 {
				diff = -diff
			}
			if diff < tt.minDiff {
				t.Errorf("FeelsLike() difference from temperature (%f) is too small; expected at least %f",
					diff, tt.minDiff)
			}
		})
	}
}

func TestFeelsLikeConsistency(t *testing.T) {
	// Test that Weather and HourlyForecast produce the same results
	// for the same conditions
	temp := 15.0
	humidity := 65.0
	windSpeed := 7.0

	weather := &Weather{
		Temperature: temp,
		Humidity:    humidity,
		WindSpeed:   windSpeed,
	}

	forecast := &HourlyForecast{
		Temperature: temp,
		Humidity:    humidity,
		WindSpeed:   windSpeed,
	}

	weatherFeels := weather.FeelsLike()
	forecastFeels := forecast.FeelsLike()

	if weatherFeels != forecastFeels {
		t.Errorf("Weather.FeelsLike() = %f; HourlyForecast.FeelsLike() = %f; expected same value",
			weatherFeels, forecastFeels)
	}
}
