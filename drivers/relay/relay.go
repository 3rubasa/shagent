package relay

import (
	"fmt"
	"sync"
	"time"
)

type Relay struct {
	targetState RelayState // "on", "off" or "neutral"
	ticker      *time.Ticker
	done        chan bool
	job         chan bool
	mux         *sync.Mutex
	period      time.Duration
	deviceAPI   DeviceAPI
}

func New(deviceAPI DeviceAPI, period time.Duration) *Relay {
	return &Relay{
		targetState: relayStateNeutral,
		done:        make(chan bool),
		job:         make(chan bool, 1),
		mux:         &sync.Mutex{},
		period:      period,
		deviceAPI:   deviceAPI,
	}
}

func (p *Relay) Start() error {
	p.ticker = time.NewTicker(p.period)
	go p.mainLoop()
	return nil
}

func (p *Relay) Stop() {
	p.ticker.Stop()
	p.done <- true
}

func (p *Relay) GetState() (string, error) {
	return p.deviceAPI.GetState()
}

func (p *Relay) TurnOn() error {
	p.mux.Lock()
	p.targetState = relayStateOn
	if len(p.job) == 0 {
		p.job <- true
	}
	p.mux.Unlock()

	return nil
}

func (p *Relay) TurnOff() error {
	p.mux.Lock()
	p.targetState = relayStateOff
	if len(p.job) == 0 {
		p.job <- true
	}
	p.mux.Unlock()

	return nil
}

func (p *Relay) mainLoop() {
	for {
		select {
		case <-p.done:
			return
		case <-p.ticker.C:
			p.mux.Lock()
			if len(p.job) == 0 {
				p.job <- true
			}
			p.mux.Unlock()
		case <-p.job:
			p.onNewJob()
		}
	}
}

func (p *Relay) onNewJob() {
	p.mux.Lock()
	ts := p.targetState
	p.mux.Unlock()

	fmt.Println("OnNewJob")
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

	fmt.Printf("Boiler's current state is %d, target state is %d", cs, ts)

	// Compare it with target state
	if cs == ts {
		return
	}

	// Apply target state if it is different from the current one
	switch ts {
	case relayStateOn:
		fmt.Printf("Trying to turn boiler on... \n")
		err = p.deviceAPI.TurnOn()
		if err != nil {
			fmt.Printf("Failed to turn boiler on: %s\n", err.Error())
			return
		}
	case relayStateOff:
		fmt.Printf("Trying to turn boiler off... \n")
		err = p.deviceAPI.TurnOff()
		if err != nil {
			fmt.Printf("Failed to turn boiler off: %s\n", err.Error())
			return
		}
	case relayStateNeutral:
		// Do nothing
	default:
		fmt.Printf("Error! Unexpected target state for boiler: %d\n", ts)
	}
}
