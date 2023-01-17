package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/3rubasa/shagent/controllers"
)

type WebServer interface {
	Initialize() error
	Start() error
	Stop() error
}

const port = 8888

type webServer struct {
	mux    *http.ServeMux
	srv    *http.Server
	boiler controllers.BoilerController
}

func New(boiler controllers.BoilerController) WebServer {
	return &webServer{
		boiler: boiler,
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
