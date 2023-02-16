package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type camLightHandler struct {
	c *APIComponents
}

func NewCamLightHandler(c *APIComponents) *camLightHandler {
	return &camLightHandler{
		c: c,
	}
}

func (h *camLightHandler) GetState(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		State int    `json:"relay_state"`
		Error string `json:"error"`
	}

	response := &Response{}

	state, err := h.c.CamLight.Get()
	if err != nil {
		fmt.Println("Failed to get relay state: ", err)
		response.Error = err.Error()
	} else {
		response.State = state
	}

	json.NewEncoder(rw).Encode(response)
}

func (h *camLightHandler) TurnOn(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error string `json:"error"`
	}

	response := &Response{}

	err := h.c.CamLight.TurnOn()
	if err != nil {
		fmt.Println("Failed to turn relay on: ", err)
		response.Error = err.Error()
	}

	json.NewEncoder(rw).Encode(response)
}

func (h *camLightHandler) TurnOff(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error string `json:"error"`
	}

	response := &Response{}

	err := h.c.CamLight.TurnOff()
	if err != nil {
		fmt.Println("Failed to turn relay off: ", err)
		response.Error = err.Error()
	}

	json.NewEncoder(rw).Encode(response)
}
