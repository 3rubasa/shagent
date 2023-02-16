package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type roomLightHandler struct {
	c *APIComponents
}

func NewRoomLightHandler(c *APIComponents) *roomLightHandler {
	return &roomLightHandler{
		c: c,
	}
}

func (h *roomLightHandler) GetState(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		State int    `json:"relay_state"`
		Error string `json:"error"`
	}

	response := &Response{}

	state, err := h.c.RoomLight.Get()
	if err != nil {
		fmt.Println("Failed to get relay state: ", err)
		response.Error = err.Error()
	} else {
		response.State = state
	}

	json.NewEncoder(rw).Encode(response)
}

func (h *roomLightHandler) TurnOn(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error string `json:"error"`
	}

	response := &Response{}

	err := h.c.RoomLight.TurnOn()
	if err != nil {
		fmt.Println("Failed to turn relay on: ", err)
		response.Error = err.Error()
	}

	json.NewEncoder(rw).Encode(response)
}

func (h *roomLightHandler) TurnOff(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error string `json:"error"`
	}

	response := &Response{}

	err := h.c.RoomLight.TurnOff()
	if err != nil {
		fmt.Println("Failed to turn relay off: ", err)
		response.Error = err.Error()
	}

	json.NewEncoder(rw).Encode(response)
}
