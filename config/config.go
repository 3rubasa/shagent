package config

type Config struct {
	BusinessLogic    BusinessLogicConfig    `json:"business_logic"`
	WeatherProvider  WeatherProviderConfig  `json:"weather_provider"`
	ForecastProvider ForecastProviderConfig `json:"forecast_provider"`
}

type BusinessLogicConfig struct {
	Consumer ConsumerConfig `json:"consumer"`
}

type ConsumerConfig struct {
	Enabled bool   `json:"enabled"`
	APIKeys string `json:"api_keys"`
	Address string `json:"address"`
	Schema  string `json:"schema"`
	URI     string `json:"uri"`
}

type WeatherProviderConfig struct {
	Enabled bool `json:"enabled"`
}

type ForecastProviderConfig struct {
	Enabled bool `json:"enabled"`
}
