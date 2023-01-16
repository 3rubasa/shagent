package controllers

type LightController interface {
	Initialize() error
	Start() error
	Stop()
}

type Watchdog interface {
	Initialize() error
	Start() error
	Stop()
}

type BoilerController interface {
	Initialize() error
	Start() error
	Stop()
	// Returns "on" or "off"
	GetState() (string, error)
	TurnOn() error
	TurnOff() error
}
