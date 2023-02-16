package mockdeviceapi

//go:generate mockgen -destination=./mock.go -package=mockdeviceapi github.com/3rubasa/shagent/drivers/relay/asyncdecorator DeviceAPI
