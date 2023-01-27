package watchdog

// Improvemens:
// 1. The application needs to be normally shut down before restart
// 2. Configuration should be provided as a struct (read from toml?)

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

type InternetChecker interface {
	IsInternetAvailable() (bool, error)
}

type Watchdog struct {
	vpnTicker         *time.Ticker
	vpnPeriod         time.Duration
	inetTicker        *time.Ticker
	inetPeriod        time.Duration
	inetPeriodOnError time.Duration
	done              chan bool
	inetErrCount      uint8
	osservices        OSServicesProvider
	inetchecker       InternetChecker
}

func New(osservices OSServicesProvider, inetchecker InternetChecker, inetPeriod, inetPeriodOnError, vpnPeriod time.Duration) *Watchdog {
	return &Watchdog{
		done:              make(chan bool),
		osservices:        osservices,
		inetchecker:       inetchecker,
		inetPeriod:        inetPeriod,
		inetPeriodOnError: inetPeriodOnError,
		vpnPeriod:         vpnPeriod,
	}
}

func (p *Watchdog) Start() error {
	p.inetTicker = time.NewTicker(p.inetPeriod)
	p.vpnTicker = time.NewTicker(p.vpnPeriod)

	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.inetTicker.C:
				p.testInternetConnection()
			case <-p.vpnTicker.C:
				p.testVPNConnection()
			}
		}
	}()

	return nil
}

func (p *Watchdog) Stop() {
	p.inetTicker.Stop()
	p.done <- true
}

func (p *Watchdog) testInternetConnection() {
	var err error

	// Check internet connection
	internetIsAvail, err := p.inetchecker.IsInternetAvailable()
	if err != nil {
		fmt.Printf("Failed to check internet availability: %s. Rebooting immediately.\n", err.Error())
		err = p.osservices.Reboot()
		if err != nil {
			fmt.Println("Reboot failed: ", err)
		} else {
			// For testing purposes
			p.Stop()
		}

		return
	}

	if internetIsAvail {
		p.inetErrCount = 0
		p.inetTicker.Reset(p.inetPeriod)
	} else {
		p.inetErrCount++
		p.inetTicker.Reset(p.inetPeriodOnError)
	}

	if p.inetErrCount >= 3 {
		fmt.Printf("Internet is not available for the 3rd time in a row, issuing reboot command \n")
		err = p.osservices.Reboot()
		if err != nil {
			fmt.Println("Failed to reboot: ", err)
		} else {
			// For testing purposes
			p.Stop()
		}
	}
}

func (p *Watchdog) testVPNConnection() {
	ifcs, err := net.Interfaces()
	if err != nil {
		fmt.Println("Failed to get the list of interfaces: ", err)
		goto reboot
	}

	for _, ifc := range ifcs {
		if strings.HasPrefix(ifc.Name, "tun") {
			return
		}
	}

	// If we are here it means we could not find interface named tun0, which means that
	// no VPN connection is set up. Restart openvpn service

	fmt.Println("Watchdog: interface with name tun0 was not found, restarting openvpn service...")

	err = p.restartVPNService()
	if err != nil {
		fmt.Println("Failed to restart openvpn service: ", err)
		goto reboot
	}

reboot:
	err = p.osservices.Reboot()
	if err != nil {
		fmt.Println("Failed to reboot: ", err)
	} else {
		// For testing purposes
		p.Stop()
	}
}

func (p *Watchdog) restartVPNService() error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	command := "sudo systemctl restart openvpn"
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println("Failed to exec command ", command, " Error:", err)
		return err
	}

	if len(stderr.String()) > 0 {
		fmt.Println("STDERR: ", stderr.String())
		return fmt.Errorf("error, STDERR not empty: %s", stderr.String())
	}

	return nil
}
