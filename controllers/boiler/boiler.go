package boiler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/3rubasa/shagent/controllers"
)

const (
	relaymDNSName = "ewelink_10012ff7ab.local" // "10.42.0.214"
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
}

var controllerSingleton *boilerController

func New() controllers.BoilerController {
	if controllerSingleton == nil {
		controllerSingleton = &boilerController{}
	}

	return controllerSingleton
}

func (p *boilerController) Initialize() error {
	return nil
}

func (p *boilerController) Start() error {
	return nil
}

func (p *boilerController) Stop() {
}

func (p *boilerController) GetState() (string, error) {
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	url := fmt.Sprintf("http://%s:%d/%s", relaymDNSName, relayPort, relayInfoPath)
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
	return nil
}
func (p *boilerController) TurnOff() error {
	return nil
}
