package businesslogic

func (bl *BusinessLogic) GetKitchenTemp() (float64, error) {
	return bl.c.KitchenTemp.Get()
}

// Power
func (bl *BusinessLogic) GetPowerState() (int, error) {
	return bl.c.Power.Get()
}

// Boiler
func (bl *BusinessLogic) GetBoilerState() (int, error) {
	return bl.c.Boiler.Get()
}

func (bl *BusinessLogic) TurnOnBoiler() error {
	return bl.c.Boiler.TurnOn()
}

func (bl *BusinessLogic) TurnOffBoiler() error {
	return bl.c.Boiler.TurnOff()
}

// Room light
func (bl *BusinessLogic) GetRoomLightState() (int, error) {
	return bl.c.RoomLight.Get()
}

func (bl *BusinessLogic) TurnOnRoomLight() error {
	return bl.c.RoomLight.TurnOn()
}

func (bl *BusinessLogic) TurnOffRoomLight() error {
	return bl.c.RoomLight.TurnOff()
}

func (bl *BusinessLogic) GetCamLightState() (int, error) {
	return bl.c.CamLight.Get()
}

func (bl *BusinessLogic) TurnOnCamLight() error {
	return bl.c.CamLight.TurnOn()
}

func (bl *BusinessLogic) TurnOffCamLight() error {
	return bl.c.CamLight.TurnOff()
}

func (bl *BusinessLogic) GetCellAccBalance() (float64, error) {
	return bl.c.LTEModule.GetAccountBalance()
}

func (bl *BusinessLogic) GetCellInetBalance() (float64, error) {
	return bl.c.LTEModule.GetInetBalance()
}

func (bl *BusinessLogic) GetCellTariff() (string, error) {
	return bl.c.LTEModule.GetTariff()
}

func (bl *BusinessLogic) GetCellPhoneNumber() (string, error) {
	return bl.c.LTEModule.GetPhoneNumber()
}
