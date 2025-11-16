package met

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kristofferrisa/sky-cli/internal/models"
)

const (
	baseURL   = "https://api.met.no/weatherapi/locationforecast/2.0/compact"
	userAgent = "sky-cli/1.0 github.com/kristofferrisa/sky-cli"
)

// Client represents a MET Norway API client
type Client struct {
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new MET Norway API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: userAgent,
	}
}

// GetForecast fetches weather forecast for the given coordinates
func (c *Client) GetForecast(ctx context.Context, lat, lon float64) (*Response, error) {
	url := fmt.Sprintf("%s?lat=%.4f&lon=%.4f", baseURL, lat, lon)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetCurrentWeather fetches current weather conditions
func (c *Client) GetCurrentWeather(ctx context.Context, loc *models.Location) (*models.Weather, error) {
	resp, err := c.GetForecast(ctx, loc.Latitude, loc.Longitude)
	if err != nil {
		return nil, err
	}

	if len(resp.Properties.Timeseries) == 0 {
		return nil, fmt.Errorf("no weather data available")
	}

	// First timeseries entry is the current/nearest weather
	current := resp.Properties.Timeseries[0]

	// Get weather symbol from next_1_hours or next_6_hours
	symbol := ""
	precipitation := 0.0
	if current.Data.Next1Hours != nil {
		symbol = current.Data.Next1Hours.Summary.SymbolCode
		precipitation = current.Data.Next1Hours.Details.PrecipitationAmount
	} else if current.Data.Next6Hours != nil {
		symbol = current.Data.Next6Hours.Summary.SymbolCode
		precipitation = current.Data.Next6Hours.Details.PrecipitationAmount
	}

	weather := &models.Weather{
		Location:      loc,
		Timestamp:     current.Time,
		UpdatedAt:     resp.Properties.Meta.UpdatedAt,
		Temperature:   current.Data.Instant.Details.AirTemperature,
		Humidity:      current.Data.Instant.Details.RelativeHumidity,
		Pressure:      current.Data.Instant.Details.AirPressureAtSeaLevel,
		CloudCover:    current.Data.Instant.Details.CloudAreaFraction,
		WindSpeed:     current.Data.Instant.Details.WindSpeed,
		WindDir:       current.Data.Instant.Details.WindFromDirection,
		Precipitation: precipitation,
		Symbol:        symbol,
	}

	return weather, nil
}

// GetHourlyForecast fetches hourly forecast for the specified number of hours
func (c *Client) GetHourlyForecast(ctx context.Context, loc *models.Location, hours int) (*models.Forecast, error) {
	resp, err := c.GetForecast(ctx, loc.Latitude, loc.Longitude)
	if err != nil {
		return nil, err
	}

	if len(resp.Properties.Timeseries) == 0 {
		return nil, fmt.Errorf("no forecast data available")
	}

	// Limit to requested hours or available data
	maxHours := hours
	if maxHours > len(resp.Properties.Timeseries) {
		maxHours = len(resp.Properties.Timeseries)
	}

	forecast := &models.Forecast{
		Location: loc,
		Hours:    make([]models.HourlyForecast, 0, maxHours),
	}

	for i := 0; i < maxHours; i++ {
		ts := resp.Properties.Timeseries[i]

		symbol := ""
		precipitation := 0.0
		if ts.Data.Next1Hours != nil {
			symbol = ts.Data.Next1Hours.Summary.SymbolCode
			precipitation = ts.Data.Next1Hours.Details.PrecipitationAmount
		} else if ts.Data.Next6Hours != nil {
			symbol = ts.Data.Next6Hours.Summary.SymbolCode
			precipitation = ts.Data.Next6Hours.Details.PrecipitationAmount
		}

		hourly := models.HourlyForecast{
			Time:          ts.Time,
			Temperature:   ts.Data.Instant.Details.AirTemperature,
			Humidity:      ts.Data.Instant.Details.RelativeHumidity,
			WindSpeed:     ts.Data.Instant.Details.WindSpeed,
			Precipitation: precipitation,
			Symbol:        symbol,
		}

		forecast.Hours = append(forecast.Hours, hourly)
	}

	return forecast, nil
}

// GetDailySummary calculates daily summary from hourly data
func (c *Client) GetDailySummary(ctx context.Context, loc *models.Location) (*models.DailySummary, error) {
	forecast, err := c.GetHourlyForecast(ctx, loc, 24)
	if err != nil {
		return nil, err
	}

	if len(forecast.Hours) == 0 {
		return nil, fmt.Errorf("no forecast data available")
	}

	summary := &models.DailySummary{
		Location: loc,
		Date:     forecast.Hours[0].Time,
	}

	// Calculate min, max, avg temperature
	minTemp := forecast.Hours[0].Temperature
	maxTemp := forecast.Hours[0].Temperature
	totalTemp := 0.0
	totalPrecip := 0.0

	for _, hour := range forecast.Hours {
		if hour.Temperature < minTemp {
			minTemp = hour.Temperature
		}
		if hour.Temperature > maxTemp {
			maxTemp = hour.Temperature
		}
		totalTemp += hour.Temperature
		totalPrecip += hour.Precipitation
	}

	summary.TemperatureMin = minTemp
	summary.TemperatureMax = maxTemp
	summary.TemperatureAvg = totalTemp / float64(len(forecast.Hours))
	summary.PrecipitationTotal = totalPrecip

	return summary, nil
}
