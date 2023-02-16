package weatherprovider

type CurWeatherResponse struct {
	Count int                `json:"count"`
	Data  []ResponseDataItem `json:"data"`
}

type ResponseDataItem struct {
	City    string  `json:"city_name"`
	Country string  `json:"country_code"`
	Temp    float64 `json:"temp"`
}
