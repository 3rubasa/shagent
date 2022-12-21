package sensors

type TemperatureProvider interface {
	Initialize() error
	GetTemperature() (float64, error)
}
