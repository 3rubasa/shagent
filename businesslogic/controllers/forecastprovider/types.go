package forecastprovider

type ForecastResponse struct {
	City    string             `json:"city_name"`
	Country string             `json:"country_code"`
	Data    []ResponseDataItem `json:"data"`
}

type ResponseDataItem struct {
	MinTemp   float64 `json:"min_temp"`
	ValidDate string  `json:"valid_date"`
}
