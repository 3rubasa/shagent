package businesslogic

type MainController interface { // Interface for UI
	// Kitchen Temperture
	GetKitchenTemp() (float64, error)

	// Power
	GetPowerState() (int, error)

	// Boiler
	GetBoilerState() (int, error)
	TurnOnBoiler() error
	TurnOffBoiler() error

	// Room Light
	GetRoomLightState() (int, error)
	TurnOnRoomLight() error
	TurnOffRoomLight() error

	// Cam Light
	GetCamLightState() (int, error)
	TurnOnCamLight() error
	TurnOffCamLight() error

	// Cellular
	GetCellAccBalance() (float64, error)
	GetCellInetBalance() (float64, error)
	GetCellTariff() (string, error)
	GetCellPhoneNumber() (string, error)
}

type RoomLightController interface {
	Start() error
	Stop() error
	Get() (int, error)
	TurnOn() error
	TurnOff() error
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

type LTEModuleController interface {
	Start() error
	Stop() error
	GetAccountBalance() (float64, error)
	GetInetBalance() (float64, error)
	GetTariff() (string, error)
	GetPhoneNumber() (string, error)
}
