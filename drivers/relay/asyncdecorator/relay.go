package asyncdecorator

import (
	"log"
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

	log.Println("Debug: Async Relay Decorator: OnNewJob()")

	// Get current state
	csStr, err := p.GetState()
	if err != nil {
		log.Printf("NOTICE: Could not to get relay current state: %v", err)
		return
	}

	cs, err := StringToRelayState(csStr)
	if err != nil {
		log.Printf("ERROR: Failed to convert relay state from string: %v", err)
		return
	}

	log.Printf("Debug: Async Relay's current state is %d, target state is %d", cs, ts)

	// Compare it with target state
	if cs == ts {
		return
	}

	// Apply target state if it is different from the current one
	switch ts {
	case relayStateOn:
		log.Println("Debug: Async Relay Decorator: turning the relay on...")
		err = p.deviceAPI.TurnOn()
		if err != nil {
			log.Printf("NOTICE: Could not turn the relay on: %v", err)
			return
		}
	case relayStateOff:
		log.Println("Debug: Async Relay Decorator: turning the relay off...")
		err = p.deviceAPI.TurnOff()
		if err != nil {
			log.Printf("NOTICE: Could not turn the relay off: %v", err)
			return
		}
	case relayStateNeutral:
		log.Println("Debug: Async Relay Decorator: skipping because target state is Neutral...")
		// Do nothing
	default:
		log.Panicf("ERROR: Async Relay Decorator: Unexpected target state of the relay: %d", ts)
		return
	}
}
