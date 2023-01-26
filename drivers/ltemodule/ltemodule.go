package ltemodule

type LTEModule struct {
	deviceAPI DeviceAPI
}

func New(deviceAPI DeviceAPI) *LTEModule {
	return &LTEModule{
		deviceAPI: deviceAPI,
	}
}

func (l *LTEModule) SendUSSD(ussd string) (string, error) {
	return l.deviceAPI.SendUSSD(ussd)
}
