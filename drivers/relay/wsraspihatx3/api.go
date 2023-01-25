package wsraspihatx3

import (
	"errors"
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

type API struct {
	channel RelayChannel
	pin     rpio.Pin
}

func New(channel RelayChannel) (*API, error) {
	err := rpio.Open()
	if err != nil {
		fmt.Printf("Failed to open rpio: %s \n", err.Error())
		return nil, err
	}

	pin := rpio.Pin(chanToGPIO[channel])
	pin.Output()
	pin.High()

	return &API{
		channel: channel,
		pin:     pin,
	}, nil
}

func (a API) GetState() (string, error) {
	s := a.pin.Read()

	switch s {
	case rpio.High:
		return "off", nil
	case rpio.Low:
		return "on", nil
	default:
		return "", errors.New("unexpected gpio state")
	}
}

func (a API) TurnOn() error {
	a.pin.Low()
	return nil
}

func (a API) TurnOff() error {
	a.pin.High()
	return nil
}
