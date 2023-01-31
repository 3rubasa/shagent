package weatherprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/procyon-projects/chrono"
)

const invalidTemperature = -100

type CurrentWeather struct {
	key         string
	city        string
	country     string
	firstPeriod time.Duration
	mainPeriod  time.Duration
	scheduler   chrono.TaskScheduler
	task        chrono.ScheduledTask
	curTemp     float64
	curTempTS   time.Time
	cacheTTL    time.Duration
	mux         *sync.Mutex
}

func New(key, city, country string, firstPeriod, mainPeriod, cacheTTL time.Duration) *CurrentWeather {
	return &CurrentWeather{
		key:         key,
		city:        city,
		country:     country,
		firstPeriod: firstPeriod,
		mainPeriod:  mainPeriod,
		curTemp:     invalidTemperature,
		curTempTS:   time.Date(1800, 1, 1, 0, 0, 0, 0, time.Local),
		cacheTTL:    cacheTTL,
		mux:         &sync.Mutex{},
	}
}

func (cw *CurrentWeather) Start() error {
	var err error
	cw.scheduler = chrono.NewDefaultTaskScheduler()
	cw.task, err = cw.scheduler.ScheduleAtFixedRate(cw.updateTemperature, cw.mainPeriod, chrono.WithTime(time.Now().Add(cw.firstPeriod)))

	if err != nil {
		fmt.Println("CurrentWeather: Failed to schedule a task: ", err)
		return err
	}

	return nil
}

func (cw *CurrentWeather) Stop() error {
	cw.task.Cancel()
	sdc := cw.scheduler.Shutdown()
	<-sdc

	return nil
}

func (cw *CurrentWeather) Get() (float64, error) {
	cw.mux.Lock()
	t := cw.curTemp
	tTS := cw.curTempTS
	cw.mux.Unlock()

	if time.Since(tTS) > cw.cacheTTL {
		return 0.0, errors.New("no fresh data")
	}

	return t, nil
}

func (cw *CurrentWeather) updateTemperature(ctx context.Context) {
	url := fmt.Sprintf("https://api.weatherbit.io/v2.0/current?city=%s&country=%s&key=%s", cw.city, cw.country, cw.key)
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

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status is not 200: %d", resp.StatusCode)
		fmt.Printf("Error: %s \n", err.Error())
		return
	}

	var respBody CurWeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Printf("Error while parsing response from weatherbit.io: %s \n", err.Error())
		return
	}

	if respBody.Count != 1 {
		fmt.Printf("Error while parsing response from weatherbit.io, number of data blocks is not 1, it is equal to: %d \n", respBody.Count)
		return
	}

	if len(respBody.Data) != 1 {
		fmt.Printf("Error while parsing response from weatherbit.io, len of array of data blocks is not 1, it is equal to: %d \n", len(respBody.Data))
		return
	}

	if !strings.EqualFold(respBody.Data[0].City, cw.city) {
		fmt.Printf("Error while parsing response from weatherbit.io, city is not equal to expected: %s \n", respBody.Data[0].City)
		return
	}

	if !strings.EqualFold(respBody.Data[0].Country, cw.country) {
		fmt.Printf("Error while parsing response from weatherbit.io, country is not equal to expected: %s \n", respBody.Data[0].Country)
		return
	}

	fmt.Println("Temperature from weather.io: ", respBody.Data[0].Temp)

	cw.mux.Lock()
	cw.curTemp = respBody.Data[0].Temp
	cw.curTempTS = time.Now()
	cw.mux.Unlock()
}
