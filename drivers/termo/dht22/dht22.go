package dht22

import (
	"fmt"
	"log"

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
		log.Println("ERROR: Failed to initialize DHT sersor: ", err)
		return err
	}

	s.dhtSensor, err = dht.NewDHT(fmt.Sprintf("GPIO%d", s.pin), dht.Celsius, sensorType)
	if err != nil {
		log.Println("ERROR: Failed to instantiate DHT struct: ", err)
		return err
	}

	return nil
}

func (s *Sensor) Get() (float64, error) {
	_, t, err := s.dhtSensor.ReadRetry(maxRetries)
	if err != nil {
		log.Println("Debug: Failed to read DHT sensor data: ", err)
		return 0.0, err
	}

	return t, nil
}
