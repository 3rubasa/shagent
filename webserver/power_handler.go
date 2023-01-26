package webserver

import "net/http"

type powerHandler struct {
	c *APIComponents
}

func NewPowerHandler(c *APIComponents) *powerHandler {
	return &powerHandler{
		c: c,
	}
}

func (h *powerHandler) GetPowerStatus(rw http.ResponseWriter, r *http.Request) {

}
