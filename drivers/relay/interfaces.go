package relay

type DeviceAPI interface {
	GetState() (string, error)
	TurnOn() error
	TurnOff() error
}
