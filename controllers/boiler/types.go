package boiler

type RelayInfo struct {
	Error int           `json:"error"`
	Data  RelayInfoData `json:"data"`
}

type RelayInfoData struct {
	Switch string `json:"switch"`
}
