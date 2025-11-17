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
