package asyncdecorator

type DeviceAPI interface {
	GetState() (string, error)
	TurnOn() error
	TurnOff() error
}
