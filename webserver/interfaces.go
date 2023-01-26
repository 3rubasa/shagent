package webserver

type TempSensorController interface {
	Get() (float64, error)
}

type PowerSensorController interface {
	Get() (int, error)
}

type RelayController interface {
	Get() (int, error)
	TurnOn() error
	TurnOff() error
}

type LTEModuleController interface {
	GetAccountBalance() (string, error)
	GetInetBalance() (string, error)
	GetTariff() (string, error)
}
