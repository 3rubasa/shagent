package ltemodule

type DeviceAPI interface {
	SendUSSD(string) (string, error)
}
