package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebServer interface {
	Initialize() error
	Start() error
	Stop() error
}

const port = 8888

type webServer struct {
	mux *http.ServeMux
	srv *http.Server
}

func New() WebServer {
	return &webServer{}
}

func (w *webServer) Initialize() error {
	w.mux = http.NewServeMux()
	w.mux.HandleFunc("/controllers/boiler/get_state", func(rw http.ResponseWriter, r *http.Request) {
		type Response struct {
			R string `json:"response"`
		}

		response := &Response{
			R: "Hello World",
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
