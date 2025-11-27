package met

import "time"

// Response represents the root structure of MET Norway API response
type Response struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

// Geometry contains location information
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"` // [longitude, latitude, altitude]
}

// Properties contains the actual weather data
type Properties struct {
	Meta       Meta         `json:"meta"`
	Timeseries []Timeseries `json:"timeseries"`
}

// Meta contains metadata about the forecast
type Meta struct {
	UpdatedAt time.Time `json:"updated_at"`
	Units     Units     `json:"units"`
}

// Units describes the units used in the data
type Units struct {
	AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level"`
	AirTemperature        string `json:"air_temperature"`
	CloudAreaFraction     string `json:"cloud_area_fraction"`
	PrecipitationAmount   string `json:"precipitation_amount"`
	RelativeHumidity      string `json:"relative_humidity"`
	WindFromDirection     string `json:"wind_from_direction"`
	WindSpeed             string `json:"wind_speed"`
}

// Timeseries represents a single time point in the forecast
type Timeseries struct {
	Time time.Time `json:"time"`
	Data Data      `json:"data"`
}

// Data contains instant and forecast data
type Data struct {
	Instant     Instant     `json:"instant"`
	Next1Hours  *NextNHours `json:"next_1_hours,omitempty"`
	Next6Hours  *NextNHours `json:"next_6_hours,omitempty"`
	Next12Hours *NextNHours `json:"next_12_hours,omitempty"`
}

// Instant contains current conditions
type Instant struct {
	Details InstantDetails `json:"details"`
}

// InstantDetails contains the actual instant weather values
type InstantDetails struct {
	AirPressureAtSeaLevel    float64 `json:"air_pressure_at_sea_level"`
	AirTemperature           float64 `json:"air_temperature"`
	CloudAreaFraction        float64 `json:"cloud_area_fraction"`
	RelativeHumidity         float64 `json:"relative_humidity"`
	WindFromDirection        float64 `json:"wind_from_direction"`
	WindSpeed                float64 `json:"wind_speed"`
	FogAreaFraction          float64 `json:"fog_area_fraction,omitempty"`
	UltravioletIndexClearSky float64 `json:"ultraviolet_index_clear_sky,omitempty"`
}

// NextNHours contains forecast for the next N hours
type NextNHours struct {
	Summary Summary         `json:"summary"`
	Details ForecastDetails `json:"details"`
}

// Summary contains weather symbol information
type Summary struct {
	SymbolCode string `json:"symbol_code"`
}

// ForecastDetails contains forecast-specific details
type ForecastDetails struct {
	PrecipitationAmount        float64 `json:"precipitation_amount,omitempty"`
	PrecipitationAmountMax     float64 `json:"precipitation_amount_max,omitempty"`
	PrecipitationAmountMin     float64 `json:"precipitation_amount_min,omitempty"`
	ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation,omitempty"`
	ProbabilityOfThunder       float64 `json:"probability_of_thunder,omitempty"`
}
