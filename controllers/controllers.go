package controllers

type LightController interface {
	Initialize() error
	Start() error
	Stop()
}
