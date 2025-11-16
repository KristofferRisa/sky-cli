# Sky CLI ☁️

A beautiful, fast, and extensible command-line weather tool written in Go.

**Powered by MET Norway (Meteorologisk institutt)**

## Features

- **Current Weather**: Get instant weather conditions for any location
- **Rich Formatting**: Beautiful terminal output with colors and emojis
- **LLM-Friendly**: Structured data output perfect for AI processing
- **Saved Locations**: Save your favorite locations for quick access
- **Flexible Input**: Use location names, coordinates, or defaults
- **Hourly Forecasts**: View upcoming weather (up to 12 hours)
- **Daily Summaries**: Get temperature ranges and precipitation totals
- **Fast & Reliable**: Compiled Go binary with minimal dependencies

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/kristofferrisa/sky-cli
cd sky-cli

# Build
go build -o sky ./cmd/sky

# Install (optional)
sudo mv sky /usr/local/bin/
```

## Quick Start

```bash
# Get current weather for default location (Stavern, Norway)
sky current

# Get weather with 12-hour forecast and daily summary
sky current --forecast --summary

# Get weather for specific coordinates (Oslo)
sky current --lat 59.9139 --lon 10.7522

# Disable colors and emojis (for scripting)
sky current --no-color --no-emoji
```

## Usage

### Current Weather

```bash
# Use default location
sky current

# Use saved location
sky current stavern

# Use coordinates
sky current --lat 59.0 --lon 10.0

# Include hourly forecast (12 hours)
sky current --forecast

# Include daily summary
sky current --summary

# Include both
sky current --forecast --summary

# Custom forecast hours
sky current --forecast --hours 24
```

### Global Flags

- `--no-color` - Disable colored output
- `--no-emoji` - Disable emoji symbols

## Configuration

Sky CLI uses a configuration file located at `~/.sky/config.yaml` (or `~/.config/sky/config.yaml`).

### Example Configuration

```yaml
default_location: stavern
default_format: full
no_color: false
no_emoji: false

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

  home:
    name: "Home"
    latitude: 59.0
    longitude: 10.0
    timezone: "Europe/Oslo"
```

### Adding Locations

Edit `~/.sky/config.yaml` and add your locations under the `locations` key.

## Output Examples

### Basic Current Weather

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CURRENT WEATHER - Stavern, Norway (59.00°N, 10.00°E)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
API: MET Norway (Meteorologisk institutt)
Coordinates: 59.00°N, 10.00°E
Request time: 2025-11-16 21:59:26

Location:     Stavern, Norway (59.00°N, 10.00°E)
Time:         Sunday, November 16, 2025 at 20:00
Updated:      Sunday, November 16, 2025 at 20:27

Conditions:   ☀️ Clear sky
Temperature:  1.1°C
Humidity:     70%
Cloud Cover:  1%
Precipitation: 0.0 mm (next hour)

Wind:         2.6 m/s from NW (Northwest)
Wind Status:  Light breeze
Pressure:     1000.4 hPa
```

### With Forecast and Summary

Includes:
- Current conditions
- 12-hour hourly forecast table
- Daily temperature ranges
- Total precipitation
- Weather warnings (if applicable)
- LLM-friendly structured data

## Architecture

Sky CLI is built with a clean, modular architecture:

- **Provider Interface**: Pluggable weather API providers (currently MET Norway)
- **Formatter Interface**: Multiple output formats (currently "full")
- **Configuration System**: YAML-based config with saved locations
- **Models**: Clean data structures for weather, forecasts, locations
- **CLI Framework**: Cobra for robust command-line interface

See [PROGRESS.md](PROGRESS.md) for detailed architecture and development roadmap.

## Development

### Project Structure

```
sky-cli/
├── cmd/sky/              # CLI entry point and commands
├── internal/
│   ├── api/met/         # MET Norway API client
│   ├── config/          # Configuration management
│   ├── formatter/       # Output formatters
│   ├── models/          # Data models
│   └── ui/              # UI helpers (colors, symbols)
├── go.mod
├── PROGRESS.md          # Development progress tracking
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o sky ./cmd/sky
```

## Roadmap

### Phase 1: Core Foundation ✅ COMPLETED
- [x] Basic CLI structure
- [x] MET Norway provider
- [x] Full formatter
- [x] Current weather command
- [x] Basic configuration

### Phase 2: Enhanced Features (Planned)
- [ ] JSON, Summary, Markdown formatters
- [ ] Forecast commands (hourly, daily)
- [ ] Cache layer
- [ ] Location management commands
- [ ] Unit conversion (metric/imperial)

### Phase 3: Extensibility (Planned)
- [ ] Additional weather providers (OpenWeather, Weather.gov)
- [ ] Weather alerts and warnings
- [ ] Historical data
- [ ] Export capabilities

### Phase 4: Distribution (Planned)
- [ ] CI/CD pipeline
- [ ] Cross-platform builds
- [ ] Homebrew formula
- [ ] Package managers (apt, snap)

## Credits

- **Weather Data**: [MET Norway](https://www.met.no/) (Meteorologisk institutt)
- **Inspiration**: Original bash script `weather-stavern.sh`
- **Author**: Kristoffer Risa

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

Special thanks to MET Norway for providing free, high-quality weather data through their excellent API.
