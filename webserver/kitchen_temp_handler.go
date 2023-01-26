package webserver

import "net/http"

type kitchenTempHandler struct {
	c *APIComponents
}

func NewKitchenTempHandler(c *APIComponents) *kitchenTempHandler {
	return &kitchenTempHandler{
		c: c,
	}
}

func (h *kitchenTempHandler) GetTemperature(rw http.ResponseWriter, r *http.Request) {

}
