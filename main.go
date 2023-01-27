package main

import (
	"fmt"

	//"os"
	//"syscall"
	"time"

	//"bitbucket.org/gmcbay/i2c"

	"github.com/3rubasa/shagent/businesslogic"
	"github.com/3rubasa/shagent/businesslogic/controllers/ltemodulecontroller"
	"github.com/3rubasa/shagent/businesslogic/controllers/power"
	scheduledrelay "github.com/3rubasa/shagent/businesslogic/controllers/relay/scheduled"
	simplerelay "github.com/3rubasa/shagent/businesslogic/controllers/relay/simple"
	"github.com/3rubasa/shagent/businesslogic/controllers/temperature"
	"github.com/3rubasa/shagent/drivers/ltemodule"
	"github.com/3rubasa/shagent/drivers/ltemodule/sim7600"
	"github.com/3rubasa/shagent/drivers/power/gpiopowersensor"
	"github.com/3rubasa/shagent/drivers/relay"
	"github.com/3rubasa/shagent/drivers/relay/sonoffr3rf"
	"github.com/3rubasa/shagent/drivers/relay/wsraspihatx3"
	"github.com/3rubasa/shagent/drivers/termo/dht22"
	"github.com/3rubasa/shagent/osservices"
	"github.com/3rubasa/shagent/watchdog"
	"github.com/3rubasa/shagent/webserver"
)

func main() {
	var err error

	// Common
	osservices := osservices.NewOSServicesProvider()

	// 1 - watchdog DONE
	inetchecker := watchdog.NewInternetChecker("http://google.com")
	wd := watchdog.New(osservices, inetchecker, 20*time.Minute, 5*time.Minute, 5*time.Minute)
	wd.Start()

	// 2 - boiler DONE
	sonoffr3rfRelayDrv := sonoffr3rf.New(osservices, "24:a1:60:1d:72:9d")
	boilerRelayDrv := relay.New(sonoffr3rfRelayDrv, 30*time.Minute)
	boilerController := simplerelay.New(boilerRelayDrv)

	// 3 - roomLight DONE
	// TODO: Later, if roomLight is nil, what are we going to do?
	var roomLightRelayDrv *relay.Relay
	wsRelayForRoomLight, err := wsraspihatx3.New(wsraspihatx3.RelayChannel1)
	if err != nil {
		fmt.Println("Failed to create WaveShare Raspi Hat relay device: ", err)
	} else {
		roomLightRelayDrv = relay.New(wsRelayForRoomLight, 30*time.Minute)
	}

	var ontimes, offtimes []string
	ontimes = append(ontimes, "0 45 06 * * *", "0 10 17 * * *")
	offtimes = append(offtimes, "0 15 08 * * *", "0 12 01 * * *")

	roomLightController := scheduledrelay.New(roomLightRelayDrv, ontimes, offtimes)

	// 4 - camLight DONE
	// TODO: Later, if cam light is nil, what are we going to do?
	var camLightRelayDrv *relay.Relay
	wsRelayForCamLight, err := wsraspihatx3.New(wsraspihatx3.RelayChannel2)
	if err != nil {
		fmt.Println("Failed to create WaveShare Raspi Hat relay device: ", err)
	} else {
		camLightRelayDrv = relay.New(wsRelayForCamLight, 30*time.Minute)
	}

	camLightController := simplerelay.New(camLightRelayDrv)

	// 5 - power DONE
	powerSensorDrv := gpiopowersensor.New(16)
	powerController := power.New(powerSensorDrv)

	// 6 - Kitchen Temperature Sensor
	tempSensorDrv := dht22.New(4)
	kitchenTempController := temperature.New(tempSensorDrv, 10*time.Minute, time.Minute)

	// 7 - LTEModule
	sim7600Drv := sim7600.New("/dev/ttyUSB2", 20*time.Second)
	lteDrv := ltemodule.New(sim7600Drv)
	lteController := ltemodulecontroller.New(lteDrv)

	// 8 - Business logic
	bl := businesslogic.New(roomLightController, kitchenTempController, powerController, time.Minute, time.Minute)
	bl.Start()

	// 9 - webserver
	components := &webserver.APIComponents{
		RoomLight:   roomLightController,
		CamLight:    camLightController,
		KitchenTemp: kitchenTempController,
		Power:       powerController,
		Boiler:      boilerController,
		LTEModule:   lteController,
	}

	ws := webserver.New(components, 8888)
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

	// TODO
	time.Sleep(1000 * time.Hour)
}
