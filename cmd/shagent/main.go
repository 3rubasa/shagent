package main

import (
	"fmt"
	"log"

	//"os"
	//"syscall"
	"time"

	//"bitbucket.org/gmcbay/i2c"

	"github.com/3rubasa/osservices"
	"github.com/3rubasa/shagent/pkg/businesslogic"
	"github.com/3rubasa/shagent/pkg/businesslogic/controllers/forecastprovider/openweather"
	"github.com/3rubasa/shagent/pkg/businesslogic/controllers/ltemodulecontroller"
	"github.com/3rubasa/shagent/pkg/businesslogic/controllers/power"
	scheduledrelay "github.com/3rubasa/shagent/pkg/businesslogic/controllers/relay/scheduled"
	simplerelay "github.com/3rubasa/shagent/pkg/businesslogic/controllers/relay/simple"
	"github.com/3rubasa/shagent/pkg/businesslogic/controllers/temperature"
	"github.com/3rubasa/shagent/pkg/businesslogic/controllers/weatherprovider"
	"github.com/3rubasa/shagent/pkg/config"
	"github.com/3rubasa/shagent/pkg/drivers/ltemodule"
	"github.com/3rubasa/shagent/pkg/drivers/ltemodule/sim7600"
	"github.com/3rubasa/shagent/pkg/drivers/power/gpiopowersensor"
	"github.com/3rubasa/shagent/pkg/drivers/relay/asyncdecorator"
	"github.com/3rubasa/shagent/pkg/drivers/relay/sonoffr3rf"
	"github.com/3rubasa/shagent/pkg/drivers/relay/wsraspihatx3"
	"github.com/3rubasa/shagent/pkg/drivers/termo/dht22"
	"github.com/3rubasa/shagent/pkg/drivers/termo/esp8266"
	"github.com/3rubasa/shagent/pkg/grpcapi"
	"github.com/3rubasa/shagent/pkg/webserver"
)

const configPath = "./shagent.json"

func main() {
	log.Println("Starting shagent...")
	var err error

	log.Println("Reading configuration from file ", configPath, " ...")
	cfg, err := config.ReadFromFile(configPath)
	if err != nil {
		log.Println("ERROR: Failed to get config: ", err)
		return
	}

	// Common
	osservices := osservices.NewOSServicesProvider()

	// 2 - boiler DONE
	sonoffr3rfRelayDrv := sonoffr3rf.New(osservices, "24:a1:60:1d:72:9d", "10012ff7ab", "10.0.0.17")
	boilerRelayDrv := asyncdecorator.New(sonoffr3rfRelayDrv, 10*time.Minute)
	boilerController := simplerelay.New(boilerRelayDrv)

	// 3 - roomLight DONE
	// TODO: Later, if roomLight is nil, what are we going to do?
	roomLightRelayDrv, err := wsraspihatx3.New(wsraspihatx3.RelayChannel2)
	if err != nil {
		log.Println("ERROR: Failed to create WaveShare Raspi Hat relay device: ", err)
	}

	var ontimes, offtimes []string
	ontimes = append(ontimes, "0 45 06 * * *", "0 10 17 * * *")
	offtimes = append(offtimes, "0 15 08 * * *", "0 12 01 * * *")

	roomLightController := scheduledrelay.New(roomLightRelayDrv, ontimes, offtimes)

	// 4 - camLight DONE
	// TODO: Later, if cam light is nil, what are we going to do?
	camLightRelayDrv, err := wsraspihatx3.New(wsraspihatx3.RelayChannel3)
	if err != nil {
		fmt.Println("ERROR: Failed to create WaveShare Raspi Hat relay device: ", err)
	}

	camLightController := simplerelay.New(camLightRelayDrv)

	// 5 - power DONE
	powerSensorDrv := gpiopowersensor.New(16)
	powerController := power.New(powerSensorDrv)

	// 6.1 - Kitchen Temperature Sensor
	kitchenTempSensorDrv := dht22.New(4)
	kitchenTempController := temperature.New(kitchenTempSensorDrv, 10*time.Minute, time.Minute)

	// 6.2 - Window Temperature Sensor
	windowTempSensorDrv := dht22.New(24)
	windowTempController := temperature.New(windowTempSensorDrv, 10*time.Minute, time.Minute)

	// 6.3 - Outdoors Temperature Sensor
	outdoorsTempSensorDrv := dht22.New(23)
	outdoorsTempController := temperature.New(outdoorsTempSensorDrv, 10*time.Minute, time.Minute)

	pantryTempSensorDrv := esp8266.New(3333, 30*time.Minute)
	pantryTempController := temperature.New(pantryTempSensorDrv, 30*time.Minute, 10*time.Minute)

	// 6.4 - Weather Temperature
	weatherTempProvider := weatherprovider.New(&cfg.WeatherProvider, "5e9e1159073f45d7a7063db8891c82b0", "stebnyk", "ua", time.Minute, 45*time.Minute, 95*time.Minute)

	// 6.5 - Forecast Provider
	// 6.4 - Weather Temperature
	forecastTempProvider := openweather.New(&cfg.ForecastProvider, "3c49c423e46752496919cbef4186df81", "stebnyk", "ua", time.Minute, 30*time.Minute, 3*time.Hour)

	// 7 - LTEModule
	sim7600Drv := sim7600.New("/dev/ttyUSB2", 20*time.Second)
	lteDrv := ltemodule.New(sim7600Drv)
	lteController := ltemodulecontroller.New(lteDrv)

	blComponents := &businesslogic.BusinessLogicComponents{
		RoomLight:    roomLightController,
		CamLight:     camLightController,
		KitchenTemp:  kitchenTempController,
		WindowTemp:   windowTempController,
		ForecastTemp: forecastTempProvider,
		OutdoorsTemp: outdoorsTempController,
		PantryTemp:   pantryTempController,
		WeatherTemp:  weatherTempProvider,
		Power:        powerController,
		Boiler:       boilerController,
		LTEModule:    lteController,
	}

	tempControlTable := map[businesslogic.ForecestTempRange]map[businesslogic.WeatherTempRange]businesslogic.RoomTempRange{
		{Min: -100, Max: -4}: {
			{Min: -100, Max: 100}: {Min: 8, Max: 100}, // Don't turn off
		},
		{Min: -4, Max: 100}: {
			{Min: -100, Max: -3}: {Min: 8, Max: 100}, // Don't turn off
			{Min: -3, Max: 100}:  {Min: 6, Max: 8},
		},
	}

	// tempControlTable := map[businesslogic.ForecestTempRange]map[businesslogic.WeatherTempRange]businesslogic.RoomTempRange{
	// 	{Min: -100, Max: -4}: {
	// 		{Min: -100, Max: 100}: {Min: 8, Max: 100}, // Don't turn off
	// 	},
	// 	{Min: -4, Max: -3}: {
	// 		{Min: -100, Max: -3}: {Min: 8, Max: 100}, // Don't turn off
	// 		{Min: -3, Max: 100}:  {Min: 6, Max: 8},
	// 	},
	// 	{Min: -3, Max: -1}: {
	// 		{Min: -100, Max: -3}: {Min: 8, Max: 100}, // Don't turn off
	// 		{Min: -3, Max: 1}:    {Min: 5, Max: 6},
	// 		{Min: 1, Max: 100}:   {Min: 5, Max: 6},
	// 	},
	// 	{Min: -1, Max: 1}: {
	// 		{Min: -100, Max: -3}: {Min: 8, Max: 100}, // Don't turn off
	// 		{Min: -3, Max: 0}:    {Min: 5, Max: 6},
	// 		{Min: 0, Max: 5}:     {Min: 5, Max: 6},
	// 		{Min: 5, Max: 100}:   {Min: 4, Max: 5},
	// 	},
	// 	{Min: 1, Max: 100}: {
	// 		{Min: -100, Max: -3}: {Min: 8, Max: 100}, // Don't turn off
	// 		{Min: -3, Max: 0}:    {Min: 5, Max: 6},
	// 		{Min: 0, Max: 1}:     {Min: 4, Max: 5},
	// 		{Min: 1, Max: 100}:   {Min: 3, Max: 4},
	// 	},
	// }

	businessRules := &businesslogic.BusinessRules{
		TempControlTable: tempControlTable,
	}

	processor := businesslogic.NewProcessor(businessRules, blComponents)

	// 8 - Business logic
	bl := businesslogic.New(&cfg.BusinessLogic, blComponents, processor, 5*time.Minute, 2*time.Minute, 5*time.Minute)
	bl.Start()

	// 9 - webserver
	components := &webserver.APIComponents{
		RoomLight:   roomLightController,
		CamLight:    camLightController,
		KitchenTemp: kitchenTempController,
		Power:       powerController,
		Boiler:      boilerController,
		MC:          bl,
	}

	ws := webserver.New(components, 8888)
	err = ws.Initialize()
	if err != nil {
		log.Println("ERROR: Failed to initialize the web server: ", err)
	}

	err = ws.Start()
	if err != nil {
		log.Println("ERROR: Failed to start the web server: ", err)
	}

	ga := grpcapi.New(bl, 50053)
	err = ga.Start()
	if err != nil {
		log.Println("ERROR: Failed to start GRPC API: ", err)
	}

	log.Println("SERVICE INTIALIZATION COMPLETE")

	// TODO
	time.Sleep(1000 * time.Hour)
}
