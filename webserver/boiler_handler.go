package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type boilerHandler struct {
	c *APIComponents
}

func NewBoilerHandler(c *APIComponents) *boilerHandler {
	return &boilerHandler{
		c: c,
	}
}

func (h *boilerHandler) GetState(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		State int    `json:"relay_state"`
		Error string `json:"error"`
	}

	response := &Response{}

	state, err := h.c.Boiler.Get()
	if err != nil {
		fmt.Println("Failed to get relay state: ", err)
		response.Error = err.Error()
	} else {
		response.State = state
	}

	json.NewEncoder(rw).Encode(response)
}

func (h *boilerHandler) TurnOn(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error string `json:"error"`
	}

	response := &Response{}

	err := h.c.Boiler.TurnOn()
	if err != nil {
		fmt.Println("Failed to turn relay on: ", err)
		response.Error = err.Error()
	}

	json.NewEncoder(rw).Encode(response)
}

func (h *boilerHandler) TurnOff(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error string `json:"error"`
	}

	response := &Response{}

	err := h.c.Boiler.TurnOff()
	if err != nil {
		fmt.Println("Failed to turn relay off: ", err)
		response.Error = err.Error()
	}

	json.NewEncoder(rw).Encode(response)
}
