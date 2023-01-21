package controllers

type RelayController interface {
	Start() error
	Stop()
	GetState() (string, error)
	TurnOn() error
	TurnOff() error
}

type LightController interface {
	Initialize() error
	Start() error
	Stop()
}
