package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadFromFile(path string) (*Config, error) {
	cfgFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open config file: ", err)
		return nil, err
	}
	defer cfgFile.Close()

	var cfg *Config
	err = json.NewDecoder(cfgFile).Decode(&cfg)
	if err != nil {
		fmt.Println("Failed to read config file: ", err)
		return nil, err
	}

	return cfg, nil
}

type Config struct {
	BusinessLogic    BusinessLogicConfig    `json:"business_logic"`
	WeatherProvider  WeatherProviderConfig  `json:"weather_provider"`
	ForecastProvider ForecastProviderConfig `json:"forecast_provider"`
}

type InetCheckerConfig struct {
	Enabled     bool   `json:"enabled"`
	URL         string `json:"url"`
	LongPeriod  int    `json:"long_period"`
	ShortPeriod int    `json:"short_period"`
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
