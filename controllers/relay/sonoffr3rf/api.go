package sonoffr3rf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type API struct {
	osSvcs  OSServicesProvider
	macAddr string
}

func New(osSvcs OSServicesProvider, macAddr string) *API {
	return &API{
		osSvcs:  osSvcs,
		macAddr: macAddr,
	}
}

func (a API) GetState() (string, error) {
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	ip, err := a.osSvcs.GetIPFromMAC(a.macAddr)
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

func (a API) TurnOn() error {
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	ip, err := a.osSvcs.GetIPFromMAC(a.macAddr)
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

func (a API) TurnOff() error {
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	ip, err := a.osSvcs.GetIPFromMAC(a.macAddr)
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
