package sensors

type TemperatureProvider interface {
	Initialize() error
	GetTemperature() (float64, error)
}

type PowerStatusProvider interface {
	Initialize() error
	GetPowerStatus() (int, error) // 1 - on, 0 - off
}
