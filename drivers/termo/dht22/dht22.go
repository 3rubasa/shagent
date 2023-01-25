package dht22

import (
	"fmt"

	"github.com/3rubasa/go-dht"
)

type Sensor struct {
	pin       uint8
	dhtSensor *dht.DHT
}

func New(Pin uint8) *Sensor {
	return &Sensor{
		pin: Pin,
	}
}

func (s *Sensor) Initialize() error {
	err := dht.HostInit()
	if err != nil {
		fmt.Printf("Error in HostInit(): %s \n", err.Error())
		return err
	}

	s.dhtSensor, err = dht.NewDHT(fmt.Sprintf("GPIO%d", s.pin), dht.Celsius, sensorType)
	if err != nil {
		fmt.Printf("Error in NewDHT(): %s \n", err.Error())
		return err
	}

	return nil
}

func (s *Sensor) Get() (float64, error) {
	_, t, err := s.dhtSensor.ReadRetry(maxRetries)
	if err != nil {
		fmt.Printf("Error while reading a sample: %s \n", err.Error())
		return -273.15, err
	}

	return t, nil
}
