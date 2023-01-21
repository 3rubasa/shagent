package watchdog

// Improvemens:
// 1. The application needs to be normally shut down before restart
// 2. Configuration should be provided as a struct (read from toml?)

import (
	"fmt"
	"time"
)

type InternetChecker interface {
	IsInternetAvailable() (bool, error)
}

type Watchdog struct {
	ticker      *time.Ticker
	period      time.Duration
	done        chan bool
	errCount    uint8
	osservices  OSServicesProvider
	inetchecker InternetChecker
}

func New(osservices OSServicesProvider, inetchecker InternetChecker, period time.Duration) *Watchdog {
	return &Watchdog{
		done:        make(chan bool),
		osservices:  osservices,
		inetchecker: inetchecker,
		period:      period,
	}
}

func (p *Watchdog) Start() error {
	p.ticker = time.NewTicker(p.period)
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				p.testInternetConnection()
			}
		}
	}()

	return nil
}

func (p *Watchdog) Stop() {
	p.ticker.Stop()
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
		p.errCount = 0
	} else {
		p.errCount++
	}

	if p.errCount >= 3 {
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
