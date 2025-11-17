package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/kristofferrisa/sky-cli/internal/models"
	"github.com/spf13/cobra"
)

// locationsCmd represents the locations command
var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "Manage saved locations",
	Long:  `Manage saved locations in the configuration file.`,
}

// listLocationsCmd lists all saved locations
var listLocationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved locations",
	Long:  `Display all saved locations with their coordinates.`,
	RunE:  runListLocations,
}

// addLocationCmd adds a new location
var addLocationCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new location",
	Long:  `Add a new location to the configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAddLocation,
}

// removeLocationCmd removes a location
var removeLocationCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a saved location",
	Long:    `Remove a location from the configuration file.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runRemoveLocation,
}

// setDefaultLocationCmd sets the default location
var setDefaultLocationCmd = &cobra.Command{
	Use:   "set-default <name>",
	Short: "Set default location",
	Long:  `Set the default location to use when no location is specified.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSetDefaultLocation,
}

var (
	addLat      float64
	addLon      float64
	addTimezone string
)

func init() {
	// Add subcommands
	locationsCmd.AddCommand(listLocationsCmd)
	locationsCmd.AddCommand(addLocationCmd)
	locationsCmd.AddCommand(removeLocationCmd)
	locationsCmd.AddCommand(setDefaultLocationCmd)

	// Add flags for add command
	addLocationCmd.Flags().Float64Var(&addLat, "lat", 0, "Latitude (required)")
	addLocationCmd.Flags().Float64Var(&addLon, "lon", 0, "Longitude (required)")
	addLocationCmd.Flags().StringVar(&addTimezone, "timezone", "", "Timezone (optional)")
	addLocationCmd.MarkFlagRequired("lat")
	addLocationCmd.MarkFlagRequired("lon")

	rootCmd.AddCommand(locationsCmd)
}

func runListLocations(cmd *cobra.Command, args []string) error {
	if len(cfg.Locations) == 0 {
		fmt.Println("No saved locations")
		return nil
	}

	// Sort locations by name
	names := make([]string, 0, len(cfg.Locations))
	for name := range cfg.Locations {
		names = append(names, name)
	}
	sort.Strings(names)

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tLOCATION\tCOORDINATES\tDEFAULT")
	fmt.Fprintln(w, "────\t────────\t───────────\t───────")

	for _, name := range names {
		loc := cfg.Locations[name]
		isDefault := ""
		if name == cfg.DefaultLocation {
			isDefault = "✓"
		}

		locationName := loc.Name
		if locationName == "" {
			locationName = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%.4f°N, %.4f°E\t%s\n",
			name,
			locationName,
			loc.Latitude,
			loc.Longitude,
			isDefault,
		)
	}

	w.Flush()
	return nil
}

func runAddLocation(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Check if location already exists
	if _, exists := cfg.Locations[name]; exists {
		return fmt.Errorf("location '%s' already exists (use 'locations remove' first to replace)", name)
	}

	// Create location
	loc := &models.Location{
		Name:      name,
		Latitude:  addLat,
		Longitude: addLon,
		Timezone:  addTimezone,
	}

	// Validate
	if err := loc.Validate(); err != nil {
		return err
	}

	// Add to config
	if err := cfg.AddLocation(name, loc); err != nil {
		return err
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Location '%s' added successfully\n", name)
	fmt.Printf("  %s\n", loc.String())

	// Suggest setting as default if no default exists
	if cfg.DefaultLocation == "" {
		fmt.Printf("\nTip: Set as default with: sky locations set-default %s\n", name)
	}

	return nil
}

func runRemoveLocation(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Check if location exists
	if _, exists := cfg.Locations[name]; !exists {
		return fmt.Errorf("location '%s' not found", name)
	}

	// Remove from config
	if err := cfg.RemoveLocation(name); err != nil {
		return err
	}

	// Clear default if this was the default
	if cfg.DefaultLocation == name {
		cfg.DefaultLocation = ""
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Location '%s' removed successfully\n", name)

	return nil
}

func runSetDefaultLocation(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Check if location exists
	if _, exists := cfg.Locations[name]; !exists {
		return fmt.Errorf("location '%s' not found (use 'locations add' first)", name)
	}

	// Set as default
	cfg.DefaultLocation = name

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Default location set to '%s'\n", name)
	fmt.Printf("  %s\n", cfg.Locations[name].String())

	return nil
}
