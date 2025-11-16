package models

import "fmt"

// Location represents a geographic location
type Location struct {
	Name      string  `yaml:"name" json:"name"`
	Latitude  float64 `yaml:"latitude" json:"latitude"`
	Longitude float64 `yaml:"longitude" json:"longitude"`
	Timezone  string  `yaml:"timezone,omitempty" json:"timezone,omitempty"`
}

// String returns a human-readable string representation
func (l *Location) String() string {
	if l.Name != "" {
		return fmt.Sprintf("%s (%.2f°N, %.2f°E)", l.Name, l.Latitude, l.Longitude)
	}
	return fmt.Sprintf("%.2f°N, %.2f°E", l.Latitude, l.Longitude)
}

// Validate checks if the location has valid coordinates
func (l *Location) Validate() error {
	if l.Latitude < -90 || l.Latitude > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", l.Latitude)
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", l.Longitude)
	}
	return nil
}
