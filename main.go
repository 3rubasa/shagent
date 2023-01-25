package main

import (
	"fmt"

	//"os"
	//"syscall"
	"time"

	//"bitbucket.org/gmcbay/i2c"

	"github.com/3rubasa/shagent/businesslogic"
	"github.com/3rubasa/shagent/businesslogic/controllers/power"
	"github.com/3rubasa/shagent/businesslogic/controllers/roomlight"
	"github.com/3rubasa/shagent/businesslogic/controllers/temperature"
	"github.com/3rubasa/shagent/drivers/power/gpiopowersensor"
	"github.com/3rubasa/shagent/drivers/relay"
	"github.com/3rubasa/shagent/drivers/relay/sonoffr3rf"
	"github.com/3rubasa/shagent/drivers/relay/wsraspihatx3"
	"github.com/3rubasa/shagent/drivers/termo/dht22"
	"github.com/3rubasa/shagent/osservices"
	"github.com/3rubasa/shagent/watchdog"
	"github.com/3rubasa/shagent/webserver"
)

const sampleInterval = time.Second * 60

func main() {
	var err error

	// Common
	osservices := osservices.NewOSServicesProvider()

	// 1 - watchdog DONE
	inetchecker := watchdog.NewInternetChecker("http://google.com")
	wd := watchdog.New(osservices, inetchecker, 30*time.Minute)
	wd.Start()

	// 2 - boiler DONE
	sonoffr3rfRelay := sonoffr3rf.New(osservices, "24:a1:60:1d:72:9d")
	boiler := relay.New(sonoffr3rfRelay, 30*time.Minute)

	// 3 - roomLight
	// TODO: Later, if roomLight is nil, what are we going to do?
	var roomLight *relay.Relay
	wsRelayForRoomLight, err := wsraspihatx3.New(wsraspihatx3.RelayChannel1)
	if err != nil {
		fmt.Println("Failed to create WaveShare Raspi Hat relay device: ", err)
	} else {
		roomLight = relay.New(wsRelayForRoomLight, 30*time.Minute)
	}

	// 4 - camLight
	// TODO: Later, if cam light is nil, what are we going to do?
	var camLight *relay.Relay
	wsRelayForCamLight, err := wsraspihatx3.New(wsraspihatx3.RelayChannel2)
	if err != nil {
		fmt.Println("Failed to create WaveShare Raspi Hat relay device for cam light: ", err)
	} else {
		camLight = relay.New(wsRelayForCamLight, 30*time.Minute)
	}

	// 5 - power
	powerDriver := gpiopowersensor.New(16)
	powerController := power.New(powerDriver)

	// 6 - webserver
	ws := webserver.New(boiler, roomLight, camLight)
	err = ws.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize the web server: ", err)
		return
	}
	err = ws.Start()
	if err != nil {
		fmt.Println("Failed to start the web server: ", err)
		return
	}

	// 5 - Business logic
	var ontimes, offtimes []string
	ontimes = append(ontimes, "0 45 06 * * *", "0 10 17 * * *")
	offtimes = append(offtimes, "0 15 08 * * *", "0 12 01 * * *")

	roomLightController := roomlight.New(roomLight, ontimes, offtimes)

	// Kitchen Temperature Sensor
	tempSensorDrv := dht22.New(4)
	kitchenTempController := temperature.New(tempSensorDrv, 10*time.Minute, time.Minute)

	bl := businesslogic.New(roomLightController, kitchenTempController, powerController, time.Minute)
	bl.Start()
}
