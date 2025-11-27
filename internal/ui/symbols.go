package ui

// WeatherSymbol returns emoji and description for a weather symbol code
func WeatherSymbol(symbolCode string) (emoji string, description string) {
	symbols := map[string]struct {
		emoji string
		desc  string
	}{
		"clearsky_day":            {"☀️", "Clear sky"},
		"clearsky_night":          {"☀️", "Clear sky"},
		"fair_day":                {"🌤️", "Fair"},
		"fair_night":              {"🌤️", "Fair"},
		"partlycloudy_day":        {"⛅", "Partly cloudy"},
		"partlycloudy_night":      {"⛅", "Partly cloudy"},
		"cloudy":                  {"☁️", "Cloudy"},
		"lightrain":               {"🌦️", "Light rain"},
		"rain":                    {"🌧️", "Rain"},
		"heavyrain":               {"⛈️", "Heavy rain"},
		"lightrainshowers_day":    {"🌦️", "Light rain showers"},
		"lightrainshowers_night":  {"🌦️", "Light rain showers"},
		"rainshowers_day":         {"🌧️", "Rain showers"},
		"rainshowers_night":       {"🌧️", "Rain showers"},
		"heavyrainshowers_day":    {"⛈️", "Heavy rain showers"},
		"heavyrainshowers_night":  {"⛈️", "Heavy rain showers"},
		"lightsleet":              {"🌨️", "Light sleet"},
		"sleet":                   {"🌨️", "Sleet"},
		"heavysleet":              {"🌨️", "Heavy sleet"},
		"lightsleetshowers_day":   {"🌨️", "Light sleet showers"},
		"lightsleetshowers_night": {"🌨️", "Light sleet showers"},
		"sleetshowers_day":        {"🌨️", "Sleet showers"},
		"sleetshowers_night":      {"🌨️", "Sleet showers"},
		"heavysleetshowers_day":   {"🌨️", "Heavy sleet showers"},
		"heavysleetshowers_night": {"🌨️", "Heavy sleet showers"},
		"lightsnow":               {"❄️", "Light snow"},
		"snow":                    {"🌨️", "Snow"},
		"heavysnow":               {"❄️", "Heavy snow"},
		"lightsnowshowers_day":    {"❄️", "Light snow showers"},
		"lightsnowshowers_night":  {"❄️", "Light snow showers"},
		"snowshowers_day":         {"🌨️", "Snow showers"},
		"snowshowers_night":       {"🌨️", "Snow showers"},
		"heavysnowshowers_day":    {"❄️", "Heavy snow showers"},
		"heavysnowshowers_night":  {"❄️", "Heavy snow showers"},
		"fog":                     {"🌫️", "Fog"},
		"lightrainandthunder":     {"⛈️", "Light rain and thunder"},
		"rainandthunder":          {"⛈️", "Rain and thunder"},
		"heavyrainandthunder":     {"⛈️", "Heavy rain and thunder"},
		"lightsleetandthunder":    {"⛈️", "Light sleet and thunder"},
		"sleetandthunder":         {"⛈️", "Sleet and thunder"},
		"heavysleetandthunder":    {"⛈️", "Heavy sleet and thunder"},
		"lightsnowandthunder":     {"⛈️", "Light snow and thunder"},
		"snowandthunder":          {"⛈️", "Snow and thunder"},
		"heavysnowandthunder":     {"⛈️", "Heavy snow and thunder"},
	}

	if s, ok := symbols[symbolCode]; ok {
		return s.emoji, s.desc
	}

	// Default for unknown symbols
	return "🌡️", symbolCode
}

// WeatherDescription returns only the description for a symbol code
func WeatherDescription(symbolCode string) string {
	_, desc := WeatherSymbol(symbolCode)
	return desc
}

// WeatherEmoji returns only the emoji for a symbol code
func WeatherEmoji(symbolCode string) string {
	emoji, _ := WeatherSymbol(symbolCode)
	return emoji
}
