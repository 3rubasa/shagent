package simplerelay

import (
	"errors"
	"fmt"

	"github.com/3rubasa/shagent/pkg/businesslogic/interfaces"
)

type SimpleRelay struct {
	driver interfaces.RelayDriver
}

func New(driver interfaces.RelayDriver) *SimpleRelay {
	r := &SimpleRelay{
		driver: driver,
	}

	return r
}

func (r *SimpleRelay) Start() error {
	return r.driver.Start()
}

func (r *SimpleRelay) Stop() error {
	r.driver.Stop()

	return nil
}

func (r *SimpleRelay) TurnOn() error {
	return r.driver.TurnOn()
}

func (r *SimpleRelay) TurnOff() error {
	return r.driver.TurnOff()
}

func (r *SimpleRelay) Get() (int, error) {
	s, err := r.driver.GetState()
	if err != nil {
		fmt.Println("Failed to get state of the room light: ", err)
		return 0, err
	}

	switch s {
	case "on":
		return 1, nil
	case "off":
		return 0, nil
	default:
		fmt.Println("Unexpected state: ", s)
		return 0, errors.New("unexpected state")
	}
}
