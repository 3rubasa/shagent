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
		return "", fmt.Errorf("failed to get relay IP address: %v", err)
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, relayPort, relayInfoPath)
	reqBody := strings.NewReader(relayInfoBody)

	log.Println("Debug: About to send http request to relay: ", url)
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Millisecond)

	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		log.Printf("ERROR: failed to create request: %v", err)
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
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

func (d Driver) TurnOn() error {
	ip, err := d.osSvcs.GetIPFromMAC(d.macAddr)
	if err == osservices.ErrNotFound {
		log.Println("NOTICE: relay ", d.macAddr, " not available")
		return ErrNotAvailable
	} else if err != nil {
		log.Println("ERROR: Failed TURN ON the relay: ", err)
		return fmt.Errorf("failed to get relay IP address: %v", err)
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, relayPort, relaySwitchPath)

	reqBody := strings.NewReader(relaySwitchOnBody)

	log.Println("Debug: About to send http request to relay: ", url)
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Millisecond)
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		log.Printf("ERROR: failed to create request: %v", err)
		return fmt.Errorf("failed to create a HTTP request: %v", err)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: failed to send request to relay: %v", err)
		return ErrNotAvailable
	}

	if resp.StatusCode != 200 {
		log.Println("ERROR: relay ", d.macAddr, " replied with HTTP status not equal to 200: ", resp.StatusCode)
		return fmt.Errorf("relay replied with unexpected HTTP status code %d", resp.StatusCode)
	}

	var respBody RelaySwitchOnResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Println("ERROR: can't parse response from relay: ", err)
		return fmt.Errorf("can't parse response from relay: %v", err)
	}

	if respBody.Error != 0 {
		log.Println("ERROR: relay replied with an error: ", respBody.Error)
		return fmt.Errorf("relay replied with an error: %d", respBody.Error)
	}

	return nil
}

func (d Driver) TurnOff() error {
	ip, err := d.osSvcs.GetIPFromMAC(d.macAddr)
	if err == osservices.ErrNotFound {
		log.Println("NOTICE: relay ", d.macAddr, " not available")
		return ErrNotAvailable
	} else if err != nil {
		log.Println("ERROR: Failed to get the relay IP address: ", err)
		return fmt.Errorf("failed to get relay IP address: %v", err)
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, relayPort, relaySwitchPath)

	reqBody := strings.NewReader(relaySwitchOffBody)

	log.Println("Debug: About to send http request to relay: ", url)
	// Enforce delay - switch can process one request per 200 ms
	time.Sleep(200 * time.Millisecond)
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		log.Printf("ERROR: failed to create request: %v", err)
		return fmt.Errorf("failed to create an HTTP request: %v", err)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: failed to send request to relay: %v", err)
		return ErrNotAvailable
	}

	if resp.StatusCode != 200 {
		log.Println("ERROR: relay ", d.macAddr, " replied with HTTP status not equal to 200: ", resp.StatusCode)
		return fmt.Errorf("relay replied with unexpected HTTP status code %d", resp.StatusCode)
	}

	var respBody RelaySwitchOnResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Println("ERROR: can't parse response from relay: ", err)
		return fmt.Errorf("can't parse response from relay: %v", err)
	}

	if respBody.Error != 0 {
		log.Println("ERROR: relay replied with an error: ", respBody.Error)
		return fmt.Errorf("relay replied with an error: %d", respBody.Error)
	}

	return nil
}
