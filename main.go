package main

import (
	"fmt"
	"time"

	"github.com/3rubasa/shagent/sensors/temperature"
)

const sampleInterval = time.Second * 5

func main() {
	tp := temperature.New()

	err := tp.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize temperature sensor: ", err)
	}

	for {
		var t float64

		t, err := tp.GetTemperature()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println("Current temperature is ", t)

		time.Sleep(sampleInterval)
	}
}
