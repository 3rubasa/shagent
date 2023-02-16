package weatherio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/3rubasa/shagent/pkg/config"
	"github.com/procyon-projects/chrono"
)

const invalidTemperature = -100

type ForecastProvider struct {
	cfg               *config.ForecastProviderConfig
	key               string
	city              string
	country           string
	firstPeriod       time.Duration
	mainPeriod        time.Duration
	scheduler         chrono.TaskScheduler
	task              chrono.ScheduledTask
	minForecastTemp   float64
	minForecastTempTS time.Time
	cacheTTL          time.Duration
	mux               *sync.Mutex
}

func New(cfg *config.ForecastProviderConfig, key, city, country string, firstPeriod, mainPeriod, cacheTTL time.Duration) *ForecastProvider {
	return &ForecastProvider{
		cfg:               cfg,
		key:               key,
		city:              city,
		country:           country,
		firstPeriod:       firstPeriod,
		mainPeriod:        mainPeriod,
		minForecastTemp:   invalidTemperature,
		minForecastTempTS: time.Date(1800, 1, 1, 0, 0, 0, 0, time.Local),
		cacheTTL:          cacheTTL,
		mux:               &sync.Mutex{},
	}
}

func (cw *ForecastProvider) Start() error {
	if !cw.cfg.Enabled {
		fmt.Println("Forecast Provider is Disabled")
		return nil
	}

	var err error
	cw.scheduler = chrono.NewDefaultTaskScheduler()
	cw.task, err = cw.scheduler.ScheduleAtFixedRate(cw.updateTemperature, cw.mainPeriod, chrono.WithTime(time.Now().Add(cw.firstPeriod)))

	if err != nil {
		fmt.Println("ForecastProvider: Failed to schedule a task: ", err)
		return err
	}

	return nil
}

func (cw *ForecastProvider) Stop() error {
	cw.task.Cancel()
	sdc := cw.scheduler.Shutdown()
	<-sdc

	return nil
}

func (cw *ForecastProvider) Get() (float64, error) {
	cw.mux.Lock()
	t := cw.minForecastTemp
	tTS := cw.minForecastTempTS
	cw.mux.Unlock()

	if time.Since(tTS) > cw.cacheTTL {
		return 0.0, errors.New("no fresh data")
	}

	return t, nil
}

func (cw *ForecastProvider) updateTemperature(ctx context.Context) {
	url := fmt.Sprintf("https://api.weatherbit.io/v2.0/forecast/daily?city=%s&country=%s&key=%s", cw.city, cw.country, cw.key)
	fmt.Printf("About to send http request to weatherbit.io: %s \n", url)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("Error while creating request: %s \n", err.Error())
		return
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %s \n", err.Error())
		return
	}

	curTime := time.Now()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status is not 200: %d", resp.StatusCode)
		fmt.Printf("Error: %s \n", err.Error())
		return
	}

	var respBody ForecastResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Printf("Error while parsing response from weatherbit.io: %s \n", err.Error())
		return
	}

	if len(respBody.Data) < 2 {
		fmt.Printf("Error while parsing response from weatherbit.io, len of array of data blocks is less than 2, it is equal to: %d \n", len(respBody.Data))
		return
	}

	if !strings.EqualFold(respBody.City, cw.city) {
		fmt.Printf("Error while parsing response from weatherbit.io, city is not equal to expected: %s \n", respBody.City)
		return
	}

	if !strings.EqualFold(respBody.Country, cw.country) {
		fmt.Printf("Error while parsing response from weatherbit.io, country is not equal to expected: %s \n", respBody.Country)
		return
	}

	todForecastDate, err := ParseDate(respBody.Data[0].ValidDate)
	if err != nil {
		fmt.Println("Failed to parse date: ", respBody.Data[0].ValidDate, "  Error: ", err)
		return
	}

	tomForecastDate, err := ParseDate(respBody.Data[1].ValidDate)
	if err != nil {
		fmt.Println("Failed to parse date: ", respBody.Data[1].ValidDate, "  Error: ", err)
		return
	}

	todDate := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, curTime.Location())
	tomDate := todDate.Add(24 * time.Hour)

	if todDate != todForecastDate || tomDate != tomForecastDate {
		fmt.Println("Dates in forcast are invalid")
		return
	}

	fmt.Printf("Forecast for today: min_tem=%f, max_temp=%f, low_tem=%f, high_temp=%f\n", respBody.Data[0].MinTemp, respBody.Data[0].MaxTemp, respBody.Data[0].LowTemp, respBody.Data[0].HighTemp)
	fmt.Printf("Forecast for tomorrow: min_tem=%f, max_temp=%f, low_tem=%f, high_temp=%f\n", respBody.Data[1].MinTemp, respBody.Data[1].MaxTemp, respBody.Data[1].LowTemp, respBody.Data[1].HighTemp)

	minTemp := GetMinTemperature(respBody.Data[0], respBody.Data[1])

	cw.mux.Lock()
	cw.minForecastTemp = minTemp
	cw.minForecastTempTS = curTime
	cw.mux.Unlock()
}

func GetMinTemperature(f1, f2 ResponseDataItem) float64 {
	return math.Min(f1.MinTemp,
		math.Min(f1.MaxTemp,
			math.Min(f1.LowTemp,
				math.Min(f1.HighTemp,
					math.Min(f2.MinTemp,
						math.Min(f2.MaxTemp,
							math.Min(f2.LowTemp, f2.HighTemp)))))))
}

func ParseDate(date string) (time.Time, error) {
	var err error

	if len(date) != 10 {
		return time.Now(), errors.New("invalid date formate, must by YYYY-MM-DD")
	}

	strY := string(([]byte(date))[0:4])
	strM := string(([]byte(date))[5:7])
	strD := string(([]byte(date))[8:10])

	y, err := strconv.ParseInt(strY, 10, 32)
	if err != nil {
		return time.Now(), errors.New("could not parse year")
	}

	if y < 0 {
		return time.Now(), errors.New("invalid year value")
	}

	m, err := strconv.ParseInt(strM, 10, 32)
	if err != nil {
		return time.Now(), errors.New("could not parse month")
	}

	if m < 1 || m > 12 {
		return time.Now(), errors.New("invalid month value")
	}

	d, err := strconv.ParseInt(strD, 10, 32)
	if err != nil {
		return time.Now(), errors.New("could not parse day")
	}
	if d < 1 || m > 31 {
		return time.Now(), errors.New("invalid day value")
	}

	return time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.Local), nil
}
