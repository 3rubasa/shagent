package drivers

type Relay interface {
	Start() error
	Stop()
	GetState() (string, error)
	TurnOn() error
	TurnOff() error
}
