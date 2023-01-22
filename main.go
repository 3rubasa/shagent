package main

import (
	"fmt"
	"net/http"

	//"os"
	//"syscall"
	"time"

	//"bitbucket.org/gmcbay/i2c"

	"github.com/3rubasa/shagent/controllers/relay"
	"github.com/3rubasa/shagent/controllers/relay/sonoffr3rf"
	"github.com/3rubasa/shagent/controllers/relay/wsraspihatx3"
	"github.com/3rubasa/shagent/controllers/watchdog"
	"github.com/3rubasa/shagent/osservices"
	"github.com/3rubasa/shagent/sensors/power"
	"github.com/3rubasa/shagent/sensors/temperature"
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
	err = boiler.Start()
	if err != nil {
		fmt.Println("Failed to start boiler relay controller: ", err)
		// TODO: Later, if boiler has failed to start, what are we going to do?
	}

	// 3 - roomLight
	// TODO: Later, if roomLight is nil, what are we going to do?
	var roomLight *relay.Relay
	wsRelayForRoomLight, err := wsraspihatx3.New(wsraspihatx3.RelayChannel1)
	if err != nil {
		fmt.Println("Failed to create WaveShare Raspi Hat relay device: ", err)
	} else {
		roomLight = relay.New(wsRelayForRoomLight, 30*time.Minute)
		err = roomLight.Start()
		if err != nil {
			fmt.Println("Failed to start boiler relay controller: ", err)
		}
	}

	// 4 - camLight
	// TODO: Later, if cam light is nil, what are we going to do?
	var camLight *relay.Relay
	wsRelayForCamLight, err := wsraspihatx3.New(wsraspihatx3.RelayChannel2)
	if err != nil {
		fmt.Println("Failed to create WaveShare Raspi Hat relay device for cam light: ", err)
	} else {
		camLight = relay.New(wsRelayForCamLight, 30*time.Minute)
		err = camLight.Start()
		if err != nil {
			fmt.Println("Failed to start cam light relay controller: ", err)
		}
	}

	// 4 - temperature
	tp := temperature.New()

	err = tp.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize temperature sensor: ", err)
		return
	}

	// 5 - power
	pp := power.New()

	err = pp.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize power sensor: ", err)
		return
	}

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

	for {
		var t float64

		t, err := tp.GetTemperature()
		if err != nil {
			fmt.Println("Error while getting temperature: ", err)
		}

		fmt.Println("Current temperature is ", t)

		var p int

		p, err = pp.GetPowerStatus()
		if err != nil {
			fmt.Println("Error while getting power status: ", err)
		}

		// Boiler
		bs := -1
		bsStr, err := boiler.GetState()
		if err != nil {
			fmt.Println("Failed to get boiler state: ", err)
			bs = -1
		} else {
			switch bsStr {
			case "on":
				bs = 1
			case "off":
				bs = 0
			default:
				fmt.Println("Error invalid boiler state: ", bsStr)
				bs = -1
			}
		}

		SendMeasurements(t, p, bs)

		time.Sleep(sampleInterval)
	}
}

func SendMeasurements(t float64, p int, bs int) error {
	//bodyReader := bytes.NewReader([]byte(body))
	url := fmt.Sprintf("https://api.thingspeak.com/update?api_key=TL9W7QIEFKFIYIS7&field1=%f&field2=%d&field3=%d", t, p, bs)
	fmt.Printf("About to send request: %s \n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error while creating request: %s \n", err.Error())
		return err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %s \n", err.Error())
		return err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status is not 200: %d", resp.StatusCode)
		fmt.Printf("Error: %s \n", err.Error())
		return fmt.Errorf("response status is not 200: %d", resp.StatusCode)
	}

	return nil
}
