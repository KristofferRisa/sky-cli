# Sky CLI ☁️

A beautiful, fast, and extensible command-line weather tool written in Go.

**Powered by MET Norway (Meteorologisk institutt)**

## Features

- **Feels Like Temperature**: Industry-standard Wind Chill and Heat Index calculations
- **Multiple Output Formats**: Full, JSON, Summary, and Markdown formats
- **Current Weather**: Get instant weather conditions for any location
- **Hourly Forecasts**: Dedicated forecast command with customizable hours
- **Daily/Weekly Forecasts**: Multi-day weather forecasts (up to 10 days)
- **Location Management**: Save and manage your favorite locations
- **Smart Caching**: File-based cache for 78x faster repeat queries
- **Rich Formatting**: Beautiful terminal output with colors and emojis
- **LLM-Friendly**: Structured data output perfect for AI processing
- **Flexible Input**: Use location names, coordinates, or defaults
- **Well Tested**: Comprehensive unit tests with high coverage
- **Fast & Reliable**: Compiled Go binary with minimal dependencies

## Installation

### Homebrew (macOS & Linux)

```bash
brew install kristofferrisa/sky/sky
```

### Manual Download (Linux / Other Platforms)

Download the latest release for your platform from [GitHub Releases](https://github.com/kristofferrisa/sky-cli/releases).

```bash
# Example for Linux (amd64)
curl -LO https://github.com/kristofferrisa/sky-cli/releases/latest/download/sky_[VERSION]_Linux_x86_64.tar.gz
tar -xzf sky_[VERSION]_Linux_x86_64.tar.gz
sudo mv sky /usr/local/bin/
```

### From Source

```bash
# Clone the repository
git clone https://github.com/kristofferrisa/sky-cli
cd sky-cli

# Build using Makefile (recommended)
make build

# Or build directly with Go
go build -o sky ./cmd/sky

# Install to system (optional)
sudo mv sky /usr/local/bin/

# Or install to GOPATH/bin
make install
```

## Quick Start

```bash
# Get current weather (default location)
sky current

# Get weather in different formats
sky current --format json      # JSON output
sky current --format summary   # One-line summary
sky current --format markdown  # Markdown format

# Get forecasts
sky forecast                   # 12-hour forecast
sky forecast --hours 24        # 24-hour forecast
sky daily                      # 7-day forecast
sky daily --days 10            # 10-day forecast

# Manage locations
sky locations add oslo --lat 59.9139 --lon 10.7522
sky locations list
sky current oslo
```

## Commands

### `sky current` - Current Weather

Get current weather conditions with optional forecast and daily summary.

```bash
# Basic usage
sky current                          # Use default location
sky current stavern                  # Use saved location
sky current --lat 59.0 --lon 10.0   # Use coordinates

# With forecast and summary
sky current --forecast               # Include 12-hour forecast
sky current --summary                # Include daily summary
sky current --forecast --summary     # Include both
sky current --forecast --hours 24    # Custom forecast hours

# Different formats
sky current --format json            # JSON output
sky current --format summary         # Brief one-line summary
sky current --format markdown        # Markdown format
sky current --format full            # Detailed output (default)
```

**Flags:**

- `--format, -f` - Output format (full, json, summary, markdown)
- `--location, -l` - Location name from config
- `--lat` - Latitude
- `--lon` - Longitude
- `--forecast` - Include hourly forecast
- `--summary` - Include daily summary
- `--hours` - Number of hours for forecast (default: 12)

### `sky forecast` - Weather Forecast

Get hourly weather forecast for a location.

```bash
# Basic usage
sky forecast                         # 12-hour forecast (default location)
sky forecast stavern                 # Forecast for saved location
sky forecast --lat 59.0 --lon 10.0  # Forecast for coordinates

# Custom hours
sky forecast --hours 6               # 6-hour forecast
sky forecast --hours 24              # 24-hour forecast

# Different formats
sky forecast --format json           # JSON output
sky forecast --format summary        # Brief summary
sky forecast --format markdown       # Markdown table
```

**Flags:**

- `--format, -f` - Output format (full, json, summary, markdown)
- `--location, -l` - Location name from config
- `--lat` - Latitude
- `--lon` - Longitude
- `--hours` - Number of hours for forecast (default: 12)

### `sky daily` - Daily Weather Forecast

Get daily weather forecast for multiple days.

```bash
# Basic usage
sky daily                       # 7-day forecast (default location)
sky daily stavern               # 7-day forecast for saved location
sky daily --lat 59.0 --lon 10.0 # Forecast for coordinates

# Custom days
sky daily --days 3              # 3-day forecast
sky daily --days 10             # 10-day forecast

# Different formats
sky daily --format json         # JSON output
sky daily --format summary      # Brief summary
sky daily --format markdown     # Markdown table
```

**Flags:**

- `--format, -f` - Output format (full, json, summary, markdown)
- `--location, -l` - Location name from config
- `--lat` - Latitude
- `--lon` - Longitude
- `--days` - Number of days for forecast (default: 7)

### `sky locations` - Location Management

Manage saved locations in your configuration.

```bash
# List all saved locations
sky locations list

# Add a new location
sky locations add bergen --lat 60.3913 --lon 5.3221 --timezone "Europe/Oslo"

# Remove a location
sky locations remove bergen

# Set default location
sky locations set-default oslo
```

**Subcommands:**

- `list` - List all saved locations
- `add <name>` - Add a new location (requires --lat and --lon)
- `remove <name>` - Remove a saved location
- `set-default <name>` - Set the default location

### Global Flags

Available on all commands:

- `--no-color` - Disable colored output
- `--no-emoji` - Disable emoji symbols
- `--help, -h` - Show help for any command

## Output Formats

Sky CLI supports four output formats for maximum flexibility:

### Full Format (default)

Beautiful, detailed terminal output with colors and emojis.

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CURRENT WEATHER - Stavern, Norway (59.00°N, 10.00°E)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Location:     Stavern, Norway (59.00°N, 10.00°E)
Time:         Sunday, November 16, 2025 at 20:00

Conditions:   ☀️ Clear sky
Temperature:  1.1°C (feels like -2.8°C)
Humidity:     70%
Wind:         2.6 m/s from NW (Northwest)
...
```

### JSON Format

Machine-readable JSON output, perfect for scripting and integration.

```bash
sky current --format json | jq '.temperature'
```

```json
{
  "location": {
    "name": "Stavern, Norway",
    "latitude": 59,
    "longitude": 10
  },
  "temperature": 1.1,
  "feels_like": -2.8,
  "humidity": 70,
  "wind_speed": 2.6,
  "conditions": "Clear sky"
}
```

### Summary Format

Brief one-line summaries, great for status bars and quick checks.

```bash
sky current --format summary
# Output: Stavern, Norway 20:00: ☀️ Clear sky 1.1°C (feels like -2.8°C), Wind: 2.6 m/s NW, Humidity: 70%
```

### Markdown Format

Documentation-friendly markdown output, perfect for reports and sharing.

```markdown
# Weather for Stavern, Norway

## Current Conditions

- **Conditions:** ☀️ Clear sky
- **Temperature:** 1.1°C (feels like -2.8°C)
- **Humidity:** 70%
```

## Feels Like Temperature

Sky CLI calculates "feels like" temperature using industry-standard formulas that match what you see on professional weather services.

### How It Works

The app automatically applies the appropriate formula based on conditions:

**Wind Chill** (Cold & Windy)
- Applied when: Temperature ≤ 10°C AND wind speed > 4.8 km/h
- Formula: North American/UK standard (adopted 2001)
- Example: `0°C with 15 km/h wind → feels like -5°C`

**Heat Index** (Hot & Humid)
- Applied when: Temperature ≥ 27°C AND humidity > 40%
- Formula: US National Weather Service Rothfusz regression
- Example: `32°C with 80% humidity → feels like 41°C`

**Moderate Conditions**
- When neither formula applies, shows actual temperature
- Example: `20°C → feels like 20°C`

### Why This Matters

For Norway and other cold climates, Wind Chill is especially important:
- Shows how cold it actually *feels* on your skin
- Helps you dress appropriately for the weather
- Matches what you see on Yr.no and other professional services

### Available in All Formats

- **Full**: `Temperature: 0.2°C (feels like -3.0°C)`
- **JSON**: `"feels_like": -3.0` field
- **Summary**: Inline with temperature
- **Markdown**: Listed with temperature

## Configuration

Sky CLI uses a configuration file located at `~/.sky/config.yaml` (or `~/.config/sky/config.yaml`).

### Example Configuration

```yaml
# Default location to use when no location is specified
default_location: stavern

# Default output format (full, json, summary, markdown)
default_format: full

# Disable colors/emojis globally
no_color: false
no_emoji: false

# Cache configuration
cache:
  enabled: true
  directory: ~/.sky/cache
  ttl_minutes: 10

# Saved locations
locations:
  stavern:
    name: "Stavern, Norway"
    latitude: 59.0
    longitude: 10.0
    timezone: "Europe/Oslo"

  oslo:
    name: "Oslo, Norway"
    latitude: 59.9139
    longitude: 10.7522
    timezone: "Europe/Oslo"

  bergen:
    name: "Bergen, Norway"
    latitude: 60.3913
    longitude: 5.3221
    timezone: "Europe/Oslo"
```

### Cache Configuration

Sky CLI caches weather data to reduce API calls and improve performance.

- **Default TTL**: 10 minutes
- **Cache Location**: `~/.sky/cache/`
- **Performance**: 78x faster on cached requests!
- **Automatic**: No user action needed

To disable caching:

```yaml
cache:
  enabled: false
```

## Usage Examples

### Quick Weather Check

```bash
# Brief summary for your location
sky current --format summary

# JSON for scripting
TEMP=$(sky current --format json | jq -r '.temperature')
echo "Current temperature: ${TEMP}°C"
```

### Planning Your Day

```bash
# Full report with forecast
sky current --forecast --summary

# 24-hour forecast
sky forecast --hours 24

# Week ahead
sky daily --days 7
```

### Managing Locations

```bash
# Add your home and work locations
sky locations add home --lat 59.0 --lon 10.0
sky locations add work --lat 59.9 --lon 10.8

# Set default
sky locations set-default home

# Quick check for work
sky current work
```

### Markdown Reports

```bash
# Generate markdown weather report
sky current --format markdown --forecast > weather-report.md
```

## Architecture

Sky CLI is built with a clean, modular architecture:

- **Provider Interface**: Pluggable weather API providers (currently MET Norway)
- **Formatter Interface**: Multiple output formats (full, JSON, summary, markdown)
- **Cache Layer**: File-based cache with TTL management
- **Configuration System**: YAML-based config with saved locations
- **CLI Framework**: Cobra for robust command-line interface

### Project Structure

```
sky-cli/
├── cmd/sky/              # CLI entry point and commands
│   ├── main.go
│   ├── root.go
│   ├── current.go       # Current weather command
│   ├── forecast.go      # Hourly forecast command
│   ├── daily.go         # Daily forecast command
│   └── locations.go     # Location management
├── internal/
│   ├── api/
│   │   ├── client.go         # Weather client interface
│   │   └── met/              # MET Norway provider
│   │       ├── client.go
│   │       ├── cached_client.go
│   │       └── models.go
│   ├── cache/                # Caching layer
│   │   ├── cache.go
│   │   └── file.go
│   ├── config/               # Configuration
│   │   └── config.go
│   ├── formatter/            # Output formatters
│   │   ├── formatter.go
│   │   ├── full.go
│   │   ├── json.go
│   │   ├── summary.go
│   │   ├── markdown.go
│   │   └── factory.go
│   ├── models/               # Data models
│   │   ├── weather.go
│   │   └── location.go
│   └── ui/                   # UI helpers
│       ├── colors.go
│       └── symbols.go
├── go.mod
├── PROGRESS.md               # Development tracking
└── README.md
```

## Development

### Quick Start with Makefile

The project includes a comprehensive Makefile for common development tasks:

```bash
# Build the binary
make build

# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Format code and run linting
make lint

# Clean up dependencies
make tidy

# Run everything (tidy, lint, test, build)
make all

# See all available commands
make help
```

### Development Commands

```bash
# Quick development iteration (no build step)
make dev

# Build and test commands
make run-current    # Build and run current weather
make run-forecast   # Build and run forecast
make run-daily      # Build and run daily forecast

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install
```

### Direct Go Commands

You can also use Go commands directly:

```bash
# Build
go build -o sky ./cmd/sky

# Run tests
go test ./...
go test ./... -cover
go test ./... -v

# Run locally
go run ./cmd/sky current

# Format and lint
gofmt -w .
go vet ./...
go mod tidy
```

**Test Coverage:**

- Models: 97.5% coverage (includes Wind Chill & Heat Index calculations)
- Cache: 56.7% coverage
- 51 test cases, all passing

### Release Testing

```bash
# Test GoReleaser configuration
make release-check

# Create a local snapshot release
make release-snapshot
```

## Roadmap

### Phase 1: Core Foundation ✅ COMPLETED

- [x] Basic CLI structure
- [x] MET Norway provider
- [x] Full formatter
- [x] Current weather command
- [x] Basic configuration

### Phase 2: Enhanced Features ✅ COMPLETED

- [x] JSON, Summary, Markdown formatters
- [x] Forecast command (hourly)
- [x] Cache layer (78x performance improvement!)
- [x] Location management commands
- [x] Format selection via --format flag

### Phase 3: Extensibility ✅ COMPLETED

- [x] Daily/weekly weather forecasts (up to 10 days)
- [x] Feels Like Temperature (Wind Chill & Heat Index)
- [x] Unit tests (97.5% coverage for models, 56.7% for cache)
- [ ] Additional weather providers (OpenWeather, Weather.gov) - deferred
- [ ] Unit conversion (metric/imperial) - deferred
- [ ] Weather alerts and warnings - deferred

### Phase 4: CI/CD & Distribution ✅ IN PROGRESS

- [x] CI/CD pipeline (GitHub Actions)
- [x] Automated testing workflow
- [x] Cross-platform builds (GoReleaser)
- [x] Automated releases with GitHub
- [x] Version command
- [ ] Homebrew formula (planned)
- [ ] Package managers (apt, snap) (planned)

## Performance

Sky CLI is designed for speed:

- **Compiled Binary**: Fast startup time
- **Smart Caching**: 78x faster on cached requests
  - First request: ~624ms (API call)
  - Cached request: ~8ms (from disk)
- **Efficient API Usage**: Only fetches what you need

## Credits

- **Weather Data**: [MET Norway](https://www.met.no/) (Meteorologisk institutt)
- **Author**: Kristoffer Risa

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

Special thanks to MET Norway for providing free, high-quality weather data through their excellent API.
