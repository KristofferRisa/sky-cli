package met

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kristofferrisa/sky-cli/internal/cache"
	"github.com/kristofferrisa/sky-cli/internal/models"
)

// CachedClient wraps the MET client with caching
type CachedClient struct {
	client *Client
	cache  cache.Cache
	ttl    time.Duration
}

// NewCachedClient creates a new cached MET client
func NewCachedClient(cache cache.Cache, ttl time.Duration) *CachedClient {
	return &CachedClient{
		client: NewClient(),
		cache:  cache,
		ttl:    ttl,
	}
}

// GetCurrentWeather fetches current weather with caching
func (c *CachedClient) GetCurrentWeather(ctx context.Context, loc *models.Location) (*models.Weather, error) {
	key := fmt.Sprintf("weather:current:%f:%f", loc.Latitude, loc.Longitude)

	// Try to get from cache
	if data, err := c.cache.Get(key); err == nil {
		var weather models.Weather
		if err := json.Unmarshal(data, &weather); err == nil {
			// Restore location reference
			weather.Location = loc
			return &weather, nil
		}
		// If unmarshal fails, fall through to fetch fresh data
	}

	// Fetch from API
	weather, err := c.client.GetCurrentWeather(ctx, loc)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(weather); err == nil {
		c.cache.Set(key, data, c.ttl)
	}

	return weather, nil
}

// GetHourlyForecast fetches hourly forecast with caching
func (c *CachedClient) GetHourlyForecast(ctx context.Context, loc *models.Location, hours int) (*models.Forecast, error) {
	key := fmt.Sprintf("weather:forecast:%f:%f:%d", loc.Latitude, loc.Longitude, hours)

	// Try to get from cache
	if data, err := c.cache.Get(key); err == nil {
		var forecast models.Forecast
		if err := json.Unmarshal(data, &forecast); err == nil {
			// Restore location reference
			forecast.Location = loc
			return &forecast, nil
		}
	}

	// Fetch from API
	forecast, err := c.client.GetHourlyForecast(ctx, loc, hours)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(forecast); err == nil {
		c.cache.Set(key, data, c.ttl)
	}

	return forecast, nil
}

// GetDailySummary fetches daily summary with caching
func (c *CachedClient) GetDailySummary(ctx context.Context, loc *models.Location) (*models.DailySummary, error) {
	key := fmt.Sprintf("weather:summary:%f:%f", loc.Latitude, loc.Longitude)

	// Try to get from cache
	if data, err := c.cache.Get(key); err == nil {
		var summary models.DailySummary
		if err := json.Unmarshal(data, &summary); err == nil {
			// Restore location reference
			summary.Location = loc
			return &summary, nil
		}
	}

	// Fetch from API
	summary, err := c.client.GetDailySummary(ctx, loc)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(summary); err == nil {
		c.cache.Set(key, data, c.ttl)
	}

	return summary, nil
}
