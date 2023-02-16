package webserver

type APIComponents struct {
	RoomLight   RelayController
	CamLight    RelayController
	KitchenTemp TempSensorController
	Power       PowerSensorController
	Boiler      RelayController
	LTEModule   LTEModuleController
}

type handlers struct {
	RoomLightHandler   *roomLightHandler
	CamLightHandler    *camLightHandler
	KitchenTempHandler *kitchenTempHandler
	PowerHandler       *powerHandler
	BoilerHandler      *boilerHandler
	LTEModuleHandler   *lteModuleHandler
}
