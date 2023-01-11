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
