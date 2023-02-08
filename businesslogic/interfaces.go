package businesslogic

type RoomLightController interface {
	Start() error
	Stop() error
	Get() (int, error)
}

type TempSensorController interface {
	Start() error
	Stop() error
	Get() (float64, error)
}

type TempForecastController interface {
	Start() error
	Stop() error
	Get() (float64, error)
}

type PowerSensorController interface {
	Start() error
	Stop() error
	Get() (int, error)
}

type BoilerController interface {
	Start() error
	Stop() error
	TurnOn() error
	TurnOff() error
	Get() (int, error)
}
