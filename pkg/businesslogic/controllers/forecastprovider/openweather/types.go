package openweather

type ForecastResponse struct {
	List []ListItem   `json:"list"`
	City ForecastCity `json:"city"`
}

type ListItem struct {
	Main         ListItemMainSection `json:"main"`
	TimeStamp    int64               `json:"dt"`
	TimeStampStr string              `json:"dt_txt"`
}

type ListItemMainSection struct {
	MinTemp float64 `json:"temp_min"`
}

type ForecastCity struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}
