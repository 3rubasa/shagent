package businesslogic

type BoilerController interface {
	Start() error
	Stop()
	// Returns "on" or "off"
	GetState() (string, error)
	TurnOn() error
	TurnOff() error
}
