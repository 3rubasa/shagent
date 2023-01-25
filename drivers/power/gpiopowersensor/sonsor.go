package gpiopowersensor

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

type Sensor struct {
	gpioPin uint8
	pin     rpio.Pin
}

func New(gpioPin uint8) *Sensor {
	return &Sensor{
		gpioPin: gpioPin,
	}
}

func (s *Sensor) Initialize() error {
	err := rpio.Open()
	if err != nil {
		fmt.Printf("Failed to open rpio: %s \n", err.Error())
		return err
	}

	s.pin = rpio.Pin(s.gpioPin)
	s.pin.Input()

	return nil
}

func (s *Sensor) Get() (int, error) {
	state := s.pin.Read()

	if state == rpio.Low {
		return 0, nil
	} else {
		return 1, nil
	}
}
