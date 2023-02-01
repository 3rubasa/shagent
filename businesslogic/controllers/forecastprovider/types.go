package forecastprovider

type ForecastResponse struct {
	City    string             `json:"city_name"`
	Country string             `json:"country_code"`
	Data    []ResponseDataItem `json:"data"`
}

type ResponseDataItem struct {
	MinTemp   float64 `json:"min_temp"`
	MaxTemp   float64 `json:"max_temp"`
	LowTemp   float64 `json:"low_temp"`
	HighTemp  float64 `json:"high_temp"`
	ValidDate string  `json:"valid_date"`
}
