package businesslogic

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type BusinessLogic struct {
	//boiler    interfaces.RelayDriver
	roomLight RoomLightController
	//camLight  interfaces.RelayDriver
	kitchenTemp   TempSensorController
	power         PowerSensorController
	done          chan bool
	pollingPeriod time.Duration
	pollingTicker *time.Ticker
	sendingPeriod time.Duration
	sendingTicker *time.Ticker
	stateMux      *sync.Mutex
	state         state
}

func New(roomLight RoomLightController, kitchenTemp TempSensorController, powerController PowerSensorController, pollingPeriod, sendingPeriod time.Duration) *BusinessLogic {
	return &BusinessLogic{
		//boiler:    boiler,
		roomLight: roomLight,
		//camLight:  camLight,
		kitchenTemp:   kitchenTemp,
		power:         powerController,
		done:          make(chan bool),
		stateMux:      &sync.Mutex{},
		pollingPeriod: pollingPeriod,
		sendingPeriod: sendingPeriod,
	}
}

func (b *BusinessLogic) Start() error {
	var err error

	// err = b.boiler.Start()
	// if err != nil {
	// 	fmt.Println("Failed to start boiler relay: ", err)
	// }

	err = b.roomLight.Start()
	if err != nil {
		fmt.Println("Failed to start room light relay controller: ", err)
	}

	// err = b.camLight.Start()
	// if err != nil {
	// 	fmt.Println("Failed to start cam light relay controller: ", err)
	// }

	err = b.kitchenTemp.Start()
	if err != nil {
		fmt.Println("Failed to start kitchen temperature sensor controller: ", err)
	}

	err = b.power.Start()
	if err != nil {
		fmt.Println("Failed to start power sensor controller: ", err)
	}

	b.pollingTicker = time.NewTicker(b.pollingPeriod)
	b.sendingTicker = time.NewTicker(b.sendingPeriod)

	go b.mainLoop()

	return nil
}

func (b *BusinessLogic) Stop() error {
	b.done <- true
	return nil
}

func (b *BusinessLogic) mainLoop() error {

	for {
		select {
		case <-b.pollingTicker.C:
			b.pollSensors()
		case <-b.sendingTicker.C:
			b.sendState()
		case <-b.done:
			return nil
		}
	}
}

func (b *BusinessLogic) pollSensors() {
	var t float64

	var s state
	t, err := b.kitchenTemp.Get()
	if err != nil {
		fmt.Println("Error while getting kitchen temperature: ", err)
	} else {
		s.KitchenTemp = t
		s.KitchenTempValid = true
	}

	p, err := b.power.Get()
	if err != nil {
		fmt.Println("Error while getting power status: ", err)
	} else {
		s.Power = p
		s.PowerValid = true
	}

	// // Boiler
	// bs := -1
	// bsStr, err := boiler.GetState()
	// if err != nil {
	// 	fmt.Println("Failed to get boiler state: ", err)
	// 	bs = -1
	// } else {
	// 	switch bsStr {
	// 	case "on":
	// 		bs = 1
	// 	case "off":
	// 		bs = 0
	// 	default:
	// 		fmt.Println("Error invalid boiler state: ", bsStr)
	// 		bs = -1
	// 	}
	// }

	b.stateMux.Lock()
	b.state = s
	b.stateMux.Unlock()
}

func (b *BusinessLogic) sendState() {
	var s state
	b.stateMux.Lock()
	s = b.state
	b.stateMux.Unlock()

	url := "https://api.thingspeak.com/update?api_key=TL9W7QIEFKFIYIS7"
	if s.KitchenTempValid {
		url += fmt.Sprintf("&field1=%f", s.KitchenTemp)
	}
	if s.PowerValid {
		url += fmt.Sprintf("&field2=%d", s.Power)
	}

	fmt.Printf("About to send request: %s \n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error while creating request: %s \n", err.Error())
		return
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %s \n", err.Error())
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status is not 200: %d", resp.StatusCode)
		fmt.Printf("Error: %s \n", err.Error())
		return
	}
}
