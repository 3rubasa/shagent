package main

import (
	"fmt"
	"net/http"

	//"os"
	//"syscall"
	"time"

	//"bitbucket.org/gmcbay/i2c"
	"github.com/3rubasa/shagent/controllers/light"
	"github.com/3rubasa/shagent/controllers/watchdog"
	"github.com/3rubasa/shagent/sensors/power"
	"github.com/3rubasa/shagent/sensors/temperature"
)

const sampleInterval = time.Second * 60

// func SwitchRalayOn() {
// 	fmt.Println("opening gpio")
// 	err := rpio.Open()
// 	if err != nil {
// 		panic(fmt.Sprint("unable to open gpio", err.Error()))
// 	}

// 	defer rpio.Close()

// 	pin := rpio.Pin(26)
// 	pin.Output()

// 	for x := 0; x < 20; x++ {
// 		pin.Toggle()
// 		time.Sleep(time.Second * 5)
// 	}
// }

func main() {
	wd := watchdog.New()
	wd.Initialize()
	wd.Start()

	l := light.New()
	err := l.Initialize()
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

	tp := temperature.New()

	err = tp.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize temperature sensor: ", err)
		return
	}

	pp := power.New()

	err = pp.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize power sensor: ", err)
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

		SendMeasurements(t, p)
		time.Sleep(sampleInterval)
	}
}

func SendMeasurements(t float64, p int) error {
	//bodyReader := bytes.NewReader([]byte(body))
	url := fmt.Sprintf("https://api.thingspeak.com/update?api_key=TL9W7QIEFKFIYIS7&field1=%f&field2=%d", t, p)
	fmt.Printf("About to send request: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error while creating request: %s", err.Error())
		return err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %s", err.Error())
		return err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status is not 200: %d", resp.StatusCode)
		fmt.Printf("Error: %s", err.Error())
		return fmt.Errorf("response status is not 200: %d", resp.StatusCode)
	}

	return nil
}
