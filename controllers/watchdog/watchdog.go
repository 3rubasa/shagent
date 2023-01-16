package watchdog

import (
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/3rubasa/shagent/controllers"
)

const period = 30 * time.Minute
const url = "https://google.com"

type watchdog struct {
	ticker   *time.Ticker
	done     chan bool
	errCount uint8
}

var watchdogSingleton *watchdog

func New() controllers.Watchdog {
	if watchdogSingleton == nil {
		watchdogSingleton = &watchdog{}
	}

	return watchdogSingleton
}

func (p *watchdog) Initialize() error {
	p.done = make(chan bool)

	return nil
}

func (p *watchdog) Start() error {
	p.ticker = time.NewTicker(period)
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				p.TestInternetConnection()
			}
		}
	}()

	return nil
}

func (p *watchdog) Stop() {
	p.ticker.Stop()
	p.done <- true
}

func InternetIsAvailable() bool {
	fmt.Printf("About to send request to check if Internet is available: %s \n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error while creating request: %s \n", err.Error())
		return false
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %s \n", err.Error())
		return false
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("response status is >= 400: %d", resp.StatusCode)
		fmt.Printf("Error: %s \n", err.Error())
		return false
	}

	return true
}
func (p *watchdog) TestInternetConnection() {
	// Check internet connection
	internetIsAvail := InternetIsAvailable()
	if internetIsAvail {
		p.errCount = 0
	} else {
		p.errCount++
	}

	if p.errCount >= 3 {
		fmt.Printf("Internet is not available for the 3rd time in a row, issuing reboot command \n")
		syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	}
}
