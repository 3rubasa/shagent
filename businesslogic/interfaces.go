package businesslogic

type RoomLightController interface {
	Start() error
	Stop() error
}

type TempSensorController interface {
	Start() error
	Stop() error
	Get() (float64, error)
}

type PowerSensorController interface {
	Start() error
	Stop() error
	Get() (int, error)
}
