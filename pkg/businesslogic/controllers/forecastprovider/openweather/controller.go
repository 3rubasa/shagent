package openweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/3rubasa/shagent/pkg/config"
	"github.com/procyon-projects/chrono"
)

const invalidTemperature = -100
const stebnykLat = 49.2942
const stebnykLon = 23.5636

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
		log.Println("NOTICE: OpenWeather Forecast Provider is Disabled in config")
		return nil
	}

	var err error
	cw.scheduler = chrono.NewDefaultTaskScheduler()
	cw.task, err = cw.scheduler.ScheduleAtFixedRate(cw.updateTemperature, cw.mainPeriod, chrono.WithTime(time.Now().Add(cw.firstPeriod)))

	if err != nil {
		log.Println("ERROR: ForecastProvider: Failed to schedule a task: ", err)
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
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", stebnykLat, stebnykLon, cw.key)
	log.Printf("Debug: About to send http request to Open Weather: %s \n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ERROR: Error while creating request: ", err)
		return
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("NOTICE: Error while sending request: ", err)
		return
	}

	curTime := time.Now()

	if resp.StatusCode != 200 {
		log.Println("ERROR: HTTP response status is not 200: ", resp.StatusCode)
		return
	}

	var respBody ForecastResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Println("ERROR: Failed to parse response from openweather: ", err)
		return
	}

	if len(respBody.List) < 8 {
		log.Println("ERROR: Response from Open Weather is unexpected, length of list items is less than 8: ", len(respBody.List))
		return
	}

	if !strings.EqualFold(respBody.City.Name, cw.city) {
		log.Println("ERROR: Error while parsing response from OpenWeather, city is not equal to expected: ", respBody.City.Name)
		return
	}

	if !strings.EqualFold(respBody.City.Country, cw.country) {
		log.Println("ERROR: Error while parsing response from OpenWeather, country is not equal to expected: ", respBody.City.Country)
		return
	}

	// Test that timestamps of the first 8 list items are within 24 + 3 hours from now
	// And find min temp
	minTemp := respBody.List[0].Main.MinTemp

	for i := 0; i < 8; i++ {
		ts, err := ParseDate(respBody.List[i].TimeStampStr, time.UTC)
		if err != nil {
			log.Println("ERROR: Failed to parse timestamp ", respBody.List[i].TimeStampStr)
			return
		}

		if curTime.Add(27*time.Hour).Unix() < ts.Unix() {
			log.Println("ERROR: timestamp ", respBody.List[i].TimeStampStr, " is not within expected range from ", curTime, " to ", curTime.Add(27*time.Hour))
			return
		}

		if minTemp > respBody.List[i].Main.MinTemp {
			minTemp = respBody.List[i].Main.MinTemp
		}
	}

	log.Println("Debug: From OpenWeather: Minimum Temperature for the next 24 hrs is ", minTemp)

	cw.mux.Lock()
	cw.minForecastTemp = minTemp
	cw.minForecastTempTS = curTime
	cw.mux.Unlock()
}

func ParseDate(date string, loc *time.Location) (time.Time, error) {
	var err error

	if len(date) != 19 {
		log.Println("ERROR: invalid date format, must be YYYY-MM-DD HH:MM:SS actual: ", date)
		return time.Now(), errors.New("invalid date format, must be YYYY-MM-DD HH:MM:SS")
	}

	strY := string(([]byte(date))[0:4])
	strM := string(([]byte(date))[5:7])
	strD := string(([]byte(date))[8:10])
	strHr := string(([]byte(date))[11:13])
	strMin := string(([]byte(date))[14:16])
	strSec := string(([]byte(date))[17:19])

	y, err := strconv.ParseInt(strY, 10, 32)
	if err != nil {
		log.Println("ERROR: could not parse year: ", strY, " Error: ", err)
		return time.Now(), errors.New("could not parse year")
	}

	if y < 0 {
		log.Println("ERROR: Invalid year value: ", y)
		return time.Now(), errors.New("invalid year value")
	}

	m, err := strconv.ParseInt(strM, 10, 32)
	if err != nil {
		log.Println("ERROR: could not parse month: ", strM, " Error: ", err)
		return time.Now(), errors.New("could not parse month")
	}

	if m < 1 || m > 12 {
		log.Println("ERROR: Invalid month value: ", m)
		return time.Now(), errors.New("invalid month value")
	}

	d, err := strconv.ParseInt(strD, 10, 32)
	if err != nil {
		log.Println("ERROR: could not parse day: ", strD, " Error: ", err)
		return time.Now(), errors.New("could not parse day")
	}
	if d < 1 || m > 31 {
		log.Println("ERROR: Invalid day value: ", d)
		return time.Now(), errors.New("invalid day value")
	}

	hr, err := strconv.ParseInt(strHr, 10, 32)
	if err != nil {
		log.Println("ERROR: could not parse hour: ", strHr, " Error: ", err)
		return time.Now(), errors.New("could not parse hour")
	}
	if hr < 0 || hr > 23 {
		log.Println("ERROR: Invalid hour value: ", hr)
		return time.Now(), errors.New("invalid hour value")
	}

	min, err := strconv.ParseInt(strMin, 10, 32)
	if err != nil {
		log.Println("ERROR: could not parse minutes: ", strMin, " Error: ", err)
		return time.Now(), errors.New("could not parse minutes")
	}
	if min < 0 || min > 59 {
		log.Println("ERROR: Invalid minutes value: ", min)
		return time.Now(), errors.New("invalid minutes value")
	}

	sec, err := strconv.ParseInt(strSec, 10, 32)
	if err != nil {
		log.Println("ERROR: could not parse seconds: ", strSec, " Error: ", err)
		return time.Now(), errors.New("could not parse seconds")
	}
	if sec < 0 || sec > 59 {
		log.Println("ERROR: Invalid seconds value: ", sec)
		return time.Now(), errors.New("invalid seconds value")
	}

	return time.Date(int(y), time.Month(m), int(d), int(hr), int(min), int(sec), 0, loc), nil
}
