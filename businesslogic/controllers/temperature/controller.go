package temperature

import (
	"fmt"
	"sync"
	"time"

	"github.com/3rubasa/shagent/businesslogic/interfaces"
)

type Controller struct {
	driver           interfaces.TempSensorDriver
	ticker           *time.Ticker
	done             chan bool
	curTemp          float64
	curTempTimestamp time.Time
	cacheTTL         time.Duration
	mux              *sync.Mutex
	pollingPeriod    time.Duration
}

func New(driver interfaces.TempSensorDriver, cacheTTL, pollingPeriod time.Duration) *Controller {
	return &Controller{
		driver:           driver,
		done:             make(chan bool),
		curTemp:          invalidTemperature,
		curTempTimestamp: time.Date(1800, 1, 1, 0, 0, 0, 0, time.Local),
		cacheTTL:         cacheTTL,
		mux:              &sync.Mutex{},
		pollingPeriod:    pollingPeriod,
	}
}

func (c *Controller) Start() error {
	err := c.driver.Initialize()
	if err != nil {
		fmt.Println("Error while starting temperature sensor controller: ", err)
		return err
	}

	c.ticker = time.NewTicker(c.pollingPeriod)

	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-c.ticker.C:
				c.onTick()
			}
		}
	}()

	return nil
}

func (c *Controller) Stop() error {
	c.ticker.Stop()
	c.done <- true

	return nil
}

func (c *Controller) Get() (float64, error) {
	c.mux.Lock()
	t := c.curTemp
	ts := c.curTempTimestamp
	c.mux.Unlock()

	if time.Since(ts) <= c.cacheTTL {
		return t, nil
	}

	return invalidTemperature, fmt.Errorf("up to date temperature is not available")
}

func (c *Controller) onTick() {
	var err error
	// Try to get cur temperture from the driver
	ts := time.Now()
	t, err := c.driver.Get()
	if err != nil {
		fmt.Println("Failed to get temperature: ", err)
		return
	}

	// If succeeded
	// 	store the got temperature and ts atomically
	c.mux.Lock()
	c.curTemp = t
	c.curTempTimestamp = ts
	c.mux.Unlock()
}
