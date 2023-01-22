package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/3rubasa/shagent/controllers"

	"github.com/sitec-systems/gmodem"
)

type WebServer interface {
	Initialize() error
	Start() error
	Stop() error
}

const port = 8888

type webServer struct {
	mux       *http.ServeMux
	srv       *http.Server
	boiler    controllers.RelayController
	roomLight controllers.RelayController
	camLight  controllers.RelayController
}

func New(boiler, roomLight, camLight controllers.RelayController) WebServer {
	return &webServer{
		boiler:    boiler,
		roomLight: roomLight,
		camLight:  camLight,
	}
}

func (w *webServer) Initialize() error {
	w.mux = http.NewServeMux()
	w.mux.HandleFunc("/controllers/boiler/get_state", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			State string `json:"relay_state"`
			Error string `json:"error"`
		}

		response := &Response{}

		state, err := w.boiler.GetState()
		if err != nil {
			fmt.Println("Failed to get relay state: ", err)
			response.Error = err.Error()
		} else {
			response.State = state
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/controllers/boiler/turn_on", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error string `json:"error"`
		}

		response := &Response{}

		err := w.boiler.TurnOn()
		if err != nil {
			fmt.Println("Failed to turn relay on: ", err)
			response.Error = err.Error()
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/controllers/boiler/turn_off", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error string `json:"error"`
		}

		response := &Response{}

		err := w.boiler.TurnOff()
		if err != nil {
			fmt.Println("Failed to turn relay off: ", err)
			response.Error = err.Error()
		}

		json.NewEncoder(rw).Encode(response)
	})
	w.mux.HandleFunc("/controllers/room_light/get_state", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			State string `json:"relay_state"`
			Error string `json:"error"`
		}

		response := &Response{}

		state, err := w.roomLight.GetState()
		if err != nil {
			fmt.Println("Failed to get relay state: ", err)
			response.Error = err.Error()
		} else {
			response.State = state
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/controllers/room_light/turn_on", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error string `json:"error"`
		}

		response := &Response{}

		err := w.roomLight.TurnOn()
		if err != nil {
			fmt.Println("Failed to turn relay on: ", err)
			response.Error = err.Error()
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/controllers/room_light/turn_off", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error string `json:"error"`
		}

		response := &Response{}

		err := w.roomLight.TurnOff()
		if err != nil {
			fmt.Println("Failed to turn relay off: ", err)
			response.Error = err.Error()
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/controllers/cam_light/get_state", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			State string `json:"relay_state"`
			Error string `json:"error"`
		}

		response := &Response{}

		state, err := w.camLight.GetState()
		if err != nil {
			fmt.Println("Failed to get relay state: ", err)
			response.Error = err.Error()
		} else {
			response.State = state
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/controllers/cam_light/turn_on", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error string `json:"error"`
		}

		response := &Response{}

		err := w.camLight.TurnOn()
		if err != nil {
			fmt.Println("Failed to turn relay on: ", err)
			response.Error = err.Error()
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/controllers/cam_light/turn_off", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error string `json:"error"`
		}

		response := &Response{}

		err := w.camLight.TurnOff()
		if err != nil {
			fmt.Println("Failed to turn relay off: ", err)
			response.Error = err.Error()
		}

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/info/cellular/balance", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   string `json:"error"`
			Balance string `json:"balance"`
		}
		response := &Response{}

		mdm := &gmodem.Modem{
			DevFile:     "/dev/ttyUSB2",
			ReadTimeout: 20 * time.Second,
		}

		err := mdm.Open()
		if err != nil {
			errMsg := fmt.Sprintln("Failed to open Serial device: ", err.Error())
			fmt.Println(errMsg)
			response.Error = errMsg
			json.NewEncoder(rw).Encode(response)
			return
		}

		str, err := mdm.SendAt(`AT+CUSD=1,"*111#",15`)
		if err != nil {
			errMsg := fmt.Sprintln("Failed to send AT command: ", err)
			fmt.Println(errMsg)
			response.Error = errMsg
			json.NewEncoder(rw).Encode(response)
			return
		}

		fmt.Println("Cellular Balance: ", str)
		response.Balance = str

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/info/cellular/internet", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   string `json:"error"`
			Balance string `json:"balance"`
		}
		response := &Response{}

		mdm := &gmodem.Modem{
			DevFile:     "/dev/ttyUSB2",
			ReadTimeout: 20 * time.Second,
		}

		err := mdm.Open()
		if err != nil {
			errMsg := fmt.Sprintln("Failed to open Serial device: ", err.Error())
			fmt.Println(errMsg)
			response.Error = errMsg
			json.NewEncoder(rw).Encode(response)
			return
		}

		str, err := mdm.SendAt(`AT+CUSD=1,"*121#",15`)
		if err != nil {
			errMsg := fmt.Sprintln("Failed to send AT command: ", err)
			fmt.Println(errMsg)
			response.Error = errMsg
			json.NewEncoder(rw).Encode(response)
			return
		}

		fmt.Println("Cellular Internet Balance: ", str)
		response.Balance = str

		json.NewEncoder(rw).Encode(response)
	})

	w.mux.HandleFunc("/info/cellular/tariff", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error  string `json:"error"`
			Tariff string `json:"tariff"`
		}
		response := &Response{}

		mdm := &gmodem.Modem{
			DevFile:     "/dev/ttyUSB2",
			ReadTimeout: 20 * time.Second,
		}

		err := mdm.Open()
		if err != nil {
			errMsg := fmt.Sprintln("Failed to open Serial device: ", err.Error())
			fmt.Println(errMsg)
			response.Error = errMsg
			json.NewEncoder(rw).Encode(response)
			return
		}

		str, err := mdm.SendAt(`AT+CUSD=1,"*161#",15`)
		if err != nil {
			errMsg := fmt.Sprintln("Failed to send AT command: ", err)
			fmt.Println(errMsg)
			response.Error = errMsg
			json.NewEncoder(rw).Encode(response)
			return
		}

		fmt.Println("Cellular Tariff: ", str)
		response.Tariff = str

		json.NewEncoder(rw).Encode(response)
	})

	w.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: w.mux,
	}

	return nil
}

func (w *webServer) Start() error {
	go func() {
		fmt.Println("Starting the server...")

		err := w.srv.ListenAndServe()
		if err != nil {
			fmt.Printf("Server exited with an error: %s \n", err.Error())
		}
	}()

	return nil
}

func (w *webServer) Stop() error {
	fmt.Println("Shutting down the server...")
	err := w.srv.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("Shutdown retured with an error: %s \n", err.Error())
	}

	return nil
}
