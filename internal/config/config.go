package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kristofferrisa/sky-cli/internal/models"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	DefaultLocation string                     `yaml:"default_location" mapstructure:"default_location"`
	DefaultFormat   string                     `yaml:"default_format" mapstructure:"default_format"`
	NoColor         bool                       `yaml:"no_color" mapstructure:"no_color"`
	NoEmoji         bool                       `yaml:"no_emoji" mapstructure:"no_emoji"`
	Locations       map[string]*models.Location `yaml:"locations" mapstructure:"locations"`
}

// Load loads configuration from file or creates default config
func Load() (*Config, error) {
	// Set config file name and locations
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add config paths
	viper.AddConfigPath("$HOME/.sky")
	viper.AddConfigPath("$HOME/.config/sky")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("default_format", "full")
	viper.SetDefault("no_color", false)
	viper.SetDefault("no_emoji", false)
	viper.SetDefault("locations", map[string]*models.Location{
		"stavern": {
			Name:      "Stavern, Norway",
			Latitude:  59.0,
			Longitude: 10.0,
			Timezone:  "Europe/Oslo",
		},
	})
	viper.SetDefault("default_location", "stavern")

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configDir := filepath.Join(os.Getenv("HOME"), ".sky")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")

	viper.Set("default_location", c.DefaultLocation)
	viper.Set("default_format", c.DefaultFormat)
	viper.Set("no_color", c.NoColor)
	viper.Set("no_emoji", c.NoEmoji)
	viper.Set("locations", c.Locations)

	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetLocation retrieves a location by name
func (c *Config) GetLocation(name string) (*models.Location, error) {
	loc, ok := c.Locations[name]
	if !ok {
		return nil, fmt.Errorf("location '%s' not found in config", name)
	}
	return loc, nil
}

// GetDefaultLocation returns the default location
func (c *Config) GetDefaultLocation() (*models.Location, error) {
	if c.DefaultLocation == "" {
		return nil, fmt.Errorf("no default location configured")
	}
	return c.GetLocation(c.DefaultLocation)
}

// AddLocation adds or updates a location
func (c *Config) AddLocation(name string, loc *models.Location) error {
	if err := loc.Validate(); err != nil {
		return err
	}

	if c.Locations == nil {
		c.Locations = make(map[string]*models.Location)
	}

	c.Locations[name] = loc
	return nil
}

// RemoveLocation removes a location by name
func (c *Config) RemoveLocation(name string) error {
	if _, ok := c.Locations[name]; !ok {
		return fmt.Errorf("location '%s' not found", name)
	}

	delete(c.Locations, name)
	return nil
}
