package api

import (
	"context"

	"github.com/kristofferrisa/sky-cli/internal/models"
)

// WeatherClient is the interface for weather API clients
type WeatherClient interface {
	GetCurrentWeather(ctx context.Context, loc *models.Location) (*models.Weather, error)
	GetHourlyForecast(ctx context.Context, loc *models.Location, hours int) (*models.Forecast, error)
	GetDailySummary(ctx context.Context, loc *models.Location) (*models.DailySummary, error)
}
