package power

import (
	"fmt"

	"github.com/3rubasa/shagent/sensors"
	"github.com/stianeikeland/go-rpio"
)

const pin = 16

type powerStatusProvider struct {
	pin rpio.Pin
}

var providerSingleton *powerStatusProvider

func New() sensors.PowerStatusProvider {
	if providerSingleton == nil {
		providerSingleton = &powerStatusProvider{}
	}

	return providerSingleton
}

func (p *powerStatusProvider) Initialize() error {
	err := rpio.Open()
	if err != nil {
		fmt.Printf("Failed to open rpio: %s", err.Error())
		return err
	}

	p.pin = rpio.Pin(pin)
	p.pin.Input()

	return nil
}

func (p *powerStatusProvider) GetPowerStatus() (int, error) {
	s := p.pin.Read()

	if s == rpio.Low {
		return 0, nil
	} else {
		return 1, nil
	}
}
