# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sky CLI is a command-line weather tool written in Go that provides current weather, hourly forecasts, and daily forecasts using the MET Norway API. The project emphasizes clean architecture with pluggable components, multiple output formats (full, JSON, summary, markdown), and performance optimization through file-based caching.

## Development Commands

### Build
```bash
go build -o sky ./cmd/sky
```

### Run Locally
```bash
go run ./cmd/sky current
go run ./cmd/sky forecast --hours 24
go run ./cmd/sky daily --days 7
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/models -v
go test ./internal/cache -v
```

### Code Quality
```bash
# Format code
gofmt -w .

# Run static analysis
go vet ./...

# Verify dependencies
go mod verify
```

### Release
The project uses GoReleaser for automated builds and releases. CI/CD is configured via GitHub Actions.

## Architecture

### Core Architectural Patterns

**1. Provider Interface Pattern (WeatherClient)**
- Location: `internal/api/client.go`
- The `WeatherClient` interface abstracts weather data providers
- Current implementation: MET Norway (`internal/api/met/`)
- Enables future support for additional providers (OpenWeather, Weather.gov, etc.)
- All commands interact through this interface, never directly with providers

**2. Decorator Pattern (Caching Layer)**
- Location: `internal/api/met/cached_client.go`
- `CachedClient` wraps the base MET client with caching behavior
- Cache keys incorporate location coordinates and request parameters
- Cache misses fall through to the underlying provider
- Created in `cmd/sky/root.go:getWeatherClient()` based on config settings

**3. Strategy Pattern (Formatters)**
- Location: `internal/formatter/`
- Four formatter implementations: Full, JSON, Summary, Markdown
- All implement the `Formatter` interface with four methods:
  - `FormatCurrent()` - current weather
  - `FormatForecast()` - hourly forecasts
  - `FormatDailySummary()` - single day summary
  - `FormatDailyForecast()` - multi-day forecasts
- Factory creates formatters by name (`internal/formatter/factory.go`)
- Formatters receive `Options` struct with NoColor, NoEmoji, and TimeFormat flags

**4. Cache Interface**
- Location: `internal/cache/cache.go`
- Simple interface: Get, Set, Delete, Clear, Has
- Implementations:
  - `FileCache`: stores cached data as JSON files in `~/.sky/cache/`
  - `NoOpCache`: disabled cache (returns cache miss for all operations)
- TTL is managed at the cache level with file modification times

### Data Flow

1. **Command Invocation** (`cmd/sky/current.go`, `forecast.go`, `daily.go`, `locations.go`)
   - Parse flags and resolve location (from config, coordinates, or default)
   - Create weather client via `getWeatherClient()` (with or without caching)
   - Create formatter via factory based on --format flag
   - Call appropriate WeatherClient method
   - Pass result to formatter's corresponding method
   - Write formatted output to stdout

2. **Weather Client Resolution** (`cmd/sky/root.go:getWeatherClient()`)
   - If caching disabled: return raw `met.NewClient()`
   - If caching enabled:
     - Create FileCache with directory from config
     - Get TTL from config (default 10 minutes)
     - Wrap MET client with `met.NewCachedClient(cache, ttl)`

3. **Cache Layer** (when enabled)
   - Generate cache key from method, location coordinates, and parameters
   - Check cache for valid (non-expired) data
   - On hit: deserialize and return cached data
   - On miss: fetch from API, serialize, store in cache with TTL, return data

### Key Components

**Configuration System** (`internal/config/config.go`)
- Uses Viper for YAML config management
- Config locations: `~/.sky/config.yaml` or `~/.config/sky/config.yaml`
- Stores saved locations map, default location, cache settings, output preferences
- Methods: `Load()`, `Save()`, `GetLocation()`, `AddLocation()`, `RemoveLocation()`

**Models** (`internal/models/`)
- `Location`: lat/lon with validation, timezone, and name
- `Weather`: current conditions with helper methods (`WindDirection()`, `WindDescription()`)
- `Forecast`: collection of `HourlyForecast` structs
- `DailySummary`: aggregated day data (min/max temps, total precip, most common symbol)
- `DailyForecast`: collection of `DailySummary` structs

**MET Norway Provider** (`internal/api/met/`)
- `client.go`: Core API client, fetches from MET Norway API
- `cached_client.go`: Wraps client with caching (decorator)
- `models.go`: MET-specific response structures
- Daily forecast aggregation: fetches hourly data and groups by date

**UI Helpers** (`internal/ui/`)
- `colors.go`: Terminal color utilities
- `symbols.go`: Weather emoji/symbol mapping

### Command Structure

All commands are Cobra commands registered in their respective files:
- `cmd/sky/root.go`: Root command, global flags, version command, client factory
- `cmd/sky/current.go`: Current weather with optional forecast and daily summary
- `cmd/sky/forecast.go`: Hourly forecast (default 12 hours)
- `cmd/sky/daily.go`: Daily forecast (default 7 days)
- `cmd/sky/locations.go`: Subcommands for location management (list, add, remove, set-default)

Each command follows the same pattern:
1. Define flags (format, location, lat/lon, hours/days)
2. In RunE function: resolve location, create client, fetch data, format, output
3. Handle errors with descriptive messages

### Configuration Details

Default configuration is embedded in `internal/config/config.go` via Viper defaults. The config includes:
- Default location: "stavern" (Norway)
- Cache enabled by default with 10-minute TTL
- Cache directory: `~/.sky/cache`
- Default format: "full"

Location resolution order (in commands):
1. Explicit --lat/--lon flags
2. --location flag (looks up in config)
3. Positional argument (looks up in config)
4. Default location from config
5. Error if none found

## Testing Strategy

- Unit tests focus on models and cache layers
- Test files: `*_test.go` alongside implementation
- Current coverage: Models 96.3%, Cache 56.7%
- Tests use table-driven approach for validation logic
- Mock external dependencies (MET API not called in tests)

## Common Patterns

**Adding a New Output Format:**
1. Create new formatter in `internal/formatter/` implementing `Formatter` interface
2. Implement all four format methods
3. Register in `factory.go:NewFormatter()`
4. Update README with new format documentation

**Adding a New Weather Provider:**
1. Implement `api.WeatherClient` interface in new package under `internal/api/`
2. Create provider-specific models for API responses
3. Add caching wrapper (similar to `met/cached_client.go`)
4. Update root command to support provider selection

**Adding a New Command:**
1. Create new file in `cmd/sky/` (e.g., `alerts.go`)
2. Define Cobra command with flags
3. Register in `root.go:init()`
4. Follow existing patterns for location resolution and formatting
5. Add corresponding method to `WeatherClient` interface if needed
6. Implement in all providers (currently just MET)
