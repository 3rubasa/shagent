package interfaces

type RelayDriver interface {
	Start() error
	Stop()
	// Returns "on" or "off"
	GetState() (string, error)
	TurnOn() error
	TurnOff() error
}

type TempSensorDriver interface {
	Start() error
	Stop() error
	Get() (float64, error)
}

type PowerDriver interface {
	Initialize() error
	Get() (int, error)
}

type LTEModuleDriver interface {
	SendUSSD(ussd string) (string, error)
}
