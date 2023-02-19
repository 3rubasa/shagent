package webserver

import "github.com/3rubasa/shagent/pkg/businesslogic"

type APIComponents struct {
	RoomLight   RelayController
	CamLight    RelayController
	KitchenTemp TempSensorController
	Power       PowerSensorController
	Boiler      RelayController
	MC          businesslogic.MainController
}

type handlers struct {
	RoomLightHandler   *roomLightHandler
	CamLightHandler    *camLightHandler
	KitchenTempHandler *kitchenTempHandler
	PowerHandler       *powerHandler
	BoilerHandler      *boilerHandler
	LTEModuleHandler   *lteModuleHandler
}
