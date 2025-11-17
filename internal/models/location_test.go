package models

import (
	"testing"
)

func TestLocationValidate(t *testing.T) {
	tests := []struct {
		name      string
		location  *Location
		shouldErr bool
	}{
		{
			name: "Valid location",
			location: &Location{
				Name:      "Oslo",
				Latitude:  59.9139,
				Longitude: 10.7522,
			},
			shouldErr: false,
		},
		{
			name: "Latitude too high",
			location: &Location{
				Latitude:  91,
				Longitude: 10,
			},
			shouldErr: true,
		},
		{
			name: "Latitude too low",
			location: &Location{
				Latitude:  -91,
				Longitude: 10,
			},
			shouldErr: true,
		},
		{
			name: "Longitude too high",
			location: &Location{
				Latitude:  59,
				Longitude: 181,
			},
			shouldErr: true,
		},
		{
			name: "Longitude too low",
			location: &Location{
				Latitude:  59,
				Longitude: -181,
			},
			shouldErr: true,
		},
		{
			name: "Valid edge case - North pole",
			location: &Location{
				Latitude:  90,
				Longitude: 0,
			},
			shouldErr: false,
		},
		{
			name: "Valid edge case - South pole",
			location: &Location{
				Latitude:  -90,
				Longitude: 0,
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.location.Validate()
			if tt.shouldErr && err == nil {
				t.Errorf("Validate() expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestLocationString(t *testing.T) {
	tests := []struct {
		name     string
		location *Location
		expected string
	}{
		{
			name: "With name",
			location: &Location{
				Name:      "Oslo, Norway",
				Latitude:  59.9139,
				Longitude: 10.7522,
			},
			expected: "Oslo, Norway (59.91°N, 10.75°E)",
		},
		{
			name: "Without name",
			location: &Location{
				Latitude:  59.0,
				Longitude: 10.0,
			},
			expected: "59.00°N, 10.00°E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.location.String()
			if result != tt.expected {
				t.Errorf("String() = %s; want %s", result, tt.expected)
			}
		})
	}
}
