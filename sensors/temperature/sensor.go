package temperature

import (
	"fmt"

	"github.com/3rubasa/go-dht"
	"github.com/3rubasa/shagent/sensors"
)

const DefaultGPIOPin = "GPIO4"
const maxRetries = 11

type temperatureProvider struct {
	sensor *dht.DHT
}

var providerSingleton *temperatureProvider

func New() sensors.TemperatureProvider {
	if providerSingleton == nil {
		providerSingleton = &temperatureProvider{}
	}

	return providerSingleton
}

func (p *temperatureProvider) Initialize() error {
	err := dht.HostInit()
	if err != nil {
		fmt.Printf("Error in HostInit(): %s", err.Error())
	}

	p.sensor, err = dht.NewDHT(DefaultGPIOPin, dht.Celsius, "dht22")
	if err != nil {
		fmt.Printf("Error in NewDHT(): %s", err.Error())
	}

	return nil
}

func (p *temperatureProvider) GetTemperature() (float64, error) {
	_, t, err := p.sensor.ReadRetry(maxRetries)
	if err != nil {
		fmt.Printf("Error while reading a sample: %s", err.Error())
	}

	return t, nil
}
