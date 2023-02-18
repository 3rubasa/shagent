package esp8266

type SensorData struct {
	Temp float64 `json:"temperature"`
	MVCC float64 `json:"mvcc"`
	RVCC float64 `json:"rvcc"`
}
