package businesslogic

type BusinessLogicComponents struct {
	RoomLight    RoomLightController
	CamLight     RoomLightController
	KitchenTemp  TempSensorController
	WindowTemp   TempSensorController
	OutdoorsTemp TempSensorController
	WeatherTemp  TempSensorController
	ForecastTemp TempForecastController
	Power        PowerSensorController
	Boiler       BoilerController
}

type state struct {
	RoomLightState      int
	RoomLightStateValid bool
	CamLightState       int
	CamLightStateValid  bool
	KitchenTemp         float64
	KitchenTempValid    bool
	WindowTemp          float64
	WindowTempValid     bool
	OutdoorsTemp        float64
	OutdoorsTempValid   bool
	WeatherTemp         float64
	WeatherTempValid    bool
	ForecastedTemp      float64
	ForecastedTempValid bool
	Power               int
	PowerValid          bool
	BoilerState         int
	BoilerStateValid    bool
}

type BusinessRules struct {
	TempControlTable map[ForecestTempRange]map[WeatherTempRange]RoomTempRange
}

type tempRange struct {
	Min float64
	Max float64
}

type WeatherTempRange tempRange
type RoomTempRange tempRange
type ForecestTempRange tempRange
