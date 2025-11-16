#!/bin/bash
#
# Weather Script for Stavern, Norway
# Fetches current weather from MET Norway (Meteorologisk institutt)
# Output format: LLM/AI-friendly structured text
#
# Location: Stavern, Norway
# Coordinates: 59.0°N, 10.0°E
#
# Usage: ./weather-stavern.sh [--json|--summary|--full]
#

set -euo pipefail

# Configuration
LATITUDE="59.0"
LONGITUDE="10.0"
LOCATION_NAME="Stavern, Norway"
API_URL="https://api.met.no/weatherapi/locationforecast/2.0/compact"
USER_AGENT="WeatherScript-Stavern/1.0"

# Output mode (default: full)
OUTPUT_MODE="${1:-full}"

# Color codes for terminal (optional, works well in most terminals)
BOLD='\033[1m'
RESET='\033[0m'
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'

# Function to print section headers
print_header() {
  echo -e "\n${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
  echo -e "${BOLD}$1${RESET}"
  echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
}

# Function to format timestamp to human-readable
format_time() {
  local timestamp=$1
  date -j -f "%Y-%m-%dT%H:%M:%SZ" "$timestamp" "+%A, %B %d, %Y at %H:%M" 2>/dev/null || echo "$timestamp"
}

# Function to format short time
format_short_time() {
  local timestamp=$1
  date -j -f "%Y-%m-%dT%H:%M:%SZ" "$timestamp" "+%H:%M" 2>/dev/null || echo "$timestamp"
}

# Function to get weather symbol description
get_weather_description() {
  local symbol=$1
  case $symbol in
  clearsky_day | clearsky_night) echo "☀️  Clear sky" ;;
  fair_day | fair_night) echo "🌤️  Fair" ;;
  partlycloudy_day | partlycloudy_night) echo "⛅ Partly cloudy" ;;
  cloudy) echo "☁️  Cloudy" ;;
  lightrain) echo "🌦️  Light rain" ;;
  rain) echo "🌧️  Rain" ;;
  heavyrain) echo "⛈️  Heavy rain" ;;
  lightrainshowers_day | lightrainshowers_night) echo "🌦️  Light rain showers" ;;
  rainshowers_day | rainshowers_night) echo "🌧️  Rain showers" ;;
  heavyrainshowers_day | heavyrainshowers_night) echo "⛈️  Heavy rain showers" ;;
  lightsleet) echo "🌨️  Light sleet" ;;
  sleet) echo "🌨️  Sleet" ;;
  heavysleet) echo "🌨️  Heavy sleet" ;;
  lightsnow) echo "❄️  Light snow" ;;
  snow) echo "🌨️  Snow" ;;
  heavysnow) echo "❄️  Heavy snow" ;;
  fog) echo "🌫️  Fog" ;;
  *) echo "🌡️  $symbol" ;;
  esac
}

# Function to get wind direction
get_wind_direction() {
  local degrees=$1
  local deg=${degrees%.*} # Remove decimal part

  if [ "$deg" -ge 338 ] || [ "$deg" -lt 23 ]; then
    echo "N (North)"
  elif [ "$deg" -ge 23 ] && [ "$deg" -lt 68 ]; then
    echo "NE (Northeast)"
  elif [ "$deg" -ge 68 ] && [ "$deg" -lt 113 ]; then
    echo "E (East)"
  elif [ "$deg" -ge 113 ] && [ "$deg" -lt 158 ]; then
    echo "SE (Southeast)"
  elif [ "$deg" -ge 158 ] && [ "$deg" -lt 203 ]; then
    echo "S (South)"
  elif [ "$deg" -ge 203 ] && [ "$deg" -lt 248 ]; then
    echo "SW (Southwest)"
  elif [ "$deg" -ge 248 ] && [ "$deg" -lt 293 ]; then
    echo "W (West)"
  elif [ "$deg" -ge 293 ] && [ "$deg" -lt 338 ]; then
    echo "NW (Northwest)"
  fi
}

# Function to get wind speed description
get_wind_description() {
  local speed=$1
  local speed_int=${speed%.*}

  if [ "$speed_int" -lt 2 ]; then
    echo "Calm"
  elif [ "$speed_int" -lt 6 ]; then
    echo "Light breeze"
  elif [ "$speed_int" -lt 12 ]; then
    echo "Moderate breeze"
  elif [ "$speed_int" -lt 20 ]; then
    echo "Strong breeze"
  elif [ "$speed_int" -lt 29 ]; then
    echo "Gale"
  else
    echo "Storm"
  fi
}

# Fetch weather data
WEATHER_DATA=$(curl -s "${API_URL}?lat=${LATITUDE}&lon=${LONGITUDE}" -A "${USER_AGENT}")

if [ -z "$WEATHER_DATA" ]; then
  echo "❌ Error: Failed to fetch weather data"
  exit 1
fi

# Check if output is JSON mode (skip all formatting)
if [ "$OUTPUT_MODE" = "json" ]; then
  echo "$WEATHER_DATA" | jq '.'
  exit 0
fi

# Print header for non-JSON modes
print_header "Fetching Weather Data for $LOCATION_NAME"
echo "API: MET Norway (Meteorologisk institutt)"
echo "Coordinates: ${LATITUDE}°N, ${LONGITUDE}°E"
echo "Request time: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Parse weather data using jq
UPDATED_AT=$(echo "$WEATHER_DATA" | jq -r '.properties.meta.updated_at')
UNITS=$(echo "$WEATHER_DATA" | jq -r '.properties.meta.units')

# Get current weather (first timeseries entry)
CURRENT=$(echo "$WEATHER_DATA" | jq -r '.properties.timeseries[0]')
CURRENT_TIME=$(echo "$CURRENT" | jq -r '.time')
CURRENT_TEMP=$(echo "$CURRENT" | jq -r '.data.instant.details.air_temperature')
CURRENT_HUMIDITY=$(echo "$CURRENT" | jq -r '.data.instant.details.relative_humidity')
CURRENT_WIND_SPEED=$(echo "$CURRENT" | jq -r '.data.instant.details.wind_speed')
CURRENT_WIND_DIR=$(echo "$CURRENT" | jq -r '.data.instant.details.wind_from_direction')
CURRENT_PRESSURE=$(echo "$CURRENT" | jq -r '.data.instant.details.air_pressure_at_sea_level')
CURRENT_CLOUD=$(echo "$CURRENT" | jq -r '.data.instant.details.cloud_area_fraction')
CURRENT_SYMBOL=$(echo "$CURRENT" | jq -r '.data.next_1_hours.summary.symbol_code // .data.next_6_hours.summary.symbol_code')
CURRENT_PRECIP=$(echo "$CURRENT" | jq -r '.data.next_1_hours.details.precipitation_amount // 0')

# Print current weather
print_header "CURRENT WEATHER"
echo -e "${BOLD}Location:${RESET}     $LOCATION_NAME"
echo -e "${BOLD}Time:${RESET}         $(format_time "$CURRENT_TIME")"
echo -e "${BOLD}Updated:${RESET}      $(format_time "$UPDATED_AT")"
echo ""
echo -e "${BOLD}${GREEN}Conditions:${RESET}   $(get_weather_description "$CURRENT_SYMBOL")"
echo -e "${BOLD}Temperature:${RESET}  ${CURRENT_TEMP}°C"
echo -e "${BOLD}Humidity:${RESET}     ${CURRENT_HUMIDITY}%"
echo -e "${BOLD}Cloud Cover:${RESET}  ${CURRENT_CLOUD}%"
echo -e "${BOLD}Precipitation:${RESET} ${CURRENT_PRECIP} mm (next hour)"
echo ""
echo -e "${BOLD}Wind:${RESET}         ${CURRENT_WIND_SPEED} m/s from $(get_wind_direction "$CURRENT_WIND_DIR")"
echo -e "${BOLD}Wind Status:${RESET}  $(get_wind_description "$CURRENT_WIND_SPEED")"
echo -e "${BOLD}Pressure:${RESET}     ${CURRENT_PRESSURE} hPa"

# Short summary mode exits here
if [ "$OUTPUT_MODE" = "summary" ]; then
  exit 0
fi

# Get next 12 hours forecast
print_header "HOURLY FORECAST (Next 12 Hours)"
echo -e "${BOLD}Time     Temp    Symbol                Precip  Wind    Humidity${RESET}"
echo "────────────────────────────────────────────────────────────────────"

for i in {0..11}; do
  HOUR_DATA=$(echo "$WEATHER_DATA" | jq -r ".properties.timeseries[$i]")
  if [ "$HOUR_DATA" = "null" ]; then
    break
  fi

  H_TIME=$(echo "$HOUR_DATA" | jq -r '.time')
  H_TEMP=$(echo "$HOUR_DATA" | jq -r '.data.instant.details.air_temperature')
  H_SYMBOL=$(echo "$HOUR_DATA" | jq -r '.data.next_1_hours.summary.symbol_code // .data.next_6_hours.summary.symbol_code // "unknown"')
  H_PRECIP=$(echo "$HOUR_DATA" | jq -r '.data.next_1_hours.details.precipitation_amount // 0')
  H_WIND=$(echo "$HOUR_DATA" | jq -r '.data.instant.details.wind_speed')
  H_HUMIDITY=$(echo "$HOUR_DATA" | jq -r '.data.instant.details.relative_humidity')

  H_DESC=$(get_weather_description "$H_SYMBOL" | sed 's/[🌤️☀️⛅☁️🌦️🌧️⛈️🌨️❄️🌫️🌡️]//g' | xargs)

  printf "%-8s %-7s %-20s %-7s %-7s %s%%\n" \
    "$(format_short_time "$H_TIME")" \
    "${H_TEMP}°C" \
    "$H_DESC" \
    "${H_PRECIP}mm" \
    "${H_WIND}m/s" \
    "$H_HUMIDITY"
done

# Get daily summary
print_header "DAILY SUMMARY"

# Calculate daily statistics
DAILY_TEMPS=$(echo "$WEATHER_DATA" | jq -r '[.properties.timeseries[0:24] | .[].data.instant.details.air_temperature] | @json')
DAILY_PRECIP=$(echo "$WEATHER_DATA" | jq -r '[.properties.timeseries[0:24] | .[].data.next_1_hours.details.precipitation_amount // 0] | add')

MIN_TEMP=$(echo "$DAILY_TEMPS" | jq 'min')
MAX_TEMP=$(echo "$DAILY_TEMPS" | jq 'max')
AVG_TEMP=$(echo "$DAILY_TEMPS" | jq 'add / length | floor')

echo -e "${BOLD}Temperature Range:${RESET}"
echo "  • Minimum: ${MIN_TEMP}°C"
echo "  • Maximum: ${MAX_TEMP}°C"
echo "  • Average: ${AVG_TEMP}°C"
echo ""
echo -e "${BOLD}Precipitation:${RESET}"
echo "  • Total (24h): ${DAILY_PRECIP} mm"
echo ""

# Check for weather warnings
if (($(echo "$CURRENT_WIND_SPEED > 15" | bc -l))); then
  echo -e "${BOLD}${YELLOW}⚠️  Weather Warnings:${RESET}"
  echo "  • Strong winds detected (${CURRENT_WIND_SPEED} m/s)"
fi

if (($(echo "$DAILY_PRECIP > 10" | bc -l))); then
  echo -e "${BOLD}${YELLOW}⚠️  Weather Warnings:${RESET}"
  echo "  • Significant precipitation expected (${DAILY_PRECIP} mm)"
fi

# LLM-friendly structured summary
print_header "STRUCTURED DATA (For LLM Processing)"
cat <<EOF
LOCATION: Stavern, Norway (${LATITUDE}°N, ${LONGITUDE}°E)
TIMESTAMP: $(date '+%Y-%m-%d %H:%M:%S')
DATA_SOURCE: MET Norway (yr.no/met.no)

CURRENT_CONDITIONS:
  temperature: ${CURRENT_TEMP}°C
  weather: $(get_weather_description "$CURRENT_SYMBOL" | sed 's/[🌤️☀️⛅☁️🌦️🌧️⛈️🌨️❄️🌫️🌡️]//g' | xargs)
  humidity: ${CURRENT_HUMIDITY}%
  wind_speed: ${CURRENT_WIND_SPEED} m/s
  wind_direction: $(get_wind_direction "$CURRENT_WIND_DIR")
  wind_description: $(get_wind_description "$CURRENT_WIND_SPEED")
  pressure: ${CURRENT_PRESSURE} hPa
  cloud_cover: ${CURRENT_CLOUD}%
  precipitation_next_hour: ${CURRENT_PRECIP} mm

DAILY_SUMMARY:
  temperature_min: ${MIN_TEMP}°C
  temperature_max: ${MAX_TEMP}°C
  temperature_avg: ${AVG_TEMP}°C
  precipitation_24h: ${DAILY_PRECIP} mm

DATA_UNITS:
  temperature: celsius
  wind_speed: meters_per_second
  pressure: hectopascal (hPa)
  precipitation: millimeters (mm)
  humidity: percent
EOF

print_header "Weather Data Retrieved Successfully"
echo "✅ Data is fresh and ready for analysis"
echo ""

# Add recommendation section
echo -e "${BOLD}AI Analysis Suggestions:${RESET}"
echo "• Current conditions are suitable for: [to be analyzed]"
echo "• Weather trend for next 12 hours: [to be analyzed]"
echo "• Recommended activities: [to be analyzed based on conditions]"
echo "• Travel advisories: [to be analyzed based on wind/precipitation]"
echo ""

exit 0
