package wsraspihatx3

import (
	"log"

	"github.com/stianeikeland/go-rpio"
)

type Driver struct {
	channel RelayChannel
	pin     rpio.Pin
}

func New(channel RelayChannel) (*Driver, error) {
	err := rpio.Open()
	if err != nil {
		log.Println("ERROR: Failed to open rpio device: ", err)
		return nil, err
	}

	pin := rpio.Pin(chanToGPIO[channel])
	pin.Output()
	pin.High()

	return &Driver{
		channel: channel,
		pin:     pin,
	}, nil
}

func (a Driver) GetState() (string, error) {
	s := a.pin.Read()

	switch s {
	case rpio.High:
		return "off", nil
	case rpio.Low:
		return "on", nil
	default:
		log.Panicln("ERROR: unexpected gpio state: ", s)
		return "", nil
	}
}

func (a Driver) TurnOn() error {
	a.pin.Low()
	return nil
}

func (a Driver) TurnOff() error {
	a.pin.High()
	return nil
}
