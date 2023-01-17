package boiler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/3rubasa/shagent/controllers"
)

const (
	// TODO: set proper period
	period = 10 * time.Second //10 * time.Minute
	//relaymDNSName = "eWeLink_10012ff7ab.local" // "10.42.0.214" // Fails to resolve for some reason
	relayMACAddr  = "24:a1:60:1d:72:9d"
	relayPort     = 8081
	relayInfoPath = "zeroconf/info"
	relayInfoBody = `{ 
		"deviceid": "", 
		"data": { } 
	 }`

	relaySwitchPath   = "zeroconf/switch"
	relaySwitchOnBody = `{ 
		"deviceid": "", 
		"data": {
			"switch": "on" 
		} 
	 }`

	relaySwitchOffBody = `{ 
		"deviceid": "", 
		"data": {
			"switch": "off" 
		} 
	 }`
)

type boilerController struct {
	targetState RelayState // "on", "off" or "neutral"
	ticker      *time.Ticker
	done        chan bool
	mux         *sync.Mutex
}

var controllerSingleton *boilerController

func New() controllers.BoilerController {
	if controllerSingleton == nil {
		controllerSingleton = &boilerController{}
	}

	return controllerSingleton
}

func (p *boilerController) Initialize() error {
	p.done = make(chan bool)
	p.mux = &sync.Mutex{}
	p.targetState = relayStateNeutral

	return nil
}

func (p *boilerController) Start() error {
	p.ticker = time.NewTicker(period)
	go p.MainLoop()
	return nil
}

func (p *boilerController) Stop() {
	p.done <- true
}

func (p *boilerController) GetState() (string, error) {
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	ip, err := GetIPFromMAC(relayMACAddr)
	if err != nil {
		fmt.Println("Failed to get relay IP address: ", err)
		return "", err
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, relayPort, relayInfoPath)
	fmt.Printf("About to send http request to relay: %s \n", url)

	reqBody := strings.NewReader(relayInfoBody)
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		fmt.Printf("Error while creating request: %s \n", err.Error())
		return "", err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %s \n", err.Error())
		return "", err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status is not 200: %d", resp.StatusCode)
		fmt.Printf("Error: %s \n", err.Error())
		return "", fmt.Errorf("response status is not 200: %d", resp.StatusCode)
	}

	var info RelayInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		fmt.Printf("Error while parsing response from relay: %s \n", err.Error())
		return "", err
	}

	if info.Error != 0 {
		err = fmt.Errorf("relay returned an error: %d", info.Error)
		fmt.Println(err.Error())
		return "", err
	}

	if info.Data.Switch != "on" && info.Data.Switch != "off" {
		err = fmt.Errorf("relay returned unexpected value for state: %s", info.Data.Switch)
		fmt.Println(err.Error())
		return "", err
	}

	return info.Data.Switch, nil
}

func (p *boilerController) TurnOn() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.targetState = relayStateOn

	// Try once to set the state synchronously
	return p.TurnOnInternal()
}

func (p *boilerController) TurnOff() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.targetState = relayStateOff

	// Try once to set the state synchronously
	return p.TurnOffInternal()
}

func (p *boilerController) TurnOnInternal() error {
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	ip, err := GetIPFromMAC(relayMACAddr)
	if err != nil {
		fmt.Println("Failed to get relay IP address: ", err)
		return err
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, relayPort, relaySwitchPath)
	fmt.Printf("About to send http request to relay: %s \n", url)

	reqBody := strings.NewReader(relaySwitchOnBody)
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
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

	var respBody RelaySwitchOnResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Printf("Error while parsing response from relay: %s \n", err.Error())
		return err
	}

	if respBody.Error != 0 {
		err = fmt.Errorf("relay returned an error: %d", respBody.Error)
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (p *boilerController) TurnOffInternal() error {
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	ip, err := GetIPFromMAC(relayMACAddr)
	if err != nil {
		fmt.Println("Failed to get relay IP address: ", err)
		return err
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, relayPort, relaySwitchPath)
	fmt.Printf("About to send http request to relay: %s \n", url)

	reqBody := strings.NewReader(relaySwitchOffBody)
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
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

	var respBody RelaySwitchOffResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Printf("Error while parsing response from relay: %s \n", err.Error())
		return err
	}

	if respBody.Error != 0 {
		err = fmt.Errorf("relay returned an error: %d", respBody.Error)
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (p *boilerController) MainLoop() {
	for {
		select {
		case <-p.done:
			return
		case <-p.ticker.C:
			p.OnTickerEvent()
		}
	}
}

func (p *boilerController) OnTickerEvent() {
	p.mux.Lock()
	defer p.mux.Unlock()

	fmt.Println("OnTickerEvent")
	// Get current state
	csStr, err := p.GetState()
	if err != nil {
		fmt.Printf("Failed to get current state: %s\n", err.Error())
		return
	}

	cs, err := StringToRelayState(csStr)
	if err != nil {
		fmt.Printf("Failed to convert relay state from string: %s\n", err.Error())
		return
	}

	fmt.Printf("Boiler's current state is %d, target state is %d", cs, p.targetState)

	// Compare it with target state
	if cs == p.targetState {
		return
	}

	ts := p.targetState

	// Apply target state if it is different from the current one
	switch ts {
	case relayStateOn:
		fmt.Printf("Trying to turn boiler on... \n")
		err = p.TurnOnInternal()
		if err != nil {
			fmt.Printf("Failed to turn boiler on: %s\n", err.Error())
			return
		}
	case relayStateOff:
		fmt.Printf("Trying to turn boiler off... \n")
		err = p.TurnOffInternal()
		if err != nil {
			fmt.Printf("Failed to turn boiler off: %s\n", err.Error())
			return
		}
	case relayStateNeutral:
		// Do nothing
	default:
		fmt.Printf("Error! Unexpected target state for boiler: %d\n", p.targetState)
	}
}
