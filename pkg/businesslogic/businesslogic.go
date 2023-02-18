package businesslogic

import (
	"fmt"
	"log"
	"net/http"
	urltools "net/url"
	"sync"
	"time"

	"github.com/3rubasa/shagent/pkg/config"
)

type BusinessLogic struct {
	c                *BusinessLogicComponents
	p                *Processor
	cfg              *config.BusinessLogicConfig
	done             chan bool
	pollingPeriod    time.Duration
	pollingTicker    *time.Ticker
	sendingPeriod    time.Duration
	sendingTicker    *time.Ticker
	processingPeriod time.Duration
	processingTicker *time.Ticker
	stateMux         *sync.Mutex
	state            State
}

func New(cfg *config.BusinessLogicConfig, c *BusinessLogicComponents, p *Processor, pollingPeriod, sendingPeriod, processingPeriod time.Duration) *BusinessLogic {
	return &BusinessLogic{
		c:                c,
		p:                p,
		done:             make(chan bool),
		stateMux:         &sync.Mutex{},
		pollingPeriod:    pollingPeriod,
		sendingPeriod:    sendingPeriod,
		processingPeriod: processingPeriod,
		cfg:              cfg,
	}
}

func (b *BusinessLogic) Start() error {
	var err error

	err = b.c.KitchenTemp.Start()
	if err != nil {
		fmt.Println("Failed to start kitchen temperature sensor controller: ", err)
	}

	err = b.c.WindowTemp.Start()
	if err != nil {
		fmt.Println("Failed to start window temperature sensor controller: ", err)
	}

	err = b.c.OutdoorsTemp.Start()
	if err != nil {
		fmt.Println("Failed to start outdoors temperature sensor controller: ", err)
	}

	err = b.c.PantryTemp.Start()
	if err != nil {
		fmt.Println("Failed to start pantry temperature sensor controller: ", err)
	}

	err = b.c.WeatherTemp.Start()
	if err != nil {
		fmt.Println("Failed to start weather temperature sensor controller: ", err)
	}

	err = b.c.ForecastTemp.Start()
	if err != nil {
		fmt.Println("Failed to start weathr forecast provider controller: ", err)
	}

	err = b.c.RoomLight.Start()
	if err != nil {
		fmt.Println("Failed to start room light relay controller: ", err)
	}

	err = b.c.CamLight.Start()
	if err != nil {
		fmt.Println("Failed to start cam light relay controller: ", err)
	}

	err = b.c.Power.Start()
	if err != nil {
		fmt.Println("Failed to start power sensor controller: ", err)
	}

	err = b.c.Boiler.Start()
	if err != nil {
		fmt.Println("Failed to start boiler controller: ", err)
	}

	b.pollingTicker = time.NewTicker(b.pollingPeriod)
	b.sendingTicker = time.NewTicker(b.sendingPeriod)
	b.processingTicker = time.NewTicker(b.processingPeriod)

	go b.mainLoop()

	return nil
}

func (b *BusinessLogic) Stop() error {
	b.done <- true
	return nil
}

func (b *BusinessLogic) GetState() State {
	b.stateMux.Lock()
	s := b.state
	b.stateMux.Unlock()

	return s
}

func (b *BusinessLogic) mainLoop() error {

	for {
		select {
		case <-b.pollingTicker.C:
			b.pollSensors()
		case <-b.sendingTicker.C:
			b.sendState()
		case <-b.processingTicker.C:
			b.processState()
		case <-b.done:
			return nil
		}
	}
}

func (b *BusinessLogic) pollSensors() {
	fmt.Println("+ pollSensors started")
	defer fmt.Println("+ pollSensors exited")
	var t float64
	var err error
	var s State

	// RoomLight
	rl, err := b.c.RoomLight.Get()
	if err != nil {
		fmt.Println("Failed to get room light state: ", err)
	} else {
		s.RoomLightState = rl
		s.RoomLightStateValid = true
	}

	t, err = b.c.KitchenTemp.Get()
	if err != nil {
		fmt.Println("Error while getting kitchen temperature: ", err)
	} else {
		s.KitchenTemp = t
		s.KitchenTempValid = true
	}

	t, err = b.c.WindowTemp.Get()
	if err != nil {
		fmt.Println("Error while getting window temperature: ", err)
	} else {
		s.WindowTemp = t
		s.WindowTempValid = true
	}

	t, err = b.c.PantryTemp.Get()
	if err != nil {
		log.Println("NOTICE: Could not get pantry temperature: ", err)
	} else {
		s.PantryTemp = t
		s.PantryTempValid = true
	}

	t, err = b.c.WeatherTemp.Get()
	if err != nil {
		fmt.Println("Error while getting weather temperature: ", err)
	} else {
		s.WeatherTemp = t
		s.WeatherTempValid = true
	}

	t, err = b.c.ForecastTemp.Get()
	if err != nil {
		fmt.Println("Error while getting forecast temperature: ", err)
	} else {
		s.ForecastedTemp = t
		s.ForecastedTempValid = true
	}

	// t, err = b.c.OutdoorsTemp.Get()
	// if err != nil {
	// 	fmt.Println("Error while outdoors window temperature: ", err)
	// } else {
	// 	s.OutdoorsTemp = t
	// 	s.OutdoorsTempValid = true
	// }

	p, err := b.c.Power.Get()
	if err != nil {
		fmt.Println("Error while getting power status: ", err)
	} else {
		s.Power = p
		s.PowerValid = true
	}

	// Boiler
	bs, err := b.c.Boiler.Get()
	if err != nil {
		fmt.Println("Failed to get boiler state: ", err)
	} else {
		s.BoilerState = bs
		s.BoilerStateValid = true
	}

	b.stateMux.Lock()
	b.state = s
	b.stateMux.Unlock()
}

func (b *BusinessLogic) sendState() {
	fmt.Println("+ SendState started")
	defer fmt.Println("+ SendState exited")

	var s State
	b.stateMux.Lock()
	s = b.state
	b.stateMux.Unlock()

	url := fmt.Sprintf("%s://%s/%s?api_key=%s", b.cfg.Consumer.Schema, b.cfg.Consumer.Address, b.cfg.Consumer.URI, b.cfg.Consumer.APIKeys)

	if s.KitchenTempValid {
		url += "&field1=" + urltools.QueryEscape(fmt.Sprintf("%f", s.KitchenTemp))
	}
	if s.WindowTempValid {
		url += "&field2=" + urltools.QueryEscape(fmt.Sprintf("%f", s.WindowTemp))
	}
	if s.PantryTempValid {
		url += "&field3=" + urltools.QueryEscape(fmt.Sprintf("%f", s.PantryTemp))
	}
	if s.PowerValid {
		url += "&field4=" + urltools.QueryEscape(fmt.Sprintf("%d", s.Power))
	}
	if s.BoilerStateValid {
		url += "&field5=" + urltools.QueryEscape(fmt.Sprintf("%d", s.BoilerState))
	}
	if s.WeatherTempValid {
		url += "&field6=" + urltools.QueryEscape(fmt.Sprintf("%f", s.WeatherTemp))
	}
	if s.ForecastedTempValid {
		url += "&field7=" + urltools.QueryEscape(fmt.Sprintf("%f", s.ForecastedTemp))
	}
	if s.RoomLightStateValid {
		url += "&field8=" + urltools.QueryEscape(fmt.Sprintf("%d", s.RoomLightState))
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

func (b *BusinessLogic) processState() {
	fmt.Println("+ processState started")
	defer fmt.Println("+ processState exited")
	b.stateMux.Lock()
	s := b.state
	b.stateMux.Unlock()

	b.p.Process(s)
}
