package sonoffr3rf

type RelayInfo struct {
	Error int           `json:"error"`
	Data  RelayInfoData `json:"data"`
}

type RelayInfoData struct {
	Switch string `json:"switch"`
}

type RelaySwitchOnResponse struct {
	Error int `json:"error"`
}

type RelaySwitchOffResponse struct {
	Error int `json:"error"`
}
