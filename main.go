package main

import (
	"fmt"
	"net/http"

	//"os"
	//"syscall"
	"time"

	//"bitbucket.org/gmcbay/i2c"

	"github.com/3rubasa/shagent/controllers/light"
	"github.com/3rubasa/shagent/controllers/relay"
	"github.com/3rubasa/shagent/controllers/relay/sonoffr3rf"
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

	// 2 - boiler
	relayDevice := sonoffr3rf.New(osservices, "24:a1:60:1d:72:9d")
	// TODO: set proper period
	b := relay.New(relayDevice, 10*time.Second)
	err = b.Start()
	if err != nil {
		fmt.Println("Failed to start boiler relay controller: ", err)
	}

	// 3 - light
	l := light.New()
	err = l.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize light controller: ", err)
		return
	}

	err = l.Start()
	if err != nil {
		fmt.Println("Failed to start light controller: ", err)
		return
	}
	defer l.Stop()

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
	ws := webserver.New(b)
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
		bsStr, err := b.GetState()
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
