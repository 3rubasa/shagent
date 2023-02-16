package webserver

import (
	"context"
	"fmt"
	"net/http"
)

type WebServer interface {
	Initialize() error
	Start() error
	Stop() error
}

type webServer struct {
	port int
	mux  *http.ServeMux
	srv  *http.Server
	c    *APIComponents
	h    *handlers
}

func New(c *APIComponents, port int) WebServer {
	return &webServer{
		c:    c,
		port: port,
		h: &handlers{
			RoomLightHandler:   NewRoomLightHandler(c),
			CamLightHandler:    NewCamLightHandler(c),
			KitchenTempHandler: NewKitchenTempHandler(c),
			PowerHandler:       NewPowerHandler(c),
			BoilerHandler:      NewBoilerHandler(c),
			LTEModuleHandler:   NewLTEModuleHandler(c),
		},
	}
}

func (w *webServer) Initialize() error {
	w.mux = http.NewServeMux()

	SetupRoutes(w.mux, w.h)

	w.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", w.port),
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
