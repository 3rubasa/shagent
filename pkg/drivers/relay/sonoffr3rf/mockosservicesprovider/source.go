package mockosservicesprovider

//go:generate mockgen -destination=./mock.go -package=mockosservicesprovider github.com/3rubasa/shagent/pkg/drivers/relay/sonoffr3rf OSServicesProvider
