package sonoffr3rf

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/3rubasa/osservices"
)

type Driver struct {
	osSvcs  OSServicesProvider
	macAddr string
}

func New(osSvcs OSServicesProvider, macAddr string) *Driver {
	return &Driver{
		osSvcs:  osSvcs,
		macAddr: macAddr,
	}
}

func (d Driver) GetState() (string, error) {
	ip, err := d.osSvcs.GetIPFromMAC(d.macAddr)
	if err == osservices.ErrNotFound {
		log.Println("NOTICE: relay ", d.macAddr, " not available")
		return "", ErrNotAvailable
	} else if err != nil {
		log.Println("ERROR: Failed to get relay state: ", err)
		return "", fmt.Errorf("failed to get relay stay: %v", err)
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, relayPort, relayInfoPath)
	reqBody := strings.NewReader(relayInfoBody)

	log.Println("Debug: About to send http request to relay: ", url)
	// TODO: change to milliseconds
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Microsecond)

	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		log.Printf("ERROR: failed to create request: %v", err)
		return "", fmt.Errorf("failed to get relay state: %v", err)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: failed to send request to relay: %v", err)
		return "", ErrNotAvailable
	}

	if resp.StatusCode != 200 {
		log.Println("ERROR: relay ", d.macAddr, " replied with HTTP status not equal to 200: ", resp.StatusCode)
		return "", fmt.Errorf("relay replied with unexpected HTTP status code %d", resp.StatusCode)
	}

	var info RelayInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		log.Println("ERROR: can't parse response from relay: ", err)
		return "", fmt.Errorf("can't parse response from relay: %v", err)
	}

	if info.Error != 0 {
		log.Println("ERROR: relay replied with an error: ", info.Error)
		return "", fmt.Errorf("relay replied with an error: %d", info.Error)
	}

	if info.Data.Switch != "on" && info.Data.Switch != "off" {
		log.Printf("ERROR: Relay returned unexpected value for state: %s", info.Data.Switch)
		return "", fmt.Errorf("relay returned unexpected value for state: %s", info.Data.Switch)
	}

	return info.Data.Switch, nil
}

func (a Driver) TurnOn() error {
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

func (a Driver) TurnOff() error {
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
